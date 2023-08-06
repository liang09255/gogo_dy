package controller

import (
	"context"
	"main/controller/ctlFunc"
	"main/service"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type favorite struct{}

var Favorite = &favorite{}

func (e *favorite) Action(c context.Context, ctx *app.RequestContext) {
	userIdStr, ok := ctx.GetQuery("userId")
	videoIdStr, ok := ctx.GetQuery("videoId")
	actionTypeStr, ok := ctx.GetQuery("action_type")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "userId is required")
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)

	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteAction Error")
		return
	}
	msg, err := service.FavoriteService.PostFavoriteAction(c, userId, videoId, actionType)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteAction Error")
		return
	}
	ctlFunc.BaseSuccessResp(ctx, msg)
}

func (e *favorite) List(c context.Context, ctx *app.RequestContext) {
	userIdStr, ok := ctx.GetQuery("userId")
	videoIdStr, ok := ctx.GetQuery("videoId")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "userId is required")
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteList Error")
		return
	}
	data, err := service.FavoriteService.GetFavoriteList(c, userId, videoId)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteList Error")
		return
	}
	ctlFunc.Response(ctx, data)
}
