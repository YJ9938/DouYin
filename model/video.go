package model

import "gorm.io/gorm"

// 视频信息表
// idx_author_id: 搜索用户所有投稿视频
type Video struct {
	gorm.Model
	AuthorID int64 `gorm:"index:idx_author_id"`
	Title    string
	PlayURL  string
	CoverURL string
}
