package service

import (
	"context"
	"main/controller/ctlModel/commentCtlModel"
	"main/controller/ctlModel/userCtlModel"
	"main/dal"
)

type commentService struct{}

var CommentService = &commentService{}

func (h *commentService) PostCommentAction(ctx context.Context, userID, videoID int64, content string, actionType int32) (res commentCtlModel.Comment, err error) {
	var comment = dal.Comment{
		UserId:  userID,
		VideoId: videoID,
		Content: content,
	}
	err = dal.CommentDal.PostCommentAction(ctx, &comment, actionType)
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
			ID:      msg.ID,
			Content: msg.Content,
			User:    userCtlModel.User{},
		}
	}

	return commentListRes
}
