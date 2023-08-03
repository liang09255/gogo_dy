package service

import (
	"errors"
	"main/dal"
	"main/global"
	"main/middleware"
)

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

//type UserInfoResponse struct {
//	ID            int64  `gorm:"primarykey" json:"id"`
//	Username      string `gorm:"not null" json:"name"`
//	FollowCount   int64  `gorm:"default:0" json:"follow_count"`
//	FollowerCount int64  `gorm:"default:0" json:"follower_count"`
//	IsFollow      bool   `gorm:"is_follow" json:"is_follow"`
//}

type userService struct{}

var UserService = &userService{}

func (u *userService) Register(userName string, passWord string) (response *LoginResponse, err error) {
	//校验用户是否已存在
	user1 := new(dal.User)
	global.MysqlDB.Where("username = ?", userName).First(&user1)
	if user1.ID != 0 {
		return nil, errors.New("用户名已存在")
	}
	//TODO 密码加密
	//调dal插入数据库
	var userLogin = &dal.User{Username: userName, Password: passWord}
	err = dal.UserDal.CreateUser(userLogin)
	if err != nil {
		return nil, errors.New("用户创建失败")
	}
	//jwt分发和返回值封装
	var user dal.User
	global.MysqlDB.Where("username = ?", userName).Find(&user)
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

func (u *userService) Login(userName string, passWord string) (response *LoginResponse, err error) {
	//校验用户是否已存在
	user1 := new(dal.User)
	global.MysqlDB.Where("username = ?", userName).First(&user1)
	if user1.ID == 0 {
		return nil, errors.New("用户名不存在")
	}
	//校验密码
	if user1.Password != passWord {
		return nil, errors.New("密码错误！")
	}
	//jwt分发和返回值封装
	var user dal.User
	global.MysqlDB.Where("username = ?", userName).Find(&user)
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
	_, claim, err := middleware.ParseToken(token)
	if err != nil {
		return nil, errors.New("token解析失败")
	}
	ParsedUserId := claim.UserId
	//比对user_id,进行查询
	if ParsedUserId != userId {
		return nil, errors.New("token身份验证失败")
	}
	var infoResponse = new(dal.UserInfoResponse)
	err = dal.UserDal.GetUserInfoById(userId, infoResponse)
	if err != nil {
		return nil, errors.New("用户信息查询失败")
	}
	return infoResponse, nil
}
