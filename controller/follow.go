package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YJ9938/DouYin/model"
	"github.com/YJ9938/DouYin/service"
	"github.com/gin-gonic/gin"
)

type FollowResponse struct {
	Response
}

type FollowListResponse struct {
	Response
	UserList []model.UserInfo `json:"user_list,omitempty"`
}

// 关注操作  apk post请求没有user_id 参数
func RelationAction(c *gin.Context) {
	// Get Parameters
	// user_id := c.Query("user_id")
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	if token == "" || to_user_id == "" || action_type == "" {
		Error(c, 1, "获取参数失败")
		// fmt.Println("user_id:", user_id, "\ntoken:", token, "\nto_user_id:", to_user_id, "\naction_type:", action_type)
		return
	}
	actiontype, _ := strconv.ParseInt(action_type, 10, 64)

	if actiontype != 1 && actiontype != 2 {
		Error(c, 1, "操作类型不符")
		return
	}

	// verify the token
	claims := parseToken(token)
	if claims == nil {
		Error(c, 1, "身份鉴权失败")
		return
	}
	// queryId, _ := strconv.ParseInt(user_id, 10, 64)
	currentId, _ := strconv.ParseInt(claims.Id, 10, 64)
	toUserId, _ := strconv.ParseInt(to_user_id, 10, 64)

	// if queryId != currentId {
	// 	Error(c, 1, "身份鉴权失败")
	// 	return
	// 	//这里有两个参数id 和 一个token的id , 应该要判断用户id 和token是否一致
	// }

	// write to database
	followService := service.FollowService{
		CurrentUser: currentId,
		ToUser:      toUserId,
		Action_type: actiontype,
	}
	if err := followService.FollowAction(); err != nil {
		Error(c, 1, err.Error())
		return
	} else {
		msg := ""
		if actiontype == 1 {
			msg = "关注成功"
		} else {
			msg = "取消关注成功"
		}

		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  msg,
		})
	}

}

// 关注列表
func FollowList(c *gin.Context) {
	// Get Parameters
	rawId := c.Query("user_id")
	token := c.Query("token")
	if rawId == "" || token == "" {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "参数获取出错"},
		})
	}

	// verify the token
	userId, _ := strconv.ParseInt(rawId, 10, 64)
	claims := parseToken(token)
	if claims == nil || claims.Id != rawId {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户鉴权失败"},
		})
	}

	// read database
	// 复用 actiontype
	followservice := service.FollowService{
		CurrentUser: userId,
		Action_type: 1, // 这里代表 查询关注列表信息
	}

	userlist, err := followservice.UserList()
	if err != nil {
		fmt.Println("查询关注列表出错, err:", err)
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "查询关注列表出错"},
			UserList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "查询关注列表成功"},
		UserList: userlist,
	})
}

// 粉丝列表
// 逻辑和上面一样 修改actiontype 即可
func FollowerList(c *gin.Context) {
	// Get Parameters
	rawId := c.Query("user_id")
	token := c.Query("token")
	if rawId == "" || token == "" {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "参数获取出错"},
		})
	}

	// verify the token
	userId, _ := strconv.ParseInt(rawId, 10, 64)
	claims := parseToken(token)
	if claims == nil || claims.Id != rawId {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户鉴权失败"},
		})
	}

	// read database
	// 复用 actiontype
	followservice := service.FollowService{
		CurrentUser: userId,
		Action_type: 2, // 这里代表 查询粉丝列表信息
	}

	userlist, err := followservice.UserList()
	if err != nil {
		fmt.Println("查询粉丝列表出错, err:", err)
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "查询粉丝列表出错"},
			UserList: nil,
		})
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "查询粉丝列表成功"},
		UserList: userlist,
	})
}
