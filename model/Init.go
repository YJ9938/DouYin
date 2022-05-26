package model

import (
	"fmt"
	"github.com/YJ9938/DouYin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

//连接数据库
func InitMySQL() (err error) {
	config.Init()
	database := config.C.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", database.Username,
		database.Password, database.Host, database.Port, database.DBName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	if err != nil {
		panic(err)
	}
	db := DB.Migrator()
	if !db.HasTable(&User{}) {
		err := db.CreateTable(&User{})
		if err != nil {
			panic(err)
		}
	}
	if !db.HasTable(&Video{}) {
		err := db.CreateTable(&Video{})
		if err != nil {
			panic(err)
		}
	}
	if !db.HasTable(&Comment{}) {
		err := db.CreateTable(&Comment{})
		if err != nil {
			panic(err)
		}
	}
	if !db.HasTable(&Favorite{}) {
		err := db.CreateTable(&Favorite{})
		if err != nil {
			panic(err)
		}
	}
	return err
}
