package service

import "github.com/YJ9938/DouYin/model"

type UserInfoService struct {
	CurrentUser int64
	QueryUser   int64
}

func (u *UserInfoService) QueryUserInfoById() (*model.UserInfo, error) {
	userDao := model.NewUserDao()
	user, err := userDao.QueryUserById(u.QueryUser)
	//查询是否有关注关系
	//user.IsFollow = model.isFollow(u.CurrentUser, u.QueryUser)
	return user, err
}
