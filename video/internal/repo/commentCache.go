package repo

import (
	"context"
	"video/internal/model"
)

type CommentCacheRepo interface {
	SetCommentList(ctx context.Context, vid int64, commentList []model.Comment) error
	GetCommentList(ctx context.Context, vid int64) ([]model.Comment, bool, error)
	DelCommentList(ctx context.Context, vid int64) error
	DelDelayCommentList(ctx context.Context, vid int64) error
	GetVideoCommentCount(ctx context.Context, vid int64) (int64, bool, error)
	SetVideoCommentCount(ctx context.Context, vid int64, value int64) error
	IncrVideoCommentCount(ctx context.Context, vid int64) error
	DecrVideoCommentCount(ctx context.Context, vid int64) error
}
