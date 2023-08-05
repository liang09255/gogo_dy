package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/dal"
	"main/middleware"
	"main/model"
	"strconv"
)

type relationService struct{}

var RelationService = &relationService{}

// RelationAction 关系操作
func (h *relationService) RelationAction(token string, userid string, actiontype string, response *model.DouyinRelationActionResponse) error {
	username := middleware.GetUsernameByToken(token)

	//操作用户的ID
	idA, err := dal.UserDal.GetUserIDbyname(username)

	if err != nil {
		hlog.Error(err)
		return err
	}
	//被关注或者取消关注用户的ID
	idB, err := strconv.Atoi(userid)

	if err != nil {
		hlog.Error(err)
		return err
	}
	if actiontype == "1" {
		err = dal.RelationDal.Follow(idA, (int64)(idB))
		response.StatusMsg = "关注成功"
	} else {
		err = dal.RelationDal.Delete(idA, (int64)(idB))
		response.StatusMsg = "取消关注"
	}

	if err != nil {
		hlog.Error(err)
		response.StatusCode = 1
	} else {
		response.StatusCode = 0
	}

	return err
}

// GetFollowList 获取关注列表
func (h *relationService) GetFollowList(token string, userid string, response *[]model.Userinfo) error {

	id := GetInt64Bystring(userid)

	err := dal.RelationDal.GetAllFollow(token, id, response)

	if err != nil {
		hlog.Error(err)
	}
	return err
}

// GetFollowerList 获取粉丝列表
func (h *relationService) GetFollowerList(token string, userid string, response *[]model.Userinfo) error {
	id := GetInt64Bystring(userid)

	err := dal.RelationDal.GetAllFollower(token, id, response)

	if err != nil {
		hlog.Error(err)
	}

	return err
}

// GetFriendList 获取好友列表
func (h *relationService) GetFriendList(token string, userid string, response *[]model.Userinfo) error {
	id := GetInt64Bystring(userid)

	err := dal.RelationDal.GetAllFollower(token, id, response)

	if err != nil {
		hlog.Error(err)
	}
	return err
}
