package controller

import (
	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/favoriteCtlModel"
	"api/controller/ctlModel/userCtlModel"
	"api/controller/ctlModel/videoCtlModel"
	"api/controller/middleware"
	"api/global"
	videoRPC "common/ggIDL/video"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

type favorite struct{}

var Favorite = &favorite{}

const (
	Add int32 = 1
	Sub int32 = 2
)

func (e *favorite) Action(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj favoriteCtlModel.ActionReq
	err := ctx.Bind(&reqObj)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteAction error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 参数校验 - action
	var action videoRPC.ActionType

	if reqObj.ActionType == Add {
		action = videoRPC.ActionType_Add
	} else if reqObj.ActionType == Sub {
		action = videoRPC.ActionType_Cancel
	} else {
		hlog.CtxErrorf(c, "参数错误,不支持Action:%d", reqObj.ActionType)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.FavoriteActionUnknown)
		return
	}

	in := &videoRPC.FavoriteActionRequest{
		VideoId:    reqObj.VideoID,
		ActionType: action,
		UserId:     userID,
	}

	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	_, err = global.VideoClient.FavoriteAction(c, in)
	if err != nil {
		hlog.CtxErrorf(c, "video set favorite action error : %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.ServerInternal.WithDetails(err.Error()))
		return
	}

	ctlFunc.Response(ctx, favoriteCtlModel.ActionResp{
		APIBaseResp: baseCtlModel.NewBaseSuccessResp("操作成功"),
	})
}

func (e *favorite) List(c context.Context, ctx *app.RequestContext) {
	var reqObj favoriteCtlModel.ListReq
	err := ctx.Bind(&reqObj)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.InvalidParams)
		return
	}

	in := &videoRPC.FavoriteListRequest{
		UserId: reqObj.UserID,
	}

	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	videoListResp, err := global.VideoClient.FavoriteList(c, in)
	if err != nil {
		hlog.CtxErrorf(c, "favoriteList error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.ServerInternal.WithDetails(err.Error()))
		return
	}

	videos := make([]videoCtlModel.Video, 0)

	for _, v := range videoListResp.VideoList {
		video := videoCtlModel.Video{
			ID:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
			Author: userCtlModel.User{
				Id:             v.Author.Id,
				Name:           v.Author.Name,
				FollowerCount:  v.Author.FollowerCount,
				FollowCount:    v.Author.FollowCount,
				IsFollow:       v.Author.IsFollow,
				TotalFavorited: v.Author.TotalFavorited,
				WorkCount:      v.Author.WorkCount,
				FavoriteCount:  v.Author.FavoriteCount,
			},
		}
		videos = append(videos, video)

	}

	ctlFunc.Response(ctx, favoriteCtlModel.ListResp{
		APIBaseResp: baseCtlModel.NewBaseSuccessResp(),
		Videos:      videos,
	})
}
