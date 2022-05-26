package controller

import (
	"fmt"
	"github.com/YJ9938/DouYin/model"
	"github.com/YJ9938/DouYin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	NextTime  int64                `json:"next_time"`
	VideoList []model.VideoDisplay `json:"video_list,omitempty"`
}

func Feed(c *gin.Context) {
	var latestTime time.Time
	timeStamp := c.Query("latest_time")
	rawTime, _ := strconv.ParseInt(timeStamp, 10, 64)
	if rawTime == 0 {
		latestTime = time.Now()
	} else {
		latestTime = time.UnixMilli(rawTime)
	}
	fmt.Println(latestTime)
	feedService := service.FeedService{
		LatestTime: latestTime,
	}
	videoList, err := feedService.QueryFeed()
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	nextTime := videoList[len(videoList)-1].CreatedAt.Unix()
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}
