package model

import "gorm.io/gorm"

// 评论信息表
// idx_video_id: 查找视频ID对应的所有评论
type Comment struct {
	gorm.Model
	UserID  int64  `gorm:"not null" json:"user_id"`
	VideoID int64  `gorm:"not null; index:idx_video_id" json:"video_id"`
	Content string `gorm:"not null" json:"content"`
}

type CommentDao struct {
}

func NewCommentDao() *CommentDao {
	return new(CommentDao)
}

func (c *CommentDao) AddComment(userid, videoid int64, content string) (int64, error) {
	comment := &Comment{
		UserID:  userid,
		VideoID: videoid,
		Content: content,
	}
	return int64(comment.ID), DB.Create(comment).Error
}

func (c *CommentDao) QueryCommentCountByVideoId(videoid int64) (int64, error) {
	var count int64
	err := DB.Model(&Comment{}).Where("video_id = ? AND deleted_at IS NULL", videoid).Count(&count).Error
	return count, err
}

func (c *CommentDao) DeleteComment(comment_id int64) (int64, error) {
	return comment_id, DB.Where("id = ? AND deleted_at IS NULL", comment_id).Delete(&Comment{}).Error
}

func (c *CommentDao) CommentList(videoid int64) ([]Comment, error) {
	list := make([]Comment, 0, 30)
	tx := DB.Model(&Comment{}).Where("video_id = ? AND deleted_at IS NULL", videoid).Find(&list)
	return list, tx.Error
}
