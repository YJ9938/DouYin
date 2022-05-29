package model

import (
	"time"

	"gorm.io/gorm"
)

type UserAPI struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"` // 用户名称
}

type Video struct {
	gorm.Model
	AuthorID int64  `gorm:"not null; index:idx_author_id" json:"author_id"`
	Title    string `gorm:"not null" json:"title"`
	PlayURL  string `gorm:"not null" json:"play_url"`
	CoverURL string `gorm:"not null" json:"coverurl"`
}

type VideoAuthorUnion struct {
	ID            int64    `json:"id"`
	Author        *UserAPI `json:"author" gorm:"embedded;embeddedPrefix:author_"`
	PlayURL       string   `json:"play_url"`  // 视频播放地址
	CoverURL      string   `json:"cover_url"` // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"`
}

type VideoDao struct {
}

//数据表Video结构

//返回响应带Author的Video结构
type VideoDisplay struct {
	Id         int64     `json:"id,omitempty"`
	UserInfoId int64     `json:"-"`
	Author     *UserInfo `json:"author" gorm:"-"`
	PlayUrl    string    `json:"play_url,omitempty"`
	CoverUrl   string    `json:"cover_url,omitempty"`
	// FavoriteCount int64     `json:"favorite_count,omitempty"`
	// CommentCount  int64     `json:"comment_count,omitempty"`
	// IsFavorite    bool      `json:"is_favorite,omitempty"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	IsFavorite    bool      `json:"is_favorite"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func NewVideoDao() *VideoDao {
	return new(VideoDao)
}

func (v *VideoDao) AddVideo(video *Video) error {
	return DB.Create(video).Error
}

func (vd *VideoDao) QueryVideosByUserId(id int64) ([]Video, error) {
	videoList := make([]Video, 0, 30)
	err := DB.Model(&Video{}).Where("author_id = ?", id).Order("created_at ASC").Limit(30).Find(&videoList).Error
	return videoList, err
}

func (vd *VideoDao) QueryVideoByLatestTime(latestTime time.Time) ([]Video, error) {
	var videoList []Video
	err := DB.Model(&Video{}).Where("created_at<=?", latestTime).Order("created_at DESC").Limit(30).Find(&videoList).Error
	return videoList, err
}
