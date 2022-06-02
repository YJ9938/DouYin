package model

import (
	"fmt"

	"gorm.io/gorm"
)

// 用户点赞表
// idx_userid_videoid: 查找用户点赞列表，查找用户是否给某视频点了赞
type Favorite struct {
	gorm.Model
	UserID  int64 `gorm:"not null; index:idx_userid_videoid" json:"user_id"`
	VideoID int64 `gorm:"not null; index:idx_userid_videoid" json:"video_id"`
	Video   Video `gorm:"foreignkey:VideoID;association_foreignkey:ID"`
	User    User  `gorm:"foreignkey:UserID;association_foreignkey:ID"`
}

type FavoriteDao struct {
}

func NewFavoriteDao() *FavoriteDao {
	return new(FavoriteDao)
}

// AddFavorite 添加点赞信息 返回值 0-点赞成功 1-已经点赞 2-数据库错误
func (f *FavoriteDao) AddFavorite(userID int64, videoID int64) int {
	if IsFavorite(userID, videoID) {
		return 1
	}
	favorite := &Favorite{
		UserID:  userID,
		VideoID: videoID,
	}
	// 在favorite表中添加该记录
	if DB.Create(favorite).Error != nil {
		return 2
	}
	return 0
}

// DeleteFavorite 删除点赞信息 返回值 0-删除成功 1-没有点赞 2-数据库错误
func (f *FavoriteDao) DeleteFavorite(userID int64, videoID int64) int {
	if !IsFavorite(userID, videoID) {
		return 1
	}
	// 在favorite表中删除该记录
	if DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&Favorite{}).Error != nil {
		return 2
	}
	return 0
}

// 获取视频点赞数
func (f *FavoriteDao) QueryFavoriteCountByVideoId(videoid int64) (int64, error) {
	var count int64
	fmt.Println("videoid is", videoid)
	err := DB.Model(&Favorite{}).Where("video_id = ? AND deleted_at IS NULL", videoid).Select("count(*)").Group("video_id").Find(&count).Error
	fmt.Println("count is", count)
	return count, err
}

// IsFavorite 判断是否已经点赞
func IsFavorite(userID int64, videoID int64) bool {
	return DB.First(&Favorite{}, "user_id = ? and video_id = ? AND deleted_at is NULL", userID, videoID).Error == nil
}

func (f *FavoriteDao) FavoriteList(userid int64) ([]Favorite, error) {
	list := make([]Favorite, 0, 30)
	tx := DB.Model(&Favorite{}).Where("user_id = ? AND deleted_at IS NULL", userid).Find(&list)
	return list, tx.Error
}
