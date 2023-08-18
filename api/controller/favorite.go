package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/favoriteCtlModel"
	"main/controller/middleware"
	"main/service"
)

type favorite struct{}

var Favorite = &favorite{}

func (e *favorite) Action(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj favoriteCtlModel.ActionReq
	err := ctx.Bind(&reqObj)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteAction Error")
		return
	}

	msg, err := service.FavoriteService.PostFavoriteAction(c, userID, reqObj.VideoID, reqObj.ActionType)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteAction Error")
		return
	}

	ctlFunc.Response(ctx, favoriteCtlModel.ActionResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(msg),
	})
}

func (e *favorite) List(c context.Context, ctx *app.RequestContext) {
	var reqObj favoriteCtlModel.ListReq
	err := ctx.Bind(&reqObj)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteList Error")
		return
	}

	videos, err := service.FavoriteService.GetFavoriteList(c, reqObj.UserID)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "favoriteList Error")
		return
	}

	ctlFunc.Response(ctx, favoriteCtlModel.ListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Videos:   videos,
	})
}
