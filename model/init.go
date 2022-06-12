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

func init() {
	if err := InitMySQL(); err != nil {
		panic(err)
	}
}

func InitMySQL() (err error) {
	database := config.C.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", database.Username,
		database.Password, database.Host, database.Port, database.DBName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db := DB.Migrator()
	tables := []interface{}{&User{}, &Video{}, &Comment{}, &Follow{}, &Favorite{}}
	for _, table := range tables {
		if !db.HasTable(&table) {
			err := db.CreateTable(&User{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
