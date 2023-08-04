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

//func (u *userService) Login(userName string, passWord string) (response *LoginResponse, err error) {
//	//校验用户是否存在
//	NotExist := dal.UserDal.IsNotExist(userName)
//	if NotExist {
//		return nil, errors.New("用户不存在")
//	}
//	//校验密码
//	var user1 dal.User
//	err = dal.UserDal.GetUserByUserName(userName, &user1)
//	if err != nil {
//		return nil, err
//	}
//	if user1.Password != passWord {
//		return nil, errors.New("密码错误！")
//	}
//	//jwt分发和返回值封装
//	var user dal.User
//	err = dal.UserDal.GetUserByUserName(userName, &user)
//	if err != nil {
//		return nil, err
//	}
//	token, err := middleware.ReleaseToken(user)
//	if err != nil {
//		return nil, errors.New("token分发失败")
//	}
//	//封装返回结果
//	var r = new(LoginResponse)
//	r.UserId = user.ID
//	r.Token = token
//	return r, nil
//}

func (u *userService) GetUserInfo(userId int64, token string) (response *dal.UserInfoResponse, err error) {
	//解析token获得user_id
	claims, ok := middleware.ParseToken(token)
	if ok == false {
		return nil, errors.New("token解析失败")
	}
	ParsedUserId := claims.UserId
	//比对user_id,进行查询
	if ParsedUserId != userId {
		return nil, errors.New("token身份验证失败")
	}
	var infoResponse = new(dal.UserInfoResponse)
	err = dal.UserDal.GetUserInfoById(userId, infoResponse)
	if err != nil {
		return nil, err
	}
	return infoResponse, nil
}
