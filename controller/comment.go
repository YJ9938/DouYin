package controller

import (
	"net/http"
	"strconv"

	"github.com/YJ9938/DouYin/service"
	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	Response
	Comment service.CommentData `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []service.CommentData `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	// rawId := c.Query("query")
	token := c.Query("token")
	rawVideoId := c.Query("video_id")
	rawActionType := c.Query("action_type")
	comment_text := c.Query("comment_text")
	rawcomment_id := c.Query("comment_id")

	if token == "" || rawVideoId == "" || rawActionType == "" {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "参数获取失败"},
			Comment:  service.CommentData{},
		})
		return
	}

	video_id, _ := strconv.ParseInt(rawVideoId, 10, 64)
	actiontype, _ := strconv.ParseInt(rawActionType, 10, 64)
	if actiontype != 1 && actiontype != 2 {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "操作类型不符"},
			Comment:  service.CommentData{},
		})
		return
	}

	claims := parseToken(token)
	if claims == nil {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户鉴权失败"},
			Comment:  service.CommentData{},
		})
		return
	}

	user_id, _ := strconv.ParseInt(claims.Id, 10, 64)
	comment_id, _ := strconv.ParseInt(rawcomment_id, 10, 64)
	commentService := service.CommentService{
		User_id:     user_id,
		Video_id:    video_id,
		Action_type: actiontype,
		Content:     comment_text,
		CommentId:   comment_id,
	}

	if comment, err := commentService.CommentAction(); err != nil {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 1, StatusMsg: "评论操作信息写入数据库出错"},
			Comment:  service.CommentData{},
		})
		return
	} else {
		msg := ""
		if actiontype == 1 {
			msg = "评论成功"
		} else {
			msg = "删除评论成功"
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0, StatusMsg: msg},
			Comment:  comment,
		})
	}

}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	rawVideoId := c.Query("video_id")
	if token == "" || rawVideoId == "" {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "参数获取出错"},
			CommentList: nil,
		})
		return
	}

	claims := parseToken(token)
	if claims == nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "token鉴权失败"},
			CommentList: nil,
		})
		return
	}

	user_id, _ := strconv.ParseInt(claims.Id, 10, 64)
	video_id, _ := strconv.ParseInt(rawVideoId, 10, 64)
	commentService := service.CommentService{
		Video_id: video_id,
		User_id:  user_id,
	}

	if commentList, err := commentService.CommentList(); err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "获取评论列表出错"},
			CommentList: nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 0, StatusMsg: "获取评论列表成功"},
			CommentList: commentList,
		})
	}

}
