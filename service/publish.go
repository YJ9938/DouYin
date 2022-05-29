package service

import (
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
	if err != nil {
		return nil, err
	}

	//获得作者信息
	var userInfo *model.UserInfo
	videoDisplayList := make([]model.VideoDisplay, 0, 30)
	userInfo, err = model.NewUserDao().QueryUserById(p.Id)
	if err != nil {
		return nil, err
	}
	// fmt.Println("videolist:", len(videoList), " videoDisplaylist:", len(videoDisplayList))

	for i := 0; i < len(videoList); i++ {
		// publishlist函数 查询的视频作者都是同一个用户
		var videoDisplay model.VideoDisplay
		videoDisplay.Title = videoList[i].Title
		videoDisplay.Id = int64(videoList[i].ID)
		videoDisplay.PlayUrl = videoList[i].PlayURL
		videoDisplay.CoverUrl = videoList[i].CoverURL
		videoDisplay.Author = userInfo
		//下面三个 需要使用函数查询
		videoDisplay.FavoriteCount = 0
		videoDisplay.CommentCount = 0
		videoDisplay.IsFavorite = false
		videoDisplayList = append(videoDisplayList, videoDisplay)
	}
	return videoDisplayList, err
}
