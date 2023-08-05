package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/service"
	"strconv"
)

type favorite struct{}

var Favorite = &favorite{}

func RegFavorite(h *server.Hertz) {
	favoriteGroup := h.Group("/douyin/favorite")
	favoriteGroup.POST("/action", Favorite.Action)
	favoriteGroup.GET("/list", Favorite.List)
}

func (e *favorite) Action(c context.Context, ctx *app.RequestContext) {
	userIdStr, ok := ctx.GetQuery("userId")
	videoIdStr, ok := ctx.GetQuery("videoId")
	actionTypeStr, ok := ctx.GetQuery("action_type")
	if !ok {
		BaseFailResponse(ctx, "userId is required")
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)

	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		BaseFailResponse(ctx, "favoriteAction Error")
		return
	}
	msg, err := service.FavoriteService.PostFavoriteAction(c, userId, videoId, actionType)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		BaseFailResponse(ctx, "favoriteAction Error")
		return
	}
	BaseSuccessResponse(ctx, msg)
}

func (e *favorite) List(c context.Context, ctx *app.RequestContext) {
	userIdStr, ok := ctx.GetQuery("userId")
	videoIdStr, ok := ctx.GetQuery("videoId")
	if !ok {
		BaseFailResponse(ctx, "userId is required")
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		BaseFailResponse(ctx, "favoriteList Error")
		return
	}
	data, err := service.FavoriteService.GetFavoriteList(c, userId, videoId)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		BaseFailResponse(ctx, "favoriteList Error")
		return
	}
	Response(ctx, data)
}
