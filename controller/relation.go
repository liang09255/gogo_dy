package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/relationCtlModel"
	"main/controller/middleware"
	"main/service"
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

	err := service.RelationService.RelationAction(userID, reqObj.ToUserID, reqObj.ActionType)
	if err != nil {
		hlog.CtxErrorf(c, "relation action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "relation action Failed")
		return
	}
	ctlFunc.Response(ctx, relationCtlModel.ActionResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("follow action success"),
	})
}

// FollowList 获取关注列表
func (r *relation) FollowList(c context.Context, ctx *app.RequestContext) {
	var reqObj relationCtlModel.FollowListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	users, err := service.RelationService.GetFollowList(reqObj.UserID)
	if err != nil {
		hlog.CtxErrorf(c, "relation action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get follow list Failed")
		return
	}

	ctlFunc.Response(ctx, relationCtlModel.FollowListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("get follow list success"),
		Users:    users,
	})
}

// FollowerList 获取关注列表
func (r *relation) FollowerList(c context.Context, ctx *app.RequestContext) {
	var reqObj relationCtlModel.FollowerListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	users, err := service.RelationService.GetFollowerList(reqObj.UserID)
	if err != nil {
		hlog.CtxErrorf(c, "relation action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get follower list Failed")
		return
	}

	ctlFunc.Response(ctx, relationCtlModel.FollowerListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("get follower list success"),
		Users:    users,
	})
}

// FriendList 获取朋友列表
func (r *relation) FriendList(c context.Context, ctx *app.RequestContext) {
	var reqObj relationCtlModel.FriendListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	users, err := service.RelationService.GetFriendList(reqObj.UserID)
	if err != nil {
		hlog.CtxErrorf(c, "relation action error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get friend list Failed")
		return
	}

	ctlFunc.Response(ctx, relationCtlModel.FriendListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("get friend list success"),
		Users:    users,
	})
}
