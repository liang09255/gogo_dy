package service

import (
	"context"
	"fmt"
	"main/controller/ctlModel/commentCtlModel"
	"main/controller/ctlModel/userCtlModel"
	"main/dal"
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
	// TODO 获取用户信息
	res.User = userCtlModel.User{}
	res.Content = comment.Content
	res.ID = comment.ID
	res.CreateDate = comment.CreatedAt.Format("2006-01-02 15:04:05")
	return res, nil
}

func (h *commentService) delComment(ctx context.Context, c commentCtlModel.ActionReq) (res commentCtlModel.Comment, err error) {
	var comment = dal.Comment{
		ID: c.CommentID,
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

	for i, msg := range commentList {
		// TODO 获取用户信息
		commentListRes[i] = commentCtlModel.Comment{
			ID:         msg.ID,
			Content:    msg.Content,
			User:       userCtlModel.User{},
			CreateDate: msg.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return commentListRes
}
