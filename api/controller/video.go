package controller

import (
	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/userCtlModel"
	"api/controller/ctlModel/videoCtlModel"
	"api/controller/middleware"
	"api/global"

	"common/ggConfig"
	videoRPC "common/ggIDL/video"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"strconv"
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

	msg := videoRPC.FeedRequest{Latest: reqObj.LatestTime, UserId: userID}
	resp, err := global.VideoClient.Feed(c, &msg)
	if err != nil {
		hlog.CtxErrorf(c, "Feed error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Feed error")
		return
	}

	videos := pbVideoList2APIVideoList(resp.VideoList)

	//获取视频列表
	//videos, nextTime, err := service.VideoService.Feed(reqObj.LatestTime, userID)
	//resp. NextTime: 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	if err != nil {
		hlog.CtxErrorf(c, "Feed error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Feed error")
		return
	}

	ctlFunc.Response(ctx, videoCtlModel.FeedResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Videos:   videos,
		//NextTime: nextTime.UnixMicro(),
		NextTime: resp.NextTime,
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
	// 上传视频
	videoUrl, coverUrl, err := uploadVideo(file, userID)
	if err != nil {
		hlog.CtxErrorf(c, "Upload video error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Upload video error")
		return
	}

	msg := videoRPC.UploadVideoRequest{
		Title:    reqObj.Title,
		CoverUrl: coverUrl,
		PlayUrl:  videoUrl,
		UserId:   userID,
	}

	// 调用视频服务
	_, err = global.VideoClient.UploadVideo(c, &msg)
	if err != nil {
		hlog.CtxErrorf(c, "Upload video error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Upload video error")
		return
	}
	//err = service.VideoService.PublishAction(file, reqObj.Title, userID)

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

	msg := videoRPC.VideoListRequest{UserId: reqObj.UserID}
	resp, err := global.VideoClient.VideoList(c, &msg)
	if err != nil {
		hlog.CtxErrorf(c, "Get publish list error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "Get publish list error")
		return
	}

	videos := pbVideoList2APIVideoList(resp.VideoList)
	//videos, err := service.VideoService.GetPublishList(reqObj.UserID)

	ctlFunc.Response(ctx, videoCtlModel.PublishListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp("get publish list success"),
		Videos:   videos,
	})
}

func uploadVideo(file *multipart.FileHeader, userID int64) (videoUrl, coverUrl string, err error) {
	f, err := file.Open()
	if err != nil {
		return
	}
	// 上传视频
	uploadFileKey := strconv.FormatInt(userID, 10) + "/" + uuid.NewV4().String() + ".mp4"
	if err = global.AliOSSBucket.PutObject(uploadFileKey, f); err != nil {
		return
	}
	// 构造URL
	coverFileKey := uploadFileKey + "?x-oss-process=video/snapshot,t_1000,f_jpg,w_720,h_1280,m_fast"
	urlPrefix := "https://" + ggConfig.Config.AliOSS.Bucket + ".oss-cn-shenzhen.aliyuncs.com/"
	videoUrl = urlPrefix + uploadFileKey
	coverUrl = urlPrefix + coverFileKey
	return
}

func pbVideoList2APIVideoList(vs []*videoRPC.Video) []videoCtlModel.Video {
	videos := make([]videoCtlModel.Video, 0, len(vs))
	for _, v := range vs {
		var video = videoCtlModel.Video{
			ID: v.Id,
			Author: userCtlModel.User{
				Id:             v.Author.Id,
				Name:           v.Author.Name,
				FollowCount:    v.Author.FollowCount,
				FollowerCount:  v.Author.FollowerCount,
				IsFollow:       v.Author.IsFollow,
				TotalFavorited: v.Author.TotalFavorited,
				WorkCount:      v.Author.WorkCount,
				FavoriteCount:  v.Author.FavoriteCount,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		}
		videos = append(videos, video)
	}
	return videos
}
