package repo

import (
	"context"
	"video/internal/model"
)

type CommentRepo interface {
	// 增加评论
	AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	// 删除评论
	DeleteComment(ctx context.Context, comment *model.Comment) error
	// 获得评论列表
	GetCommentList(ctx context.Context, videoId int64) ([]model.Comment, error)
}
