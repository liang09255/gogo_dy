package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/messageCtlModel"
	"main/controller/middleware"
	"main/service"
)

type message struct{}

var Message = &message{}

func (m *message) Action(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	// bind params
	var reqObj messageCtlModel.ActionReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
	}
	toUserID := reqObj.ToUseID
	actionType := reqObj.ActionType
	content := reqObj.Content

	response, err := service.MessageService.SendMessage(userID, toUserID, actionType, content)
	if err != nil {
		hlog.CtxErrorf(c, "send message error: %v", err) // 记录错误到日志
		ctlFunc.BaseFailedResp(ctx, "send message error")
		return
	}

	ctlFunc.BaseSuccessResp(ctx, response)
}

func (m *message) Chat(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj messageCtlModel.ChatReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	// 参数验证完毕，传入service层处理
	response, err := service.MessageService.GetChatMessages(userID, reqObj.ToUseID, reqObj.PreMsgTime)
	if err != nil {
		hlog.CtxErrorf(c, "get chat messages error: %v", err) // 记录错误到日志
		ctlFunc.BaseFailedResp(ctx, "get chat messages error")
		return
	}

	ctlFunc.Response(ctx, messageCtlModel.ChatResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Messages: response,
	})
}
