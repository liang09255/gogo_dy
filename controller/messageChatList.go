package controller

import (
	"context"
	"strconv"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/service"
)

type MessageChatController struct{}

var MessageChatControllerInstance = &MessageChatController{}

func RegMessageChat(h *server.Hertz) {
	h.GET("/douyin/message/chat/", MessageChatControllerInstance.Chat)
}

func (m *MessageChatController) Chat(c context.Context, ctx *app.RequestContext) {
	// 必需参数token
	token, ok := ctx.GetQuery("token")
	if !ok {
		BaseFailResponse(ctx, "token is required")
		return
	}

	// 必需参数to_user_id
	toUserIDStr, ok := ctx.GetQuery("to_user_id")
	if !ok {
		BaseFailResponse(ctx, "to_user_id is required")
		return
	}
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		BaseFailResponse(ctx, "to_user_id must be an integer")
		return
	}

	// 必需参数pre_msg_time
	preMsgTimeStr, ok := ctx.GetQuery("pre_msg_time")
	if !ok {
		BaseFailResponse(ctx, "pre_msg_time is required")
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
		BaseFailResponse(ctx, "get chat messages error")
		return
	}
	
	MessageSuccessResponse(ctx, "success",response)
}

type MessageChatResponse struct {
	BaseResponse
	MessageList interface{} `json:"message_list"`
}

func MessageSuccessResponse(ctx *app.RequestContext, msg string , data interface{}) {
	Response(ctx, MessageChatResponse{
		BaseResponse:BaseResponse{StatusCode: SuccessCode, StatusMsg: msg},
		MessageList :data,
	})
}
