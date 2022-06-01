package service

import "DouYin/model"

type UserInfoService struct {
	CurrentUser int64
	QueryUser   int64
}

// 查询用户信息
func (u *UserInfoService) QueryUserInfoById() (*model.UserInfo, error) {
	userDao := model.NewUserDao()
	user, err := userDao.QueryUserById(u.QueryUser)
	//查询是否有关注关系
	//user.IsFollow = model.isFollow(u.CurrentUser, u.QueryUser)
	user.IsFollow = false
	return user, err
}
