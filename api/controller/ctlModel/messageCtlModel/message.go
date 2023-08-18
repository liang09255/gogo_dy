package messageCtlModel

import "main/controller/ctlModel/baseCtlModel"

// 发送消息

type ActionReq struct {
	ToUseID    int64  `query:"to_user_id, required" vd:"$>0"`
	ActionType int32  `query:"action_type, required" vd:"$>0"`
	Content    string `query:"content, required" vd:"len($)>0"`
}

type ActionResp struct {
	baseCtlModel.BaseResp
}

// 获取消息列表

type ChatReq struct {
	ToUseID    int64 `query:"to_user_id, required" vd:"$>0"`
	PreMsgTime int64 `query:"pre_msg_time" vd:"$>0"`
}

type ChatResp struct {
	baseCtlModel.BaseResp
	Messages []Message `json:"message_list"`
}

type Message struct {
	ID         int64  `json:"id"`
	ToUserID   int64  `json:"to_user_id"`
	FromUserID int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
