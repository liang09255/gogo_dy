package controller

import (
	"context"
	"main/controller/ctlFunc"
	"main/service"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type relation struct{}

var Relation = &relation{}

// Action 登录用户对其他用户进行关注或取消关注
func (r *relation) Action(c context.Context, ctx *app.RequestContext) {
	//登录用户对其他用户进行关注或取消关注。
	var req service.DouyinRelationActionRequest
	var resp service.DouyinRelationActionResponse

	req.ActionType = ctx.Query("action_type")
	var ok bool
	req.Token, ok = ctx.GetQuery("token")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "token is required")
		return
	}

	req.ToUserID, ok = ctx.GetQuery("to_user_id")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "to_user_id is required")
		return
	}

	err := service.RelationService.RelationAction(req.Token, req.ToUserID, req.ActionType, &resp)

	if err != nil {
		resp.StatusCode = ctlFunc.FailCode
		resp.StatusMsg = "Action Fail"
		hlog.CtxErrorf(c, "Relation action error: %v", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// FollowList 获取关注列表
func (r *relation) FollowList(c context.Context, ctx *app.RequestContext) {
	//登录用户关注的所有用户列表。
	var req service.DouyinRelationFollowListRequest
	var resp service.DouyinRelationFollowListResponse
	var ok bool
	req.UserID, ok = ctx.GetQuery("user_id")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "user_id is required")
		return
	}

	req.Token, ok = ctx.GetQuery("token")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "token is required")
		return
	}

	err := service.RelationService.GetFollowList(req.Token, req.UserID, &resp.UserList)

	if err != nil {
		resp.StatusCode = "1"
		resp.StatusMsg = String("Get List Fail")
		hlog.CtxErrorf(c, "Relation action error: %v", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.StatusCode = "0"
	resp.StatusMsg = String("Get List successfully")
	ctx.JSON(http.StatusOK, resp)
}

// FollowerList 获取关注列表
func (r *relation) FollowerList(c context.Context, ctx *app.RequestContext) {
	//所有关注登录用户的粉丝列表。
	var req service.DouyinRelationFollowerListRequest
	var resp service.DouyinRelationFollowerListResponse
	var ok bool
	req.UserID, ok = ctx.GetQuery("user_id")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "user_id is required")
		return
	}

	req.Token, ok = ctx.GetQuery("token")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "token is required")
		return
	}

	err := service.RelationService.GetFollowerList(req.Token, req.UserID, &resp.UserList)

	if err != nil {
		resp.StatusCode = "1"
		resp.StatusMsg = String("Get List Fail")
		hlog.CtxErrorf(c, "Relation action error: %v", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.StatusCode = "0"
	resp.StatusMsg = String("Get List successfully")
	ctx.JSON(http.StatusOK, resp)
}

// FriendList 获取朋友列表
func (r *relation) FriendList(c context.Context, ctx *app.RequestContext) {
	//所有关注登录用户的粉丝列表。
	var req service.DouyinRelationFriendListRequest
	var resp service.DouyinRelationFriendListResponse
	var ok bool

	req.UserID, ok = ctx.GetQuery("user_id")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "to_user_id is required")
		return
	}

	req.Token, ok = ctx.GetQuery("token")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "to_user_id is required")
		return
	}

	err := service.RelationService.GetFriendList(req.Token, req.UserID, &resp.UserList)

	if err != nil {
		resp.StatusCode = "1"
		resp.StatusMsg = String("Action Fail")
		hlog.CtxErrorf(c, "Relation action error: %v", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.StatusCode = "0"
	resp.StatusMsg = String("Action successfully")
	ctx.JSON(http.StatusOK, resp)
}
