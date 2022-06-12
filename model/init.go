package model

import (
	"fmt"

	"github.com/YJ9938/DouYin/config"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func init() {
	if err := initMySQL(); err != nil {
		panic(err)
	}

	initRedis()
}

func initRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.C.Redis.Host, config.C.Redis.Port),
		Password: "",
		DB:       0,
	})
}

func initMySQL() (err error) {
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
			if err := db.CreateTable(&User{}); err != nil {
				return err
			}
		}
	}
	return nil
}
