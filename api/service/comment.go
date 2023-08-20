package service

import (
	"api/controller/ctlModel/commentCtlModel"
	"api/controller/ctlModel/userCtlModel"
	"api/dal"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type commentService struct{}

var CommentService = &commentService{}

const (
	PostCommentCode = 1
	DelCommentCode  = 2
)

func (h *commentService) PostCommentAction(ctx context.Context, userID int64, c commentCtlModel.ActionReq) (res commentCtlModel.Comment, err error) {
	if c.ActionType == PostCommentCode {
		return h.postComment(ctx, userID, c)
	} else if c.ActionType == DelCommentCode {
		return h.delComment(ctx, c)
	}
	return commentCtlModel.Comment{}, fmt.Errorf("comment action type error, action_type: %d", c.ActionType)
}

func (h *commentService) postComment(ctx context.Context, userID int64, c commentCtlModel.ActionReq) (res commentCtlModel.Comment, err error) {
	var comment = dal.Comment{
		UserId:  userID,
		VideoId: c.VideoID,
		Content: c.CommentText,
	}
	err = dal.CommentDal.PostComment(ctx, &comment)
	if err != nil {
		return commentCtlModel.Comment{}, err
	}

	if err := VideoService.AddCommentCount(c.VideoID); err != nil {
		hlog.CtxErrorf(ctx, "add comment count error, err: %v", err)
	}

	userInfos, err := UserService.MGetUserInfo([]int64{userID}, userID)
	if err != nil {
		hlog.CtxErrorf(ctx, "get user info error, err: %v", err)
	} else {
		res.User = userInfos[0]
	}
	res.Content = comment.Content
	res.ID = comment.ID
	res.CreateDate = comment.CreatedAt.Format("2006-01-02 15:04:05")
	return res, nil
}

func (h *commentService) delComment(ctx context.Context, c commentCtlModel.ActionReq) (res commentCtlModel.Comment, err error) {
	var comment = dal.Comment{
		ID: c.CommentID,
	}

	if err := VideoService.ReduceCommentCount(c.VideoID); err != nil {
		hlog.CtxErrorf(ctx, "reduce comment count error, err: %v", err)
	}

	return commentCtlModel.Comment{}, dal.CommentDal.DelComment(ctx, &comment)
}

func (h *commentService) GetCommentList(ctx context.Context, videoId int64) ([]commentCtlModel.Comment, error) {
	arr, err := dal.CommentDal.GetCommentList(ctx, videoId)
	if err != nil {
		return nil, err
	}
	return convertCommentListToCommentListResponse(arr), nil
}

func convertCommentListToCommentListResponse(commentList []*dal.Comment) []commentCtlModel.Comment {

	commentListRes := make([]commentCtlModel.Comment, len(commentList))

	uids := make([]int64, len(commentList))
	for _, c := range commentList {
		uids = append(uids, c.UserId)
	}
	userInfos, err := UserService.MGetUserInfosMap(uids)
	if err != nil {
		hlog.Errorf("get user info error, err: %v", err)
		userInfos = make(map[int64]userCtlModel.User)
	}

	for i, msg := range commentList {
		commentListRes[i] = commentCtlModel.Comment{
			ID:         msg.ID,
			Content:    msg.Content,
			User:       userInfos[msg.UserId],
			CreateDate: msg.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return commentListRes
}
