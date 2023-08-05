package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/dal"
	"main/service"
	"strconv"
)

type comment struct{}

var Comment = &comment{}

func RegComment(h *server.Hertz) {

	comment := h.Group("/douyin/comment")
	comment.POST("/action", Comment.Action)
	comment.GET("/list", Comment.List)
}

func (e *comment) Action(c context.Context, ctx *app.RequestContext) {
	var comment dal.Comment
	err := ctx.Bind(&comment)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		BaseFailResponse(ctx, "comment Error")
		return
	}
	err = service.CommentService.PostCommentAction(c, &comment)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		BaseFailResponse(ctx, "comment Error")
		return
	}
	BaseSuccessResponse(ctx, "评论成功")
}

func (e *comment) List(c context.Context, ctx *app.RequestContext) {
	userIdStr, ok := ctx.GetQuery("userId")
	videoIdStr, ok := ctx.GetQuery("videoId")
	if !ok {
		BaseFailResponse(ctx, "userId is required")
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		BaseFailResponse(ctx, "comment Error")
		return
	}
	data, err := service.CommentService.GetCommentList(c, userId, videoId)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		BaseFailResponse(ctx, "comment Error")
		return
	}
	Response(ctx, data)
}
