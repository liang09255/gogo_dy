package repo

import (
	"context"
	"video/internal/model"
)

type CommentCacheRepo interface {
	SetCommentList(ctx context.Context, vid int64, commentList []model.Comment) error
	GetCommentList(ctx context.Context, vid int64) ([]model.Comment, bool, error)
}
