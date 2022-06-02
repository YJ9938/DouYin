package service

import (
	"fmt"
	"time"

	"github.com/YJ9938/DouYin/model"
)

type CommentData struct {
	Id         int64           `json:"id"`
	User       *model.UserInfo `json:"user"`
	Content    string          `json:"content"`
	CreateDate string          `json:"create_date"`
}

type CommentService struct {
	User_id     int64
	Video_id    int64
	Action_type int64
	Content     string
	CommentId   int64
}

func (c *CommentService) CommentAction() (CommentData, error) {
	commentDao := model.NewCommentDao()
	var err error
	var CommentId int64
	commentData := CommentData{}
	if c.Action_type == 1 {
		CommentId, err = commentDao.AddComment(c.User_id, c.Video_id, c.Content)
	} else {
		CommentId, err = commentDao.DeleteComment(c.CommentId)
	}
	if err != nil {
		return commentData, err
	}

	commentData.Id = CommentId
	userInfoService := UserInfoService{
		CurrentUser: c.User_id,
		QueryUser:   c.User_id,
	}
	commentData.User, _ = userInfoService.QueryUserInfoById()
	commentData.Content = c.Content
	commentData.CreateDate = time.Now().Format("2006-01-02 15:04:05")[5:10]
	return commentData, nil
}

func (c *CommentService) CommentList() ([]CommentData, error) {
	commentDao := model.NewCommentDao()
	CommentList := make([]CommentData, 0, 30)
	list, err := commentDao.CommentList(c.Video_id)
	if err != nil {
		fmt.Println("获取评论列表出错, err:", err)
		return nil, err
	}

	userInfoService := UserInfoService{
		CurrentUser: c.User_id,
	}
	for _, v := range list {
		commentData := CommentData{}
		commentData.Id = int64(v.ID)
		userInfoService.QueryUser = v.UserID
		commentData.User, _ = userInfoService.QueryUserInfoById()
		commentData.Content = v.Content
		commentData.CreateDate = v.CreatedAt.Format("2006-01-02 15:04:05")[5:10]
		CommentList = append(CommentList, commentData)
	}
	return CommentList, nil
}
