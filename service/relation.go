package service

import (
	"fmt"
	"main/controller/ctlModel/userCtlModel"
	"main/dal"
	"main/utils/conv"
)

type relationService struct{}

var RelationService = &relationService{}

const (
	FollowAction   = 1
	UnFollowAction = 2
)

// RelationAction 关系操作
func (r *relationService) RelationAction(userID, toUserID int64, actionType int32) error {
	if actionType == FollowAction {
		return dal.RelationDal.Follow(userID, toUserID)
	} else if actionType == UnFollowAction {
		return dal.RelationDal.Delete(userID, toUserID)
	}
	return fmt.Errorf("invalid action type, action_type: %d", actionType)
}

// GetFollowList 获取关注列表
func (r *relationService) GetFollowList(userid int64) (users []userCtlModel.User, err error) {
	var uids []int64
	uids, err = dal.RelationDal.GetAllFollow(userid)
	if err != nil {
		return
	}
	userInfos, err := UserService.MGetUserInfo(uids)
	for _, user := range userInfos {
		user.IsFollow = true
		users = append(users, user)
	}
	return
}

// GetFollowerList 获取粉丝列表
func (r *relationService) GetFollowerList(userID int64) (users []userCtlModel.User, err error) {
	var uids []int64
	uids, err = dal.RelationDal.GetAllFollower(userID)
	if err != nil {
		return
	}
	return UserService.MGetUserInfo(uids, userID)
}

// GetFriendList 获取好友列表
func (r *relationService) GetFriendList(userID int64) (users []userCtlModel.User, err error) {
	var uids []int64
	uids, err = dal.RelationDal.GetAllFriend(userID)
	userInfos, err := UserService.MGetUserInfo(uids)
	if err != nil {
		return
	}
	// 把所有IsFollow设置为true
	for _, userInfo := range userInfos {
		userInfo.IsFollow = true
		users = append(users, userInfo)
	}
	return
}

// MGetRelation 批量获取关系 返回的Map中包含了关注的用户
func (r *relationService) MGetRelation(userID int64, toUserIDs []int64) (map[int64]struct{}, error) {
	// 这个地方后面可以优化一下
	followIds, err := dal.RelationDal.GetAllFollow(userID)
	if err != nil {
		return nil, err
	}
	toUserIDMap := conv.Array2Map(toUserIDs)
	followMap := make(map[int64]struct{})
	for _, id := range followIds {
		if _, ok := toUserIDMap[id]; ok {
			followMap[id] = struct{}{}
		}
	}
	return followMap, nil
}
