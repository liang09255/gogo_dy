package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/dal"
	"main/service"
	"net/http"
	"os"
	"strings"
	"time"
)

type vedio struct{}

var Video = &vedio{}

func String(value string) *string {
	return &value
}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func (e *vedio) Feed(c context.Context, ctx *app.RequestContext) {

	var req service.DouyinFeedRequest
	var resp service.DouyinFeedResponse

	req.Token = String(ctx.Query("token"))
	req.LatestTime = String(ctx.Query("latest_time"))
	//可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	latest := time.Now().Unix()

	//获取视频列表
	err := service.VideoService.Feed(latest, *req.Token, &resp.VideoList)

	//resp. NextTime: 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	if err != nil {
		hlog.CtxErrorf(c, "Feed error: %v", err)
		resp.StatusCode = FailCode
		resp.StatusMsg = String("Get feed fail")
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.StatusCode = SuccessCode
	resp.StatusMsg = String("Get feed successfully")

	ctx.JSON(http.StatusOK, resp)
}

// PublishAction 登录用户选择视频上传
func (e *vedio) PublishAction(c context.Context, ctx *app.RequestContext) {

	Title := ctx.FormValue("title")
	Token := ctx.FormValue("token")
	var resp service.DouyinPublishActionResponse

	//获取用户上传的视频文件
	file, _ := ctx.FormFile("data")
	fmt.Println(file.Filename)

	err := ctx.SaveUploadedFile(file, fmt.Sprintf("./videofile/%s", file.Filename))

	fmt.Println("上传文件成功\n")

	if err != nil {
		println("视频发布错误\n")
		if err != nil {
			resp.StatusCode = FailCode
			resp.StatusMsg = String("publish video fail")
			ctx.JSON(http.StatusOK, resp)
		}
	}

	//将视频信息存入数据库中
	var video dal.Video
	video.Title = string(Title)
	video.FavoriteCount = 0
	video.CommentCount = 0
	video.PlayUrl = "./videofile/" + file.Filename
	video.CreateTime = time.Now().Unix()

	err = service.VideoService.PublishAction(video, string(Token))

	if err != nil {
		hlog.CtxErrorf(c, "Publish error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// Publishlist 用户的视频发布列表，直接列出用户所有投稿过的视频
func (e *vedio) Publishlist(c context.Context, ctx *app.RequestContext) {

	var req service.DouyinPublishListRequest
	var resp service.DouyinPublishListResponse
	var ok bool
	req.Token, ok = ctx.GetQuery("token")

	if !ok {
		BaseFailResponse(ctx, "token is required")
		return
	}

	req.UserID, ok = ctx.GetQuery("user_id")

	if !ok {
		BaseFailResponse(ctx, "user_id is required")
		return
	}

	err := service.RelationService.GetPublishList(req.Token, req.UserID, &resp.VideoList)

	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = String("Action Fail")
		hlog.CtxErrorf(c, "Relation action error: %v", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.StatusCode = 0
	resp.StatusMsg = String("Action successfully")
	ctx.JSON(http.StatusOK, resp)
}

// PlayVideo 播放视频
func (e *vedio) PlayVideo(c context.Context, ctx *app.RequestContext) {
	var ok bool
	fileName, ok := ctx.GetQuery("play_url")

	if !ok {
		BaseFailResponse(ctx, "play_url is required")
		return
	}

	content, err := os.ReadFile(fileName)

	if err != nil {
		hlog.CtxErrorf(c, "Play video error: %v", err)
		return
	}

	//获取文件名后缀
	words := strings.Split(fileName, ".")
	suffix := words[len(words)-1]
	//文件为视频
	if suffix == "mp4" {
		ctx.Data(http.StatusOK, "video/mp4", content)
	}
	//文件是图片
	if suffix == "jpg" {
		ctx.Data(http.StatusOK, "image/png", content)
	}
}
