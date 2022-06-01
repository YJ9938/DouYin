package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YJ9938/DouYin/model"
	"github.com/YJ9938/DouYin/service"
	"github.com/gin-gonic/gin"
)

type commentInfo struct {
	ID         int64
	User       model.UserInfo
	Content    string
	CreateDate string `json:"create_date"`
}

// ActionType represents the type of action on the comment.
type ActionType int

const (
	AddComment ActionType = 1
	DelComment ActionType = 2
)

// GetComments is a router to get all the comments
// and its owner's info by the specified video ID.
func GetComments(c *gin.Context) {
	// query the video ID
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		Error(c, http.StatusForbidden, "未提供视频ID")
		return
	}

	// get all of the matched comments by the video ID
	var commentService service.CommentService
	comments, err := commentService.QueryComments(videoID)
	if err != nil {
		log.Printf("Error while retriving comments by video id %d: %s\n", videoID, err)
		Error(c, http.StatusInternalServerError, "服务内部错误")
		return
	}
	commentInfos := make([]commentInfo, len(comments))
	userInfos := make(map[int64]model.UserInfo)

	// get each comment's user info with least DB query
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		commentInfos[i] = commentInfo{
			ID:         int64(comment.ID),
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.String(),
		}

		userID := comments[i].UserID
		if _, ok := userInfos[userID]; !ok {
			// TODO: Can't get current logined user's ID without middleware
			userService := service.UserInfoService{
				CurrentUser: 1,
				QueryUser:   userID,
			}

			userInfo, err := userService.QueryUserInfoById()
			if err != nil {
				log.Printf("Error while retriving user info by userID %d: %s\n", userID, err)
				Error(c, http.StatusInternalServerError, "服务内部错误")
				return
			}
			userInfos[userID] = *userInfo
		}
		commentInfos[i].User = userInfos[userID]
	}

	c.JSON(http.StatusOK, struct {
		Response
		Comment []commentInfo
	}{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "请求成功",
		},
		Comment: commentInfos,
	})
}

// CommentAction is a handler to handle actions about comments.
func CommentAction(c *gin.Context) {
	// TODO: Can't get current logined user's ID without middleware
	// TODO: Reduce code repeation
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		Error(c, http.StatusForbidden, "ID参数非法")
		return
	}
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		Error(c, http.StatusForbidden, "视频参数非法")
		return
	}
	ActionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil || (ActionType != int(AddComment) && ActionType != int(DelComment)) {
		Error(c, http.StatusForbidden, "操作参数非法")
		return
	}

	if ActionType == int(AddComment) {
		addComment(c, userID, videoID)
	} else {
		delComment(c, userID)
	}
}

func addComment(c *gin.Context, userID, videoID int64) {
	content := c.Query("comment_text")
	if content == "" {
		Error(c, http.StatusForbidden, "评论内容不能为空")
		return
	}

	// create a new comment and return the new row
	var commentService service.CommentService
	comment, err := commentService.AddComment(userID, videoID, content)
	if err != nil {
		log.Printf("error while adding the comment %q from user %d to video %d\n",
			content, userID, videoID)
		Error(c, http.StatusInternalServerError, "系统内部错误")
		return
	}

	// query the user info by its ID to return
	userService := service.UserInfoService{
		CurrentUser: userID,
		QueryUser:   userID,
	}
	userInfo, err := userService.QueryUserInfoById()
	if err != nil {
		log.Printf("error while querying userinfo by id: %v\n => %s\n", userInfo, err)
	}

	c.JSON(http.StatusOK, struct {
		Response
		Comment commentInfo
	}{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "评论成功",
		},
		Comment: commentInfo{
			ID:         int64(comment.ID),
			User:       *userInfo,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Local().String(),
		},
	})
}

func delComment(c *gin.Context, userID int64) {
	commentID, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	if err != nil {
		Error(c, http.StatusForbidden, "评论参数非法")
		return
	}

	var commentService service.CommentService
	if err := commentService.DelComment(userID, commentID); err != nil {
		if err == service.ErrCommentPermissionDenied {
			Error(c, http.StatusForbidden, "这不是你的评论哦")
		} else {
			log.Printf("error while deleting the comment %d\n", commentID)
			Error(c, http.StatusInternalServerError, "系统内部错误")
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "删除评论成功",
	})
}
