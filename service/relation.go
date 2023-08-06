package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/dal"
	"main/middleware"
	"strconv"
)

type DouyinRelationActionRequest struct {
	ActionType string `json:"action_type"` // 1-关注，2-取消关注
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

type DouyinRelationActionResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type DouyinRelationFollowerListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type DouyinRelationFollowerListResponse struct {
	StatusCode string                 `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string                `json:"status_msg"`  // 返回状态描述
	UserList   []dal.UserInfoResponse `json:"user_list"`   // 用户列表
}

type DouyinRelationFollowListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type DouyinRelationFollowListResponse struct {
	StatusCode string                 `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string                `json:"status_msg"`  // 返回状态描述
	UserList   []dal.UserInfoResponse `json:"user_list"`   // 用户信息列表
}

type DouyinRelationFriendListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type DouyinRelationFriendListResponse struct {
	StatusCode string                 `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string                `json:"status_msg"`  // 返回状态描述
	UserList   []dal.UserInfoResponse `json:"user_list"`   // 用户列表
}

type relationService struct{}

var RelationService = &relationService{}

// RelationAction 关系操作
func (h *relationService) RelationAction(token string, userid string, actiontype string, response *DouyinRelationActionResponse) error {

	//todo 根据token去获取ID
	id := middleware.IdentityKey

	idA, _ := strconv.Atoi(id)

	//操作用户的ID

	//被关注或者取消关注用户的ID
	idB, err := strconv.Atoi(userid)

	if err != nil {
		hlog.Error(err)
		return err
	}
	if actiontype == "1" {
		err = dal.RelationDal.Follow(int64(idA), (int64)(idB))
		response.StatusMsg = "关注成功"
	} else {
		err = dal.RelationDal.Delete(int64(idA), (int64)(idB))
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
func (h *relationService) GetFollowList(token string, userid string, response *[]dal.UserInfoResponse) error {

	id, _ := strconv.Atoi(userid)

	err := dal.RelationDal.GetAllFollow(token, int64(id), response)

	if err != nil {
		hlog.Error(err)
	}
	return err
}

// GetFollowerList 获取粉丝列表
func (h *relationService) GetFollowerList(token string, userid string, response *[]dal.UserInfoResponse) error {
	id, _ := strconv.Atoi(userid)

	err := dal.RelationDal.GetAllFollower(token, int64(id), response)

	if err != nil {
		hlog.Error(err)
	}

	return err
}

// GetFriendList 获取好友列表
func (h *relationService) GetFriendList(token string, userid string, response *[]dal.UserInfoResponse) error {

	id, _ := strconv.Atoi(userid)

	err := dal.RelationDal.GetAllFollower(token, int64(id), response)

	if err != nil {
		hlog.Error(err)
	}
	return err
}
