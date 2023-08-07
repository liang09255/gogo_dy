package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/commentCtlModel"
	"main/controller/middleware"
	"main/service"
)

type comment struct{}

var Comment = &comment{}

func (e *comment) Action(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	var reqObj commentCtlModel.ActionReq
	if err := ctx.Bind(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	comment, err := service.CommentService.PostCommentAction(c, userID, reqObj)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}

	ctlFunc.Response(ctx, commentCtlModel.ActionResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Comment:  comment,
	})
}

func (e *comment) List(c context.Context, ctx *app.RequestContext) {
	var reqObj commentCtlModel.ListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	data, err := service.CommentService.GetCommentList(c, reqObj.VideoID)
	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}

	ctlFunc.Response(ctx, commentCtlModel.ListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Comments: data,
	})
}
