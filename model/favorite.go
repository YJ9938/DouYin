package model

import "gorm.io/gorm"

// 用户点赞表
// idx_userid_videoid: 查找用户点赞列表，查找用户是否给某视频点了赞
type Faverite struct {
	gorm.Model
	UserID  int64 `gorm:"index:idx_userid_videoid"`
	VideoID int64 `gorm:"index:idx_userid_videoid"`
}
