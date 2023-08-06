package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/service"
	"strconv"
)

type favorite struct{}

var Favorite = &favorite{}

func (e *favorite) Action(c context.Context, ctx *app.RequestContext) {
	// FixMe 喜欢操作需要鉴权
	userIdStr, ok := ctx.GetQuery("user_id")
	videoIdStr, ok := ctx.GetQuery("video_id")
	actionTypeStr, ok := ctx.GetQuery("action_type")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "user_id is required")
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
	userIdStr, ok := ctx.GetQuery("user_id")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "user_id is required")
		return
	}
	videoIdStr, ok := ctx.GetQuery("video_id")

	if !ok {
		ctlFunc.BaseFailedResp(ctx, "video_id is required")
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
	ctlFunc.ResponseWithData(ctx, "", data)
}
