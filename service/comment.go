package service

import (
	"errors"

	"github.com/YJ9938/DouYin/model"
	"gorm.io/gorm"
)

type CommentService struct{}

var (
	ErrCommentPermissionDenied = errors.New("comment permission denied")
)

// Query comments by the video's ID.
func (cs CommentService) QueryComments(videoID int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := model.DB.Where("video_id = ?", videoID).Find(&comments).Error
	return comments, err
}

// Insert comments.
func (cs CommentService) AddComment(userID, videoID int64, content string) (*model.Comment, error) {
	comment := model.Comment{
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}
	if err := model.DB.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete comments by its ID.
func (cs CommentService) DelComment(userID, commentID int64) error {
	if err := model.DB.Where("id = ? AND user_id = ?", commentID, userID).
		First(&model.Comment{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrCommentPermissionDenied
		}
		return err
	}

	return model.DB.Delete(&model.Comment{}, commentID).Error
}
