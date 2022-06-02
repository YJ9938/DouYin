package controller

import (
	"DouYin/model"
	"DouYin/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
 url     --> controller  --> service   -->    model
请求来了  -->  控制器      --> 业务逻辑  --> 模型层的增删改查
*/

type FavoriteActionResponse struct {
	Response
}

type FavoriteListResponse struct {
	Response
	VideoList []model.VideoDisplay `json:"video_list,omitempty"`
}

func FavoriteAction(c *gin.Context) {
	// 判断是否成功获取请求参数
	token := c.Query("token")
	rawVideoId := c.Query("video_id")
	rawActionType := c.Query("action_type")
	if token == "" || rawVideoId == "" || rawActionType == "" {
		Error(c, 1, "参数获取失败")
		return
	}
	video_id, _ := strconv.ParseInt(rawVideoId, 10, 64)
	action_type, _ := strconv.ParseInt(rawActionType, 10, 64)
	// 判断点赞类型action_type是否合法
	if action_type != 1 && action_type != 2 {
		Error(c, 1, "操作类型不符")
		return
	}
	// 判断用户token是否合法
	claims := parseToken(token)
	if claims == nil {
		Error(c, 1, "身份鉴权失败")
		return
	}
	user_id, _ := strconv.ParseInt(claims.Id, 10, 64)

	// 将请求中的数据存入数据库
	fmt.Println("user_id: ", user_id)
	fmt.Println("video_id: ", video_id)

	// 1-点赞，2-取消点赞
	favoriteService := service.FavoriteService{
		User_id:     user_id,
		Video_id:    video_id,
		Action_type: action_type,
	}
	if err := favoriteService.FavoriteAction(); err != 0 {
		if action_type == 1 {
			if err == 1 {
				Error(c, 1, "你已经对该视频点过赞")
			}
			if err == 2 {
				Error(c, 2, "点赞操作信息写入数据库出错")
			}
		} else {
			if err == 1 {
				Error(c, 1, "你没有对该视频点过赞")
			}
			if err == 2 {
				Error(c, 2, "点赞操作信息写入数据库出错")
			}
		}
		return
	} else {
		msg := ""
		if action_type == 1 {
			msg = "点赞成功"
		} else {
			msg = "取消点赞成功"
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  msg,
		})
	}
}

// FavoriteList 获取点赞列表
// 登录用户的所有点赞视频列表
// GET /douyin/favorite/list/
// https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18902464
func FavoriteList(c *gin.Context) {
	rawId := c.Query("user_id")
	token := c.Query("token")

	claims := parseToken(token)
	if claims == nil || claims.Id != rawId {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response:  Response{StatusCode: 1, StatusMsg: "token鉴权失败"},
			VideoList: nil,
		})
		return
	}

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	favoroteService := service.FavoriteService{
		User_id: user_id,
	}
	videoList, err := favoroteService.FavoriteList()
	if err != nil {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response:  Response{StatusCode: 0, StatusMsg: "获取点赞列表失败"},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response:  Response{StatusCode: 1, StatusMsg: "获取点赞列表成功"},
			VideoList: videoList,
		})
	}
}
