package controller

import (
	"context"
	"main/controller/ctlFunc"
	"main/controller/ctlModel"
	"main/service"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type message struct{}

var Message = &message{}

func (m *message) Action(c context.Context, ctx *app.RequestContext) {
	token, ok := ctx.GetQuery("token")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "token is required")
		return
	}
	toUserIDStr, ok := ctx.GetQuery("to_user_id")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "to_user_id is required")
		return
	}
	actionTypeStr, ok := ctx.GetQuery("action_type")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "action_type is required")
		return
	}
	content, ok := ctx.GetQuery("content")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "content is required")
		return
	}

	// 将 toUserID 和 actionType 从字符串转换为正确的类型
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		ctlFunc.BaseFailedResp(ctx, "to_user_id must be an integer")
		return
	}
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 32)
	if err != nil {
		ctlFunc.BaseFailedResp(ctx, "action_type must be an integer")
		return
	}

	response, err := service.MessageService.SendMessage(token, toUserID, int32(actionType), content)
	if err != nil {
		hlog.CtxErrorf(c, "send message error: %v", err) // 记录错误到日志
		ctlFunc.BaseFailedResp(ctx, "send message error")
		return
	}

	ctlFunc.BaseSuccessResp(ctx, response)
}

func (m *message) Chat(c context.Context, ctx *app.RequestContext) {
	// 必需参数token
	token, ok := ctx.GetQuery("token")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "token is required")
		return
	}

	// 必需参数to_user_id
	toUserIDStr, ok := ctx.GetQuery("to_user_id")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "to_user_id is required")
		return
	}
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		ctlFunc.BaseFailedResp(ctx, "to_user_id must be an integer")
		return
	}

	// 必需参数pre_msg_time
	preMsgTimeStr, ok := ctx.GetQuery("pre_msg_time")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "pre_msg_time is required")
		return
	}
	preMsgTime, err := strconv.ParseInt(preMsgTimeStr, 10, 64)
	if err != nil {
		preMsgTime = 0
	}

	// 参数验证完毕，传入service层处理
	response, err := service.MessageService.GetChatMessages(token, toUserID, preMsgTime)
	if err != nil {
		hlog.CtxErrorf(c, "get chat messages error: %v", err) // 记录错误到日志
		ctlFunc.BaseFailedResp(ctx, "get chat messages error")
		return
	}

	MessageSuccessResponse(ctx, "success", response)
}

type MessageChatResponse struct {
	ctlModel.BaseResp
	MessageList interface{} `json:"message_list"`
}

func MessageSuccessResponse(ctx *app.RequestContext, msg string, data interface{}) {
	ctlFunc.Response(ctx, MessageChatResponse{
		BaseResp:    ctlModel.BaseResp{StatusCode: ctlFunc.SuccessCode, StatusMsg: msg},
		MessageList: data,
	})
}
