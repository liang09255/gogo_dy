package controller

import (
	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/relationCtlModel"
	"api/controller/ctlModel/userCtlModel"
	"api/controller/middleware"
	"api/global"
	relationRPC "common/ggIDL/relation"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/jinzhu/copier"
)

type relation struct{}

var Relation = &relation{}

// Action 登录用户对其他用户进行关注或取消关注
func (r *relation) Action(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj relationCtlModel.ActionReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	// 禁止关注自己
	if userID == reqObj.ToUserID {
		ctlFunc.BaseSuccessResp(ctx)
		return
	}

	var msg = &relationRPC.ActionRequest{
		ActionType: reqObj.ActionType,
		MyId:       userID,
		ToUserId:   reqObj.ToUserID,
	}

	ActionResponse, err := global.RelationClient.Action(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "action error")
		return
	}

	hlog.Infof("UserInfoResponse: %v", ActionResponse)

	// 封装返回信息
	var resp = &relationCtlModel.ActionResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
	}

	ctlFunc.Response(ctx, resp)
}

// FollowList 获取关注列表
func (r *relation) FollowList(c context.Context, ctx *app.RequestContext) {
	var reqObj relationCtlModel.FollowListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	userID := middleware.GetUserID(ctx)

	var msg = &relationRPC.FollowListRequest{
		MyId: userID,
	}

	FollowListResponse, err := global.RelationClient.FollowList(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "get followlist error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get followlist error")
		return
	}

	hlog.Infof("FollowListResponse: %v", FollowListResponse)

	var userInfo []userCtlModel.User
	tmp := &userCtlModel.User{}
	for _, user := range FollowListResponse.UserInfo {
		if err := copier.Copy(tmp, user); err != nil {
			hlog.CtxErrorf(c, "get followlist error: %v", err)
			ctlFunc.BaseFailedResp(ctx, "get followlist error")
			return
		}
		userInfo = append(userInfo, *tmp)
	}
	var resp = &relationCtlModel.FollowListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Users:    userInfo,
	}
	ctlFunc.Response(ctx, resp)
}

// FollowerList 获取关注列表
func (r *relation) FollowerList(c context.Context, ctx *app.RequestContext) {
	var reqObj relationCtlModel.FollowerListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}
	userID := middleware.GetUserID(ctx)

	var msg = &relationRPC.FollowerListRequest{
		MyId: userID,
	}

	FollowerListResponse, err := global.RelationClient.FollowerList(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "get followerlist error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get followerlist error")
		return
	}

	hlog.Infof("FollowerListResponse: %v", FollowerListResponse)

	var userInfo []userCtlModel.User
	tmp := &userCtlModel.User{}
	for _, user := range FollowerListResponse.UserInfo {
		if err := copier.Copy(tmp, user); err != nil {
			hlog.CtxErrorf(c, "get followerlist error: %v", err)
			ctlFunc.BaseFailedResp(ctx, "get followerlist error")
			return
		}
		userInfo = append(userInfo, *tmp)
	}
	var resp = &relationCtlModel.FollowerListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Users:    userInfo,
	}
	ctlFunc.Response(ctx, resp)

}

// FriendList 获取朋友列表
func (r *relation) FriendList(c context.Context, ctx *app.RequestContext) {
	var reqObj relationCtlModel.FriendListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	userID := middleware.GetUserID(ctx)

	var msg = &relationRPC.FriendListRequest{
		MyId: userID,
	}

	FriendListResponse, err := global.RelationClient.FriendList(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "get friendlist error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get friendlist error")
		return
	}

	hlog.Infof("FriendListResponse: %v", FriendListResponse)

	var userInfo []userCtlModel.User
	tmp := &userCtlModel.User{}
	for _, user := range FriendListResponse.UserInfo {
		if err := copier.Copy(tmp, user); err != nil {
			hlog.CtxErrorf(c, "get friendlist error: %v", err)
			ctlFunc.BaseFailedResp(ctx, "get friendlist error")
			return
		}
		userInfo = append(userInfo, *tmp)
	}
	var resp = &relationCtlModel.FollowListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Users:    userInfo,
	}
	ctlFunc.Response(ctx, resp)
}
