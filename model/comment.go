package model

import "gorm.io/gorm"

// 评论信息表
// idx_video_id: 查找视频ID对应的所有评论
type Comment struct {
	gorm.Model
	UserID  int64
	VideoID int64 `gorm:"index:idx_video_id"`
	Content string
}
