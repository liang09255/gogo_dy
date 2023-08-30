package controller

import (
	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/commentCtlModel"
	"api/controller/ctlModel/userCtlModel"
	"api/controller/middleware"
	"api/global"
	videoRPC "common/ggIDL/video"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
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

	var action videoRPC.ActionType

	if reqObj.ActionType == Add {
		action = videoRPC.ActionType_Add
	} else if reqObj.ActionType == Sub {
		action = videoRPC.ActionType_Cancel
	} else {
		hlog.CtxErrorf(c, "参数错误，不支持Action:%d", reqObj.ActionType)
		ctlFunc.BaseFailedResp(ctx, "invalid params,commentAction Error")
		return
	}

	in := &videoRPC.CommentActionRequest{
		VideoId:     reqObj.VideoID,
		ActionType:  action,
		CommentId:   reqObj.CommentID,
		UserId:      userID,
		CommentText: reqObj.CommentText,
	}

	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	comment, err := global.VideoClient.CommentAction(c, in)

	if err != nil {
		hlog.CtxErrorf(c, "comment error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment Error")
		return
	}

	commentResp := commentCtlModel.Comment{
		ID:         comment.Comment.Id,
		Content:    comment.Comment.Content,
		CreateDate: comment.Comment.CreateDate,
		User: userCtlModel.User{
			Id:             comment.Comment.User.Id,
			Name:           comment.Comment.User.Name,
			FollowerCount:  comment.Comment.User.FollowerCount,
			FollowCount:    comment.Comment.User.FollowCount,
			IsFollow:       comment.Comment.User.IsFollow,
			TotalFavorited: comment.Comment.User.TotalFavorited,
			WorkCount:      comment.Comment.User.WorkCount,
			FavoriteCount:  comment.Comment.User.FavoriteCount,
		},
	}

	ctlFunc.Response(ctx, commentCtlModel.ActionResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Comment:  commentResp,
	})
}

func (e *comment) List(c context.Context, ctx *app.RequestContext) {
	var reqObj commentCtlModel.ListReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}

	in := &videoRPC.CommentListRequest{
		VideoId: reqObj.VideoID,
	}
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	commentListResp, err := global.VideoClient.CommentList(c, in)
	if err != nil {
		hlog.CtxErrorf(c, "comment list error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "comment list Error")
		return
	}

	// 构造返回结构
	comments := make([]commentCtlModel.Comment, 0)
	for _, v := range commentListResp.CommentList {
		comment := commentCtlModel.Comment{
			ID:         v.Id,
			Content:    v.Content,
			CreateDate: v.CreateDate,
			User: userCtlModel.User{
				Id:             v.User.Id,
				Name:           v.User.Name,
				FollowerCount:  v.User.FollowerCount,
				IsFollow:       v.User.IsFollow,
				TotalFavorited: v.User.TotalFavorited,
				WorkCount:      v.User.WorkCount,
				FavoriteCount:  v.User.FavoriteCount,
			},
		}
		comments = append(comments, comment)
	}

	ctlFunc.Response(ctx, commentCtlModel.ListResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		Comments: comments,
	})
}
