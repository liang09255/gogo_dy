package service

import (
	"fmt"
	"main/controller/ctlModel/userCtlModel"
	"main/dal"
)

type relationService struct{}

var RelationService = &relationService{}

const (
	FollowAction   = 1
	UnFollowAction = 2
)

// RelationAction 关系操作
func (h *relationService) RelationAction(userID, toUserID int64, actionType int32) error {
	if actionType == FollowAction {
		return dal.RelationDal.Follow(userID, toUserID)
	} else if actionType == UnFollowAction {
		return dal.RelationDal.Delete(userID, toUserID)
	}
	return fmt.Errorf("invalid action type, action_type: %d", actionType)
}

// GetFollowList 获取关注列表
func (h *relationService) GetFollowList(userid int64) (users []userCtlModel.User, err error) {
	var uids []int64
	uids, err = dal.RelationDal.GetAllFollow(userid)
	if err != nil {
		return
	}
	for _, uid := range uids {
		// TODO查询用户信息
		var user = userCtlModel.User{
			ID:            uid,
			Username:      "",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		users = append(users, user)
	}
	return
}

// GetFollowerList 获取粉丝列表
func (h *relationService) GetFollowerList(userID int64) (users []userCtlModel.User, err error) {
	var uids []int64
	uids, err = dal.RelationDal.GetAllFollower(userID)
	for _, uid := range uids {
		// TODO查询用户信息
		var user = userCtlModel.User{
			ID:            uid,
			Username:      "",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		users = append(users, user)
	}
	return
}

// GetFriendList 获取好友列表
func (h *relationService) GetFriendList(userID int64) (users []userCtlModel.User, err error) {
	var uids []int64
	uids, err = dal.RelationDal.GetAllFriend(userID)

	for _, uid := range uids {
		// TODO查询用户信息
		var user = userCtlModel.User{
			ID:            uid,
			Username:      "",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		users = append(users, user)
	}
	return
}
