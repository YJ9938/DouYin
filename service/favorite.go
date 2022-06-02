package service

import (
	"DouYin/model"
	"fmt"
)

type FavoriteService struct {
	User_id     int64
	Video_id    int64
	Action_type int64
}

func (f *FavoriteService) FavoriteAction() int {
	favoriteDao := model.NewFavoriteDao()
	if f.Action_type == 1 {
		return favoriteDao.AddFavorite(f.User_id, f.Video_id)
	} else {
		return favoriteDao.DeleteFavorite(f.User_id, f.Video_id)
	}
}

func (f *FavoriteService) FavoriteVideoIdList() ([]int64, error) {
	var list []model.Favorite
	var err error
	favoriteDao := model.NewFavoriteDao()
	list, err = favoriteDao.FavoriteList(f.User_id)
	if err != nil {
		return nil, err
	}
	videoidlist := make([]int64, 0, len(list))
	for _, v := range list {
		videoidlist = append(videoidlist, v.VideoID)
	}
	return videoidlist, nil
}

func (f *FavoriteService) FavoriteList() ([]model.VideoDisplay, error) {
	videoidlist, err := f.FavoriteVideoIdList()
	if err != nil {
		fmt.Println("查找点赞视频idlist出错,err:", err)
		return nil, err
	}

	videoInfoList := make([]model.VideoDisplay, 0, len(videoidlist))
	for _, id := range videoidlist {
		videoInfo := model.VideoDisplay{}
		video, _ := model.NewVideoDao().QueryVideoByVideoId(id)
		videoInfo.Id = int64(video.ID)
		videoInfo.PlayUrl = video.PlayURL
		videoInfo.CoverUrl = video.CoverURL
		videoInfo.Title = video.Title
		userInfoService := UserInfoService{
			CurrentUser: f.User_id,
			QueryUser:   video.AuthorID,
		}
		videoInfo.Author, _ = userInfoService.QueryUserInfoById()
		videoInfo.FavoriteCount, _ = model.NewFavoriteDao().QueryFavoriteCountByVideoId(id)

		videoInfoList = append(videoInfoList, videoInfo)
	}
	return videoInfoList, err
}
