package service

import (
	"common/ggIDL/relation"
	"context"
	"encoding/json"
	"user/internal/domain"
	"user/pkg/utils"
)

const (
	FollowAction   = 1
	UnFollowAction = 2
)

type RelationService struct {
	relation.UnsafeRelationServer
	relationDomain *domain.RelationDomain
	userDomain     *domain.UserDomain
}

var _ relation.RelationServer = (*RelationService)(nil)

func New2() *RelationService {
	return &RelationService{
		relationDomain: domain.NewRelationDomain(),
		userDomain:     domain.NewUserDomain(),
	}
}

func (rs *RelationService) Action(ctx context.Context, msg *relation.ActionRequest) (*relation.ActionResponse, error) {

	myid := msg.MyId
	touserid := msg.ToUserId

	//自己不能关注自己
	if myid == touserid {
		return &relation.ActionResponse{}, nil
	}

	err := rs.relationDomain.RelationAction(ctx, myid, touserid, msg.ActionType)

	if err == nil {
		// 如果操作成功，发送消息到 Kafka
		relationMsg := map[string]interface{}{
			"userID":   myid,
			"targetID": touserid,
			"action":   msg.ActionType,
		}
		msgData, _ := json.Marshal(relationMsg)              // 转换为JSON格式
		utils.SendMessageToKafka("relation_action", msgData) // 发送到Kafka
	}

	return &relation.ActionResponse{}, err

}

func (rs *RelationService) FollowList(ctx context.Context, request *relation.FollowListRequest) (*relation.FollowListResponse, error) {

	uids, err := rs.relationDomain.GetFollowList(ctx, request.MyId)
	if err != nil {
		return nil, err
	}
	var res []*relation.UserInfoModel
	userinfo, err := rs.userDomain.MGetUserInfo(ctx, uids)
	var tmp relation.UserInfoModel
	for _, user := range userinfo {
		tmp.IsFollow = true
		tmp.Id = user.Id
		tmp.Name = user.Name
		tmp.FollowCount = user.FollowCount
		tmp.FollowerCount = user.FollowerCount
		tmp.FavoriteCount = user.FavoriteCount
		tmp.TotalFavorited = user.TotalFavorited
		tmp.Avatar = user.Avatar
		tmp.Signature = user.Signature
		tmp.BackgroundImage = user.BackgroundImage
		tmp.WorkCount = user.WorkCount
		res = append(res, &tmp)
	}

	if err != nil {
		return nil, err
	}
	return &relation.FollowListResponse{UserInfo: res}, nil
}

func (rs *RelationService) FollowerList(ctx context.Context, request *relation.FollowerListRequest) (*relation.FollowerListResponse, error) {
	uids, err := rs.relationDomain.GetFollowerList(ctx, request.MyId)
	if err != nil {
		return nil, err
	}
	var res []*relation.UserInfoModel
	userinfo, err := rs.userDomain.MGetUserInfo(ctx, uids)
	isFollowMap := rs.relationDomain.GetIsFollow(ctx, request.MyId, uids)
	var tmp relation.UserInfoModel
	for _, user := range userinfo {
		tmp.IsFollow = isFollowMap[user.Id]
		tmp.Id = user.Id
		tmp.Name = user.Name
		tmp.FollowCount = user.FollowCount
		tmp.FollowerCount = user.FollowerCount
		tmp.FavoriteCount = user.FavoriteCount
		tmp.TotalFavorited = user.TotalFavorited
		tmp.Avatar = user.Avatar
		tmp.Signature = user.Signature
		tmp.BackgroundImage = user.BackgroundImage
		tmp.WorkCount = user.WorkCount
		res = append(res, &tmp)
	}

	if err != nil {
		return nil, err
	}
	return &relation.FollowerListResponse{UserInfo: res}, nil
}

func (rs *RelationService) FriendList(ctx context.Context, request *relation.FriendListRequest) (*relation.FriendListResponse, error) {
	uids, err := rs.relationDomain.GetFriendList(ctx, request.MyId)
	if err != nil {
		return nil, err
	}
	var res []*relation.UserInfoModel
	userinfo, err := rs.userDomain.MGetUserInfo(ctx, uids)
	var tmp relation.UserInfoModel
	for _, user := range userinfo {
		tmp.IsFollow = true
		tmp.Id = user.Id
		tmp.Name = user.Name
		tmp.FollowCount = user.FollowCount
		tmp.FollowerCount = user.FollowerCount
		tmp.FavoriteCount = user.FavoriteCount
		tmp.TotalFavorited = user.TotalFavorited
		tmp.Avatar = user.Avatar
		tmp.Signature = user.Signature
		tmp.BackgroundImage = user.BackgroundImage
		tmp.WorkCount = user.WorkCount
		res = append(res, &tmp)
	}

	if err != nil {
		return nil, err
	}

	return &relation.FriendListResponse{UserInfo: res}, nil
}
