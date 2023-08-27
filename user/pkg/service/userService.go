package service

import (
	"common/ggIDL/user"
	"context"
	"fmt"
	"user/internal/domain"
)

type UserService struct {
	user.UnsafeUserServer
	userDomain *domain.UserDomain
}

var _ user.UserServer = (*UserService)(nil)

func New() *UserService {
	return &UserService{
		userDomain: domain.NewUserDomain(),
	}
}

func (us *UserService) Login(ctx context.Context, request *user.LoginRequest) (*user.LoginResponse, error) {
	userName := request.Username
	passWord := request.Password

	uid, err := us.userDomain.Login(ctx, userName, passWord)
	return &user.LoginResponse{UserId: uid}, err
}

func (us *UserService) Register(ctx context.Context, msg *user.RegisterRequest) (*user.RegisterResponse, error) {
	userName := msg.Username
	passWord := msg.Password

	uid, err := us.userDomain.Register(ctx, userName, passWord)
	if err != nil {
		return &user.RegisterResponse{}, err
	}

	// 在业务处理逻辑最后，调用kafka的SendLog()，将删除缓存的消息发送到kafka，可以被goroutine捕获到。
	// utils.SendCache([]byte("sign"))

	return &user.RegisterResponse{UserId: uid}, nil
}

func (us *UserService) MGetUserInfo(ctx context.Context, request *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	userId := request.UserId
	// todo: 获取用户关系
	//myId := request.MyId

	userInfo, err := us.userDomain.MGetUserInfo(ctx, userId)
	if err != nil {
		return &user.UserInfoResponse{}, err
	}
	if len(userInfo) == 0 {
		return &user.UserInfoResponse{}, fmt.Errorf("用户不存在")
	}
	return &user.UserInfoResponse{UserInfo: userInfo}, nil
}
