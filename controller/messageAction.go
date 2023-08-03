package controller

import (
	"context"
	"main/service"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type MessageActionController struct{}

var MessageActionControllerInstance = &MessageActionController{}

func RegMessageAction(h *server.Hertz) {
	h.POST("/douyin/message/action/", MessageActionControllerInstance.Action)
}

func (m *MessageActionController) Action(c context.Context, ctx *app.RequestContext) {
	token, ok := ctx.GetQuery("token")
	if !ok {
		BaseFailResponse(ctx, "token is required")
		return
	}
	toUserIDStr, ok := ctx.GetQuery("to_user_id")
	if !ok {
		BaseFailResponse(ctx, "to_user_id is required")
		return
	}
	actionTypeStr, ok := ctx.GetQuery("action_type")
	if !ok {
		BaseFailResponse(ctx, "action_type is required")
		return
	}
	content, ok := ctx.GetQuery("content")
	if !ok {
		BaseFailResponse(ctx, "content is required")
		return
	}

	// 将 toUserID 和 actionType 从字符串转换为正确的类型
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		BaseFailResponse(ctx, "to_user_id must be an integer")
		return
	}
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 32)
	if err != nil {
		BaseFailResponse(ctx, "action_type must be an integer")
		return
	}

	response, err := service.MessageService.SendMessage(token, toUserID, int32(actionType), content)
	if err != nil {
		hlog.CtxErrorf(c, "send message error: %v", err) // 记录错误到日志
		BaseFailResponse(ctx, "send message error")
		return
	}

	BaseSuccessResponse(ctx, response)
}
