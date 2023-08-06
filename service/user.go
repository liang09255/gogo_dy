package service

import (
	"errors"
	"main/dal"
	"main/middleware"
)

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type userService struct{}

var UserService = &userService{}

func (u *userService) Register(userName string, passWord string) (response *LoginResponse, err error) {
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
	token, err := middleware.ReleaseToken(user)
	if err != nil {
		return nil, errors.New("token分发失败")
	}
	//封装返回结果
	var r = new(LoginResponse)
	r.UserId = user.ID
	r.Token = token
	return r, nil
}

func (u *userService) GetUserInfo(userId int64, token string) (response *dal.UserInfoResponse, err error) {
	//解析token获得user_id
	_, ok := middleware.ParseToken(token)
	if ok == false {
		return nil, errors.New("token解析失败")
	}
	var infoResponse = new(dal.UserInfoResponse)
	err = dal.UserDal.GetUserInfoById(userId, infoResponse)
	if err != nil {
		return nil, err
	}
	// TODO 查询关注状态
	infoResponse.IsFollow = false
	return infoResponse, nil
}
