package controller

import (
	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/messageCtlModel"
	"api/controller/middleware"
	"api/global"
	chatRpc "common/ggIDL/chat"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/jinzhu/copier"
	"time"
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

	msg := &chatRpc.ChatActionRequest{
		FromUserId: userID,
		ActionType: chatRpc.ActionType(reqObj.ActionType),
		Content:    reqObj.Content,
		ToUserId:   reqObj.ToUseID,
	}
	var err error
	if err != nil {
		hlog.CtxErrorf(c, "message action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "message action error")
		return
	}
	c, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()
	ActionResponse, err := global.ChatClient.Action(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "message action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "message action error")
		return
	}
	msgResponse := &messageCtlModel.ActionResp{}
	if err := copier.Copy(msgResponse, ActionResponse); err != nil {
		hlog.CtxErrorf(c, "message action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "message action error")
		return
	}
	if err != nil {
		hlog.CtxErrorf(c, "send message error: %v", err) // 记录错误到日志
		ctlFunc.BaseFailedResp(ctx, "send message error")
		return
	}

	var response = messageCtlModel.ActionResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
	}
	ctlFunc.Response(ctx, &response)
}

func (m *message) Chat(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj messageCtlModel.ChatReq
	if err := ctx.Bind(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	var msg = &chatRpc.ListRequest{
		FromUserId: userID,
		ToUserId:   reqObj.ToUseID,
		PreMsgTime: reqObj.PreMsgTime,
	}
	MessageInfoResponse, err := global.ChatClient.List(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "get chat messages error: %v", err) // 记录错误到日志
		ctlFunc.BaseFailedResp(ctx, "get chat messages error")
		return
	}
	var messageInfo []messageCtlModel.Message
	for _, m := range MessageInfoResponse.List {
		messageInfo = append(messageInfo, messageCtlModel.Message{
			ID:         m.Id,
			Content:    m.Content,
			CreateTime: m.CreateTime.AsTime().UnixMilli(),
			FromUserID: m.FromUserId,
			ToUserID:   m.ToUserId,
		})
	}
	var resp = &messageCtlModel.ChatResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Messages: messageInfo,
	}
	ctlFunc.Response(ctx, resp)
}
