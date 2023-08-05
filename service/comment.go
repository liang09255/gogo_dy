package service

import (
	"context"
	"main/dal"
)

type commentService struct{}
type CommentResponse struct {
	ID      int64  `json:"id"`
	UserId  int64  `json:"user_id"`
	Content string `json:"content"`
}

var CommentService = &commentService{}

func (h *commentService) GetCommentList(ctx context.Context, userId int64, videoId int64) ([]*CommentResponse, error) {
	arr, err := dal.CommentDal.GetCommentList(ctx, userId, videoId)
	if err != nil {
		return nil, err
	}
	return convertCommentListToCommentListResponse(arr), nil
}
func (h *commentService) PostCommentAction(ctx context.Context, comment *dal.Comment) error {
	return dal.CommentDal.PostCommentAction(ctx, comment)
}
func convertCommentListToCommentListResponse(commentList []*dal.Comment) []*CommentResponse {

	commentListRes := make([]*CommentResponse, len(commentList))

	for i, msg := range commentList {
		commentListRes[i] = &CommentResponse{
			ID:      msg.ID,
			Content: msg.Content,
			UserId:  msg.UserId,
		}
	}

	return commentListRes
}
