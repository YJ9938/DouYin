package model

import "gorm.io/gorm"

// 关注关系表
// idx_follower: 查找用户所有关注的人
// idx_followee: 查找用户所有的粉丝
type Follow struct {
	gorm.Model
	FollowerID int64 `gorm:"index:idx_follower"`
	FolloweeID int64 `gorm:"index:idx_followee"`
}

// Get the follower count of a user by its ID.
func GetFollowerCount(id int64) (int64, error) {
	var count int64
	err := db.Table("follows").Where("id = ?", id).Count(&count).Error
	if err != nil {
		return 0, nil
	}

	return count, err
}
