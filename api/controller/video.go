package controller

import (
	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/videoCtlModel"
	"api/controller/middleware"
	"api/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type video struct{}

var Video = &video{}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func (e *video) Feed(c context.Context, ctx *app.RequestContext) {
	var reqObj videoCtlModel.FeedReq
	err := ctx.Bind(&reqObj)
	if err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	var userID int64
	// 单独处理token
	token := ctx.Query("token")
	if token != "" {
		userID, err = middleware.GetUserIDFromToken(token)
		if err != nil {
			ctlFunc.BaseFailedResp(ctx, "token error")
			return
		}
	}

	//获取视频列表
	videos, nextTime, err := service.VideoService.Feed(reqObj.LatestTime, userID)
	//resp. NextTime: 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	if err != nil {
		hlog.CtxErrorf(c, "Feed error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Feed error")
		return
	}

	ctlFunc.Response(ctx, videoCtlModel.FeedResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Videos:   videos,
		NextTime: nextTime.UnixMicro(),
	})
}

// PublishAction 登录用户选择视频上传
func (e *video) PublishAction(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj videoCtlModel.PublishReq
	err := ctx.BindAndValidate(&reqObj)
	if err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}
	// 获取视频文件
	file, err := ctx.FormFile("data")
	if err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}
	// 保存视频
	err = service.VideoService.PublishAction(file, reqObj.Title, userID)
	if err != nil {
		hlog.CtxErrorf(c, "Save video error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Save video error")
		return
	}
	ctlFunc.Response(ctx, videoCtlModel.PublishResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("upload video success"),
	})
	return
}

// PublishList 用户的视频发布列表，直接列出用户所有投稿过的视频
func (e *video) PublishList(c context.Context, ctx *app.RequestContext) {
	var reqObj videoCtlModel.PublishListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	videos, err := service.VideoService.GetPublishList(reqObj.UserID)
	if err != nil {
		hlog.CtxErrorf(c, "Get publish list error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Get publish list error")
		return
	}

	ctlFunc.Response(ctx, videoCtlModel.PublishListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("get publish list success"),
		Videos:   videos,
	})
}
