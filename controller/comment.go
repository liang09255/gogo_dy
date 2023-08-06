package controller

import (
	"context"
	"main/controller/ctlFunc"
	"main/dal"
	"main/service"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type comment struct{}

var Comment = &comment{}

func (e *comment) Action(c context.Context, ctx *app.RequestContext) {
	var comment dal.Comment
	err := ctx.Bind(&comment)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}
	err = service.CommentService.PostCommentAction(c, &comment)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}
	ctlFunc.BaseSuccessResp(ctx, "评论成功")
}

func (e *comment) List(c context.Context, ctx *app.RequestContext) {
	userIdStr, ok := ctx.GetQuery("userId")
	videoIdStr, ok := ctx.GetQuery("videoId")
	if !ok {
		ctlFunc.BaseFailedResp(ctx, "userId is required")
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}
	data, err := service.CommentService.GetCommentList(c, userId, videoId)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}
	ctlFunc.Response(ctx, data)
}
