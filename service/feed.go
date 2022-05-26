package service

import (
	"github.com/YJ9938/DouYin/model"
	"time"
)

type FeedService struct {
	LatestTime time.Time
}

func (f *FeedService) QueryFeed() ([]model.VideoDisplay, error) {
	videoList, err := model.NewVideoDao().QueryVideoByLatestTime(f.LatestTime)
	//获得作者信息
	videoDisplayList := make([]model.VideoDisplay, 0, 30)
	userDao := model.NewUserDao()
	for i := range videoList {
		var videoDisplay model.VideoDisplay
		videoDisplay.Title = videoList[i].Title
		videoDisplay.Id = int64(videoList[i].ID)
		videoDisplay.CreatedAt = videoList[i].CreatedAt
		videoDisplay.PlayUrl = videoList[i].PlayURL
		videoDisplay.CoverUrl = videoList[i].CoverURL
		videoDisplay.Author, _ = userDao.QueryUserById(videoList[i].AuthorID)
		//videoDisplayList[i].IsFavorite
		videoDisplayList = append(videoDisplayList, videoDisplay)
	}
	return videoDisplayList, err
}
