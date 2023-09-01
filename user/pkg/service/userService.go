package service

import (
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"common/ggRPC"
	"context"
	"fmt"
	"user/internal/domain"
)

type UserService struct {
	//user.UnimplementedUserServer
	user.UnsafeUserServer
	userDomain     *domain.UserDomain
	relationDomain *domain.RelationDomain
	videoClient    video.VideoServiceClient
}

var _ user.UserServer = (*UserService)(nil)

func New() *UserService {
	return &UserService{
		userDomain:     domain.NewUserDomain(),
		relationDomain: domain.NewRelationDomain(),
		videoClient:    ggRPC.GetVideoClient(),
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
	userInfo, err := us.userDomain.MGetUserInfo(ctx, userId)

	if err != nil {
		return &user.UserInfoResponse{}, err
	}
	if len(userInfo) == 0 {
		return &user.UserInfoResponse{}, fmt.Errorf("用户不存在")
	}
	//获取用户关系
	myId := request.MyId
	targetUids := make([]int64, 0, len(userInfo))
	for _, u := range userInfo {
		targetUids = append(targetUids, u.Id)
	}
	isFollowMap := us.relationDomain.GetIsFollow(ctx, myId, targetUids)
	ggLog.Debug("isFollowMap", isFollowMap)
	for _, u := range userInfo {
		u.IsFollow = isFollowMap[u.Id]
	}
	ggLog.Debugf("userInfo:%+v", userInfo)
	// TODO 可能需要优化
	for k, uid := range userId {
		// 获取用户的关注数
		followList, _ := us.relationDomain.GetFollowList(ctx, uid)
		userInfo[k].FollowCount = int64(len(followList))
		// 获取用户的粉丝数
		followerList, _ := us.relationDomain.GetFollowerList(ctx, uid)
		userInfo[k].FollowerCount = int64(len(followerList))
		// 获取用户的获赞数
		favoriteCount, err := us.videoClient.GetTotalFavoriteCount(ctx, &video.GetTotalFavoriteCountRequest{UserId: uid})
		if err != nil {
			ggLog.Error("调用video服务获取获赞数失败:", err)
		}
		userInfo[k].TotalFavorited = favoriteCount.Count
		// 获取用户的作品数
		workCount, err := us.videoClient.GetTotalVideoCount(ctx, &video.GetTotalVideoCountRequest{UserId: uid})
		if err != nil {
			ggLog.Error("调用video服务获取作品数失败:", err)
		}
		userInfo[k].WorkCount = workCount.Count
		// 获取用户的喜欢数
		likeCount, err := us.videoClient.GetTotalLikeCount(ctx, &video.GetTotalLikeCountRequest{UserId: uid})
		if err != nil {
			ggLog.Error("调用video服务获取喜欢数失败:", err)
		}
		userInfo[k].FavoriteCount = likeCount.Count
	}
	return &user.UserInfoResponse{UserInfo: userInfo}, nil
}
