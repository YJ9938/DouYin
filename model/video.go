package model

import "gorm.io/gorm"

// 视频信息表
// idx_author_id: 搜索用户所有投稿视频
type Video struct {
	gorm.Model
	AuthorID int64  `gorm:"not null; index:idx_author_id" json:"author_id"`
	Title    string `gorm:"not null" json:"title"`
	PlayURL  string `gorm:"not null" json:"play_url"`
	CoverURL string `gorm:"not null" json:"coverurl"`
}
