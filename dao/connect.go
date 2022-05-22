package dao

import (
	"fmt"

	"github.com/YJ9938/DouYin/config"
	"github.com/YJ9938/DouYin/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func CreateTables() {
	CreateTable(model.Comment{})
	CreateTable(model.Favorite{})
	CreateTable(model.Follow{})
	CreateTable(model.User{})
	CreateTable(model.Video{})
}

func CreateTable(Model interface{}) {
	config.Init()
	database := config.C.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", database.Username,
		database.Password, database.Host, database.Port, database.DBName)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&Model)
}
