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
	// FixMe 喜欢操作需要鉴权
	userIdStr, ok := ctx.GetQuery("user_id")
	videoIdStr, ok := ctx.GetQuery("video_id")
	actionTypeStr, ok := ctx.GetQuery("action_type")
	if !ok {
		BaseFailResponse(ctx, "user_id is required")
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
	userIdStr, ok := ctx.GetQuery("user_id")
	if !ok {
		BaseFailResponse(ctx, "user_id is required")
		return
	}
	videoIdStr, ok := ctx.GetQuery("video_id")

	if !ok {
		BaseFailResponse(ctx, "video_id is required")
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
	ResponseWithData(ctx, "", data)
}
