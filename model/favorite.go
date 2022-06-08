package model

import "gorm.io/gorm"

// 用户点赞表
// idx_userid_videoid: 查找用户点赞列表，查找用户是否给某视频点了赞
type Favorite struct {
	gorm.Model
	UserID  int64 `gorm:"not null; index:idx_userid_videoid" json:"user_id"`
	VideoID int64 `gorm:"not null; index:idx_userid_videoid" json:"video_id"`
}

type FavoriteDao struct {
}

func NewFavoriteDao() *FavoriteDao {
	return new(FavoriteDao)
}

func (f *FavoriteDao) AddFavorite(userid, videoid int64) error {
	var count int64 = 0
	if err := DB.Table("favorites").Where("user_id = ? AND video_id = ? AND deleted_at IS NULL", userid, videoid).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil // 重复点赞不操作
	} else {
		favorite := &Favorite{
			UserID:  userid,
			VideoID: videoid,
		}
		return DB.Create(favorite).Error
	}
}

func (f *FavoriteDao) IsFavorite(userid, videoid int64) (bool, error) {
	var count int64 = 0
	err := DB.Model(&Favorite{}).Where("user_id = ? AND video_id = ? AND deleted_at IS NULL", userid, videoid).Count(&count).Error
	if err != nil || count == 0 {
		return false, err
	}
	return true, nil
}

// 获取视频点赞数
func (f *FavoriteDao) QueryFavoriteCountByVideoId(videoid int64) (int64, error) {
	var count int64 = 0
	err := DB.Model(&Favorite{}).Where("video_id = ? AND deleted_at IS NULL", videoid).Count(&count).Error
	return count, err
}

func (f *FavoriteDao) DeleteFavorite(userid, videoid int64) error {
	return DB.Where("user_id = ? AND video_id = ? AND deleted_at IS NULL", userid, videoid).Delete(&Favorite{}).Error
}

func (f *FavoriteDao) FavoriteList(userid int64) ([]Favorite, error) {
	list := make([]Favorite, 0, 30)
	tx := DB.Model(&Favorite{}).Where("user_id = ? AND deleted_at IS NULL", userid).Find(&list)
	return list, tx.Error
}
