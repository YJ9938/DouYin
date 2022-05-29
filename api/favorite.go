package api

import (
	"DouYin/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
 url     --> controller  --> logic   -->    model
请求来了  -->  控制器      --> 业务逻辑  --> 模型层的增删改查
*/

func FavoriteAction(c *gin.Context) {
	// 前端页面填写待办事项 点击提交 会发请求到这里
	// 1. 从请求中把数据拿出来
	var favo model.Favorite

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)

	favo.UserID = user_id
	favo.VideoID = video_id

	fmt.Println("favo.UserID: ", favo.UserID)
	fmt.Println("favo.VideoID: ", favo.VideoID)
	// 2. 存入数据库
	// action_type：1-点赞，2-取消点赞
	if action_type == 1 {
		flag := model.AddFavorite(favo.UserID, favo.VideoID)
		switch flag {
		case 0:
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"data": favo,
				"msg":  "favorite succeed!",
			})
		case 1:
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"data": favo,
				"msg":  "record already exists!",
			})
		case 2:
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"data": favo,
				"msg":  "database error!",
			})
		}
	}
	if action_type == 2 {
		flag := model.DeleteFavorite(favo.UserID, favo.VideoID)
		switch flag {
		case 0:
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"data": favo,
				"msg":  "delete succeed!",
			})
		case 1:
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"data": favo,
				"msg":  "record not exists!",
			})
		case 2:
			c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"data": favo,
				"msg":  "database error!",
			})
		}
	}
}

// GetFavoriteList 获取点赞列表
// 登录用户的所有点赞视频列表
// GET /douyin/favorite/list/
// https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18902464
func GetFavoriteList(c *gin.Context) {
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	data, err := model.GetFavoriteVideoList(user_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "get video list failed!",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"data": data,
		})
	}
}
