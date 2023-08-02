package service

import (
	"context"
	"main/dal"
)

type commentService struct{}

var CommentService = &commentService{}

func (h *commentService) GetCommentList(ctx context.Context, userId int64, videoId int64) error {
	dal.GetCommentList(ctx, videoId)
	//if err != nil {
	//	return nil, err
	//}
	return nil
}
func (h *commentService) PostCommentAction(ctx context.Context, comment *dal.Comment) error {
	dal.PostCommentAction(ctx, comment)
	//if err != nil {
	//	return nil, err
	//}
	return nil
}
