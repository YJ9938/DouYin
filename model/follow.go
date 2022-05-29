package model

import "gorm.io/gorm"

// 关注用户信息表
// idx_author_id: 搜索用户所有投稿视频
type Follow struct {
	gorm.Model
	FollowerID int64 `gorm:"not null; index:idx_follower" json:"follower_id"`
	FolloweeID int64 `gorm:"not null; index:idx_followee" json:"followee_id"`
}
