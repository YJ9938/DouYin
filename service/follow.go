package service

import (
	"fmt"

	"github.com/YJ9938/DouYin/model"
)

type FollowService struct {
	CurrentUser   int64
	ToUser        int64
	Action_type   int64
	FolloweeCount int64
	FollowerCount int64
	IsFollow      bool
}

func (f *FollowService) FollowAction() error {
	followActionDao := model.NewFollowActionDao()
	if f.Action_type == 1 {
		return followActionDao.AddFollow(f.CurrentUser, f.ToUser)
	} else {
		return followActionDao.DeleteFollow(f.CurrentUser, f.ToUser)
	}
}

// 查找关注数
func (f *FollowService) CountFollowee() error {
	followActionDao := model.NewFollowActionDao()
	list, err := followActionDao.CountFollowee(f.CurrentUser)
	if err != nil {
		fmt.Println("查找关注人数出错, err:", err)
	}
	f.FolloweeCount = int64(len(list))
	return err
}

// 查找粉丝数
func (f *FollowService) CountFollower() error {
	followActionDao := model.NewFollowActionDao()
	list, err := followActionDao.CountFollower(f.CurrentUser)
	if err != nil {
		fmt.Println("查找粉丝人数出错, err:", err)
	}
	f.FollowerCount = int64(len(list))
	return err
}

func (f *FollowService) IsFollowFunc() {
	followActionDao := model.NewFollowActionDao()
	f.IsFollow = followActionDao.IsFollow(f.CurrentUser, f.ToUser)
}

func (f *FollowService) UserIdList() ([]int64, error) {
	// 根据actiontype判断是查找 关注 or 粉丝 idlist
	// 设置为 1 or 2
	// 1查找 关注者idlist  2查找 粉丝list
	var list []model.Follow
	var err error
	followActionDao := model.NewFollowActionDao()
	if f.Action_type == 1 {
		list, err = followActionDao.CountFollowee(f.CurrentUser)
	} else {
		list, err = followActionDao.CountFollower(f.CurrentUser)
	}
	if err != nil {
		return nil, err
	}

	idlist := make([]int64, 0, len(list))
	for _, v := range list {
		if f.Action_type == 1 {
			idlist = append(idlist, v.FolloweeID)
		} else {
			idlist = append(idlist, v.FollowerID)
		}
	}
	return idlist, nil
}

func (f *FollowService) UserList() ([]model.UserInfo, error) {
	// 根据当前用户id 查询所有关注用户id
	// 根据关注用户id 去查询对应的信息
	// actiontype = 1 关注列表
	// actiontype = 2 粉丝列表
	idlist, err := f.UserIdList()
	if err != nil {
		fmt.Println("查找idlist 出错, err:", err)
		return nil, err
	}

	userInfoList := make([]model.UserInfo, 0, len(idlist))
	userinfoservice := UserInfoService{
		CurrentUser: f.CurrentUser,
		// QueryUser:   id,
	}
	for _, id := range idlist {
		userinfoservice.QueryUser = id
		user, _ := userinfoservice.QueryUserInfoById()
		userInfoList = append(userInfoList, *user)
	}
	return userInfoList, nil
}
