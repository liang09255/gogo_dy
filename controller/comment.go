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
	// FixMe 评论操作需要鉴权
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
	userIdStr, ok := ctx.GetQuery("user_id")
	videoIdStr, ok := ctx.GetQuery("video_id")
	if !ok {
		BaseFailResponse(ctx, "user_id is required")
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
	ResponseWithData(ctx, "", data)
}
