package model

import "gorm.io/gorm"

// 视频信息表
// idx_author_id: 搜索用户所有投稿视频
type Video struct {
	gorm.Model
	AuthorID      int64 `gorm:"index:idx_author_id"`
	Title         string
	PlayURL       string `json:"play_url"`  // 视频播放地址
	CoverURL      string `json:"cover_url"` // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"`
}

type UserAPI struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"` // 用户名称
}

type VideoAuthorUnion struct {
	ID            int64    `json:"id"`
	Author        *UserAPI `json:"author" gorm:"embedded;embeddedPrefix:author_"`
	PlayURL       string   `json:"play_url"`  // 视频播放地址
	CoverURL      string   `json:"cover_url"` // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"`
}
