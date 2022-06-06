package service

import (
	"time"

	"github.com/YJ9938/DouYin/model"
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
		// feed函数 每个视频作者可能不一样
		var videoDisplay model.VideoDisplay
		videoDisplay.Title = videoList[i].Title
		videoDisplay.Id = int64(videoList[i].ID)
		videoDisplay.CreatedAt = videoList[i].CreatedAt
		videoDisplay.PlayUrl = videoList[i].PlayURL
		videoDisplay.CoverUrl = videoList[i].CoverURL
		videoDisplay.Author, _ = userDao.QueryUserById(videoList[i].AuthorID)
		// 下面三个需要查表
		// videoDisplay.IsFavorite = false
		// videoDisplay.CommentCount = 0
		// videoDisplay.FavoriteCount = 0
		videoDisplay.IsFavorite, _ = model.NewFavoriteDao().IsFavorite(videoList[i].AuthorID, videoDisplay.Id)
		videoDisplay.CommentCount, _ = model.NewCommentDao().QueryCommentCountByVideoId(videoDisplay.Id)
		videoDisplay.FavoriteCount, _ = model.NewFavoriteDao().QueryFavoriteCountByVideoId(videoDisplay.Id)

		videoDisplayList = append(videoDisplayList, videoDisplay)
	}
	return videoDisplayList, err
}
