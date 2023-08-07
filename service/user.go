package service

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlModel/userCtlModel"
	"main/controller/middleware"
	"main/dal"
	"main/global"
	"main/utils/encrypts"
)

type userService struct{}

var UserService = &userService{}

func (u *userService) Register(userName string, passWord string) (response userCtlModel.RegisterResponse, err error) {
	//校验用户是否已存在
	NotExist := dal.UserDal.IsNotExist(userName)
	if NotExist == false {
		err = fmt.Errorf("用户已存在")
		return
	}
	hlog.Error("password: ", passWord, "salt: ", global.Config.PasswordSalt)
	// 密码加密
	passWord = encrypts.Md5(passWord + global.Config.PasswordSalt)
	hlog.Error("password", passWord)

	//调dal插入数据库
	var userLogin = &dal.User{Username: userName, Password: passWord}
	err = dal.UserDal.CreateUser(userLogin)
	if err != nil {
		return
	}

	//jwt分发和返回值封装
	var user = new(dal.User)
	err = dal.UserDal.GetUserByUserName(userName, user)
	if err != nil {
		return
	}

	token, err := middleware.ReleaseToken(user.ID)
	if err != nil {
		err = fmt.Errorf("token分发失败")
		return
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

// MGetUserInfo 获取多个用户信息，供内部调用
// userIds: 用户id列表
// uid: 当前用户id 如果传入，会查询当前用户是否关注了列表中的用户(有性能损耗)
func (u *userService) MGetUserInfo(userIds []int64, uid ...int64) (users []userCtlModel.User, err error) {
	// 获取用户基本信息
	userInfos, err := dal.UserDal.MGetUser(userIds)
	if err != nil {
		return nil, err
	}
	// 获取关注关系
	var getIsFollow bool
	var followMap map[int64]struct{}
	var isFollow func(uid int64) bool
	if len(uid) != 0 {
		getIsFollow = true
	}
	if getIsFollow {
		followMap, err = RelationService.MGetRelation(uid[0], userIds)
		if err != nil {
			hlog.Error("get user relation status failed", err)
		}
	}
	isFollow = func(uid int64) bool {
		_, ok := followMap[uid]
		return ok
	}
	// 整合到userCtlModel.User
	for _, userInfo := range userInfos {
		var user = userCtlModel.User{
			ID:            userInfo.ID,
			Username:      userInfo.Username,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      isFollow(userInfo.ID),
		}
		users = append(users, user)
	}
	return
}

func (u *userService) MGetUserInfosMap(uids []int64, myUid ...int64) (users map[int64]userCtlModel.User, err error) {
	usersArray, err := u.MGetUserInfo(uids, myUid...)
	if err != nil {
		return nil, err
	}
	users = make(map[int64]userCtlModel.User)
	for _, user := range usersArray {
		users[user.ID] = user
	}
	return
}
