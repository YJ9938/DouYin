package model

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

// 评论信息表
// idx_video_id: 查找视频ID对应的所有评论
type Comment struct {
	gorm.Model
	UserID  int64  `gorm:"not null" json:"user_id"`
	VideoID int64  `gorm:"not null; index:idx_video_id" json:"video_id"`
	Content string `gorm:"not null" json:"content"`
}

func (c Comment) MarshalBinary() (data []byte, err error) {
	return json.Marshal(c)
}

type CommentDao struct {
}

func NewCommentDao() *CommentDao {
	return new(CommentDao)
}

/*
	commentKey maps videoID to Redis key.
	Video comments are stored in redis hash seperately, the hash key is comment:{{videoID}}
	and the hash field of each comment is its comment id.
*/
func commentKey(videoID int64) string {
	return fmt.Sprintf("comment:%d", videoID)
}

func (c *CommentDao) AddComment(userid, videoid int64, content string) (int64, error) {
	// Insert into MySQL
	comment := &Comment{
		UserID:  userid,
		VideoID: videoid,
		Content: content,
	}
	if err := DB.Create(comment).Error; err != nil {
		return 0, err
	}

	// Insert into existing cache
	ctx := context.Background()
	keyString := commentKey(videoid)
	if RDB.Exists(ctx, keyString).Val() == 0 {
		return int64(comment.ID), RDB.HSet(ctx, keyString, fmt.Sprint(comment.ID), comment).Err()
	}
	return int64(comment.ID), nil
}

func (c *CommentDao) QueryCommentCountByVideoId(videoid int64) (int64, error) {
	// If cache hits, we just return it
	ctx := context.Background()
	keyString := commentKey(videoid)
	if RDB.Exists(ctx, keyString).Val() == 1 {
		return RDB.LLen(ctx, keyString).Result()
	}

	// Otherwise we must get comments from MySQL(where cache will be updated)
	comments, err := c.CommentList(videoid)
	if err != nil {
		return 0, err
	}

	return int64(len(comments)), nil
}

func (c *CommentDao) DeleteComment(comment_id int64) (int64, error) {
	var comment Comment
	if err := DB.Where("id = ?", comment_id).Find(&comment).Error; err != nil {
		return 0, err
	}
	ctx := context.Background()
	keyString := commentKey(comment.VideoID)

	if RDB.Exists(ctx, keyString).Val() == 1 {
		if err := RDB.HDel(ctx, keyString, fmt.Sprint(comment.ID)).Err(); err != nil {
			return comment_id, err
		}
	}
	return comment_id, DB.Where("id = ? AND deleted_at IS NULL", comment_id).Delete(&Comment{}).Error
}

func (c *CommentDao) CommentList(videoid int64) ([]Comment, error) {
	ctx := context.Background()
	keyString := commentKey(videoid)
	if RDB.Exists(ctx, keyString).Val() == 1 {
		commentStrings, err := RDB.HGetAll(ctx, keyString).Result()
		if err != nil {
			return nil, err
		}

		comments := make([]Comment, len(commentStrings))
		var insertCount int
		for _, value := range commentStrings {
			json.Unmarshal([]byte(value), &comments[insertCount])
			insertCount++
		}

		// Sort the comments
		sort.Slice(comments, func(i, j int) bool {
			return comments[i].ID > comments[j].ID
		})
		return comments, nil
	}

	comments := make([]Comment, 0, 30)
	if err := DB.Model(&Comment{}).Where("video_id = ? AND deleted_at IS NULL", videoid).
		Order("id desc").Find(&comments).Error; err != nil {
		return nil, err
	}

	pairs := make([]string, 0, len(comments)*2)
	for i := 0; i < len(comments); i++ {
		bytes, _ := json.Marshal(comments[i])
		pairs = append(pairs, fmt.Sprint(comments[i].ID), string(bytes))
	}
	return comments, RDB.HSet(ctx, keyString, pairs).Err()
}
