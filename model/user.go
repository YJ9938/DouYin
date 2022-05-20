package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUserExists    = errors.New("error user exists")
	ErrUserNotExists = errors.New("error user not exists")
)

// 用户信息表
// 每个用户注册时分配随机的salt，表中存SHA256(password, salt)
// idx_username_password: 注册、登录时根据用户名查找用户
type User struct {
	ID       int64 `gorm:"primary key"`
	Name     string
	Username string `gorm:"not null; unique; size:32; index:idx_username"`
	Password []byte `gorm:"not null; type:varbinary(256)"`
	Salt     []byte `gorm:"not null; type:varbinary(32)"`
}

// AddUser adds a user to the database.
// If there is already a user having the username, ErrUserExists is returned.
func AddUser(user *User) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Check if the account already exists.
		var count int64
		if err := tx.Table("users").Where("username = ?", user.Username).Count(&count).Error; err != nil {
			return err
		}
		if count != 0 {
			return ErrUserExists
		}

		// Create a new user.
		return tx.Save(user).Error
	})
}

// GetUserByUsername gets the user by its name.
// If there is no such a user, ErrUserNotExists is returned.
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.Where(&User{Username: username}, "username").First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrUserNotExists
	}
	return &user, err
}
