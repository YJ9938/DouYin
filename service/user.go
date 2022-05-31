package service

import (
	"fmt"

	"github.com/YJ9938/DouYin/model"
)

type UserInfoService struct {
	CurrentUser int64
	QueryUser   int64
}

// 查询用户信息
func (u *UserInfoService) QueryUserInfoById() (*model.UserInfo, error) {
	userDao := model.NewUserDao()
	user, err := userDao.QueryUserById(u.QueryUser)
	follow := FollowService{
		CurrentUser: u.QueryUser,
		// 查询 queryuser的信息 这里注意 #######
	}
	follow.CountFollowee()
	follow.CountFollower()

	follow.CurrentUser = u.CurrentUser
	follow.ToUser = u.QueryUser
	// 这里查询 current和 query的关系   ######
	follow.IsFollowFunc()
	user.FollowCount = follow.FolloweeCount
	user.FollowerCount = follow.FollowerCount
	user.IsFollow = follow.IsFollow
	fmt.Println("followeecount = ", user.FollowCount, " followercount = ", user.FollowerCount)

	return user, err
}
