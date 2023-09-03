package relationCtlModel

import (
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/userCtlModel"
)

// 关注/取消关注

type ActionReq struct {
	ToUserID   int64 `query:"to_user_id, required" vd:"$>0"`
	ActionType int32 `query:"action_type, required" vd:"$>0"` // 1-关注，2-取消关注
}

type ActionResp struct {
	baseCtlModel.APIBaseResp
}

// 获取关注列表

type FollowListReq struct {
	UserID int64 `query:"user_id, required" vd:"$>0"`
}

type FollowListResp struct {
	baseCtlModel.APIBaseResp
	Users []userCtlModel.User `json:"user_list"`
}

// 获取粉丝列表

type FollowerListReq struct {
	UserID int64 `query:"user_id, required" vd:"$>0"`
}

type FollowerListResp struct {
	baseCtlModel.APIBaseResp
	Users []userCtlModel.User `json:"user_list"`
}

// 获取好友列表

type FriendListReq struct {
	UserID int64 `query:"user_id, required" vd:"$>0"`
}

type FriendListResp struct {
	baseCtlModel.APIError
	Users []userCtlModel.User `json:"user_list"`
}
