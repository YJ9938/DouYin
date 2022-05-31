package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 关注用户信息表
// idx_author_id: 搜索用户所有投稿视频
type Follow struct {
	gorm.Model
	FolloweeID int64 `gorm:"not null; index:idx_followee" json:"followee_id"`
	FollowerID int64 `gorm:"not null; index:idx_follower" json:"follower_id"`
}

// er 关注 ee

type FollowActionDao struct {
}

func NewFollowActionDao() *FollowActionDao {
	return new(FollowActionDao)
}

// from 关注 to
func (f *FollowActionDao) AddFollow(from, to int64) error {
	var count int64 = 0
	if err := DB.Table("follows").Where("followee_id = ? AND follower_id = ?", to, from).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		fmt.Println("已关注,请勿重复关注")
		return errors.New("已关注,请勿重复关注")
	}
	follow := &Follow{
		FollowerID: from,
		FolloweeID: to,
	}
	return DB.Create(follow).Error
}

func (f *FollowActionDao) DeleteFollow(from, to int64) error {
	return DB.Where("follower_id = ? AND followee_id = ?", from, to).Delete(&Follow{}).Error
}

// 查找关注数 根据er的id 查找 ee list
func (f *FollowActionDao) CountFollowee(id int64) ([]Follow, error) {
	list := make([]Follow, 0, 30)
	tx := DB.Model(&Follow{}).Where("follower_id = ?", id).Find(&list)
	return list, tx.Error
}

// 查找粉丝数 根据ee的id 查找 er list
func (f *FollowActionDao) CountFollower(id int64) ([]Follow, error) {
	list := make([]Follow, 0, 30)
	tx := DB.Model(&Follow{}).Where("followee_id = ?", id).Find(&list)
	return list, tx.Error
}

func (f *FollowActionDao) IsFollow(from, to int64) bool {
	var num int64 = 0
	DB.Model(&Follow{}).Where("follower_id = ? AND followee_id = ?", from, to).Count(&num)
	if num == 0 {
		return false
	} else {
		return true
	}
	// return num, tx.Error
}
