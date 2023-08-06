package favoriteCtlModel

import (
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/videoCtlModel"
)

// 点赞

type ActionReq struct {
	VideoID    int64 `query:"video_id, required" vd:"$>0"`
	ActionType int32 `query:"action_type, required" vd:"$>0"` // 1:点赞 2:取消点赞
}

type ActionResp struct {
	baseCtlModel.BaseResp
}

// 获取点赞列表

type ListReq struct {
	UserID int64 `query:"user_id, required" vd:"$>0"`
}

type ListResp struct {
	baseCtlModel.BaseResp
	Videos []videoCtlModel.Video `json:"video_list"`
}
