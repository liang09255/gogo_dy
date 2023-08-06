package service

import (
	"errors"
	"main/controller/ctlModel/userCtlModel"
	"main/controller/middleware"
	"main/dal"
)

type userService struct{}

var UserService = &userService{}

func (u *userService) Register(userName string, passWord string) (response *userCtlModel.RegisterResponse, err error) {
	//校验用户是否已存在
	NotExist := dal.UserDal.IsNotExist(userName)
	if NotExist == false {
		return nil, errors.New("用户已存在")
	}
	//TODO 密码加密

	//调dal插入数据库
	var userLogin = &dal.User{Username: userName, Password: passWord}
	err = dal.UserDal.CreateUser(userLogin)
	if err != nil {
		return nil, err
	}

	//jwt分发和返回值封装
	var user dal.User
	err = dal.UserDal.GetUserByUserName(userName, &user)
	if err != nil {
		return nil, err
	}

	token, err := middleware.ReleaseToken(user.ID)
	if err != nil {
		return nil, errors.New("token分发失败")
	}

	//封装返回结果
	response.UserId = user.ID
	response.Token = token
	return
}

func (u *userService) GetUserInfo(userId int64) (response *userCtlModel.User, err error) {
	var user *dal.User
	user, err = dal.UserDal.GetUserInfoById(userId)
	if err != nil {
		return nil, err
	}
	// TODO 查询关注状态
	response = &userCtlModel.User{
		ID:            user.ID,
		Username:      user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}
	return
}
