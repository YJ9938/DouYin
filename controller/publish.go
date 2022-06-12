package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/YJ9938/DouYin/config"
	"github.com/YJ9938/DouYin/model"
	"github.com/YJ9938/DouYin/service"
	"github.com/YJ9938/DouYin/utils/cover"
	"github.com/gin-gonic/gin"
)

//发布视频响应
type PublishResponse struct {
	Response
}

//获取登录用户的视频发布列表响应
type PublishListResponse struct {
	Response
	VideoList []model.VideoDisplay `json:"video_list,omitempty"`
}

func Publish(c *gin.Context) {
	// 从POST请求 读取参数
	title := c.PostForm("title")
	userId := c.GetInt64("id")
	// 读取data数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, PublishResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	// 存放到本地
	videoName := filepath.Base(data.Filename) // filename 应该 与 title 对应
	videoPath := filepath.Join("./public/video", videoName)
	if err := c.SaveUploadedFile(data, videoPath); err != nil {
		c.JSON(http.StatusOK, PublishResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	// 生成视频 封面 url链接
	coverPath := "./public/cover/" + videoName
	coverName, _ := cover.GenerateCover(videoPath, coverPath, 1)

	playURL := "http://" + config.C.LocalIp.Ip + ":" + config.C.LocalIp.Port + "/static/video/" + videoName
	coverURL := "http://" + config.C.LocalIp.Ip + ":" + config.C.LocalIp.Port + "/static/cover/" + coverName
	video := &model.Video{
		AuthorID: userId,
		Title:    title,
		PlayURL:  playURL,
		CoverURL: coverURL,
	}
	//调取Service层函数
	publishService := service.PublishService{
		Video: video,
	}
	if err := publishService.Publish(); err != nil {
		c.JSON(http.StatusOK, PublishResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "视频上传成功"})
}

func PublishList(c *gin.Context) {
	//获取用户id
	rawId := c.Query("user_id")
	userId, _ := strconv.ParseInt(rawId, 10, 64)
	// fmt.Println(userId)

	publishService := service.PublishService{
		Id: userId,
	}
	videoList, err := publishService.PublishList()
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, PublishListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "查询成功"},
		VideoList: videoList,
	})
}
