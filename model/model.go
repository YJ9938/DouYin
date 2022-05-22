package model

import (
	"fmt"
	"log"

	"github.com/YJ9938/DouYin/config"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
)

func mySQLDSN(c *config.Config) string {
	mysql := &config.C.MySQL
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql.Username, mysql.Password, mysql.Host, mysql.Port, mysql.DBName)
}

func init() {
	// Connect to the MySQL instance.
	var err error
	dsn := mySQLDSN(&config.C)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err)
	}

	// Migrate MySQL database schema.
	if err := db.AutoMigrate(&User{}, &Video{}, &Follow{}, &Faverite{}, &Comment{}); err != nil {
		log.Fatalf("failed to migrate DB schema: %s", err)
	}

	// Connet to the Redis instance.
	if config.C.Redis.Host != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", config.C.Redis.Host, config.C.Redis.Port),
			DB:   config.C.Redis.DB,
		})
		err = rdb.Set("key", "value1", 0).Err()
		if err != nil {
			panic(err)
		}
	}
}
