package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/model"
	"main/service"
	"net/http"
)

type relation struct{}

var Relation = &relation{}

func RegRelation(h *server.Hertz) {
	h.POST("/douyin/relation/action/", Relation.RelationAction)
	h.GET("/douyin/relation/follow/list/", Relation.GetFollowList)
	h.GET("/douyin/relation/follower/list/", Relation.GetFollowerList)
	h.GET("/douyin/relation/friend/list/", Relation.GetFriendList)
}

// RelationAction 登录用户对其他用户进行关注或取消关注
func (r *relation) RelationAction(c context.Context, ctx *app.RequestContext) {
	//登录用户对其他用户进行关注或取消关注。
	var req model.DouyinRelationActionRequest
	var resp model.DouyinRelationActionResponse

	req.ActionType = ctx.Query("action_type")
	req.Token = ctx.Query("token")
	req.ToUserID = ctx.Query("to_user_id")

	fmt.Println(req)

	err := service.RelationService.RelationAction(req.Token, req.ToUserID, req.ActionType, &resp)

	if err != nil {
		resp.StatusCode = FailCode
		resp.StatusMsg = "Action Fail"
		hlog.CtxErrorf(c, "Relation action error: %v", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetFollowList 获取关注列表
func (r *relation) GetFollowList(c context.Context, ctx *app.RequestContext) {
	//登录用户关注的所有用户列表。
	var req model.DouyinRelationFollowListRequest
	var resp model.DouyinRelationFollowListResponse

	req.UserID = ctx.Query("user_id")
	req.Token = ctx.Query("token")

	fmt.Println(req)

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

// GetFollowerList 获取关注列表
func (r *relation) GetFollowerList(c context.Context, ctx *app.RequestContext) {
	//所有关注登录用户的粉丝列表。
	var req model.DouyinRelationFollowerListRequest
	var resp model.DouyinRelationFollowerListResponse

	req.UserID = ctx.Query("user_id")
	req.Token = ctx.Query("token")

	fmt.Println(req)

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

// GetFriendList 获取朋友列表
func (r *relation) GetFriendList(c context.Context, ctx *app.RequestContext) {
	//所有关注登录用户的粉丝列表。
	var req model.DouyinRelationFriendListRequest
	var resp model.DouyinRelationFriendListResponse

	req.UserID = ctx.Query("user_id")
	req.Token = ctx.Query("token")

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
