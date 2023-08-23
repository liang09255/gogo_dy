package commentCtlModel

import (
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/userCtlModel"
)

// 发布评论

type ActionReq struct {
	VideoID     int64  `query:"video_id, required" vd:"$>0"`
	ActionType  int32  `query:"action_type, required" vd:"$>0"` // 1-发布评论，2-删除评论
	CommentText string `query:"comment_text" vd:"len($)>0"`     // 用户填写的评论内容，在action_type=1的时候使用
	CommentID   int64  `query:"comment_id" vd:"$>0"`            // 要删除的评论id，在action_type=2的时候使用
}

type ActionResp struct {
	baseCtlModel.BaseResp
	Comment Comment `json:"comment"`
}

// 获取评论列表

type ListReq struct {
	VideoID int64 `query:"video_id, required" vd:"$>0"`
}

type ListResp struct {
	baseCtlModel.BaseResp
	Comments []Comment `json:"comment_list"`
}

type Comment struct {
	ID         int64             `json:"id"`
	User       userCtlModel.User `json:"user"`
	Content    string            `json:"content"`
	CreateDate string            `json:"create_date"`
}
