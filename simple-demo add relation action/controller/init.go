package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

type UserFollow struct {
	AccountID  int64 `json:"account_id" gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL"`
	FollowerID int64 `json:"follower_id" gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL"` //联合主键
	Status     bool  `json:"status" gorm:"default:1"`
}

type UserInfoTable struct {
	AccountId int64  `json:"id,omitempty" gorm:"primary_key"`
	Username  string `json:"name,omitempty"`
	// FollowCount   int64 `json:"follow_count,omitempty"`
	// FollowerCount int64 `json:"follower_count,omitempty"`
	// IsFollow      bool  `json:"is_follow,omitempty"`
	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
	IsFollow      bool  `json:"is_follow"`
}

func AddRelation(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:123@(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&UserFollow{})
	token := c.Query("token")
	ActionType := c.Query("action_type")
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	relation := UserFollow{AccountID: toUserId, FollowerID: usersLoginInfo[token].Id, Status: true}
	if ActionType == "1" {
		//关注
		result := db.Where(&UserFollow{AccountID: toUserId, FollowerID: usersLoginInfo[token].Id}).First(&relation).RowsAffected
		if result == 0 {
			db.Create(&relation)
		} else {
			db.Model(&UserFollow{}).Where(&UserFollow{AccountID: toUserId, FollowerID: usersLoginInfo[token].Id}).Update("Status", true)
		}
	} else if ActionType == "2" {
		//取消关注
		result := db.Where(&UserFollow{AccountID: toUserId, FollowerID: usersLoginInfo[token].Id}).First(&relation).RowsAffected
		if result > 0 {
			db.Model(&UserFollow{}).Where(&UserFollow{AccountID: toUserId, FollowerID: usersLoginInfo[token].Id}).Update("Status", false)
		}
	}

	//进行关注或取关操作之后需要更新对应的UserLoginInfo
	//db.AutoMigrate(&User{})
	//db.Model(&User{}).Create(map[string]interface{}{
	//	"Id": usersLoginInfo[token].Id,
	//})
	//u1 := UserInfoTable{
	//	AccountId: toUserId,
	//	IsFollow:  false,
	//}
	//db.Create(&u1)
}

