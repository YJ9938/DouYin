package service

import (
	"fmt"
	"github.com/YJ9938/DouYin/model"
)

type PublishService struct {
	Video *model.Video
	Id    int64
}

func (p *PublishService) Publish() error {
	videoDao := model.NewVideoDao()
	return videoDao.AddVideo(p.Video)
}

func (p *PublishService) PublishList() ([]model.VideoDisplay, error) {
	videoList, err := model.NewVideoDao().QueryVideosByUserId(p.Id)
	//获得作者信息
	var userInfo *model.UserInfo
	videoDisplayList := make([]model.VideoDisplay, 0, 30)
	userInfo, err = model.NewUserDao().QueryUserById(p.Id)
	fmt.Println("videolist:", len(videoList), " videoDisplaylist:", len(videoDisplayList))

	for i := 0; i < len(videoList); i++ {
		var videoDisplay model.VideoDisplay
		videoDisplay.Title = videoList[i].Title
		videoDisplay.Id = int64(videoList[i].ID)
		videoDisplay.PlayUrl = videoList[i].PlayURL
		videoDisplay.CoverUrl = videoList[i].CoverURL
		videoDisplay.Author = userInfo
		videoDisplayList = append(videoDisplayList, videoDisplay)
	}
	return videoDisplayList, err
}
