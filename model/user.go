package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUserExists    = errors.New("error user exists")
	ErrUserNotExists = errors.New("error user not exists")
)

type UserDao struct {
}

type UserInfoDao struct {
}

// 用户表
type User struct {
	Id       int64  `gorm:"primary key" json:"id,omitempty"`
	Username string `gorm:"not null; unique; size:32; index:idx_username" json:"name,omitempty"`
	Password []byte `gorm:"not null; type:varbinary(256)" json:"password"`
	// FollowCount   int64  `json:"follow_count,omitempty"`
	// FollowerCount int64  `json:"follower_count,omitempty"`
	Salt []byte `gorm:"not null; type:varbinary(32)" json:"-"`
}

//response 返回的用户信息表
type UserInfo struct {
	Id       int64  `json:"id,omitempty"`
	Username string `json:"name,omitempty"`
	// FollowCount   int64 `json:"follow_count,omitempty"`
	// FollowerCount int64 `json:"follower_count,omitempty"`
	// IsFollow      bool  `json:"is_follow,omitempty"`
	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
	IsFollow      bool  `json:"is_follow"`
}

func NewUserDao() *UserDao {
	return new(UserDao)
}

func NewUserInfoDao() *UserInfoDao {
	return new(UserInfoDao)
}

func (u *UserDao) QueryUserById(id int64) (*UserInfo, error) {

	user := &User{}
	userInfo := &UserInfo{}
	if err := DB.First(user, "Id = ?", id).Error; err != nil {
		return userInfo, err
	}
	userInfo.Id = user.Id
	userInfo.Username = user.Username
	// 这里剩下三个数据需要查表获得
	// userInfo.FollowCount = 0
	// userInfo.FollowerCount = 0
	// userInfo.IsFollow = false

	return userInfo, nil
}

func (u *UserDao) AddUser(user *User) error {
	var count int64
	if err := DB.Table("users").Where("username =?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return ErrUserExists
	}
	return DB.Create(user).Error
}

func (u *UserDao) QueryUserByName(name string) (*User, error) {
	var user User
	err := DB.Where(&User{Username: name}, "username").First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrUserNotExists
	}
	return &user, err
}
