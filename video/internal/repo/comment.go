package repo

import (
	"context"
	"video/internal/model"
)

type CommentRepo interface {
	AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) // 增加评论
	DeleteComment(ctx context.Context, comment *model.Comment) error                // 删除评论
	GetCommentList(ctx context.Context, videoId int64) ([]model.Comment, error)     // 获得评论列表
	MGetCommentCount(ctx context.Context, videoId []int64) map[int64]int64          // 获取视频的评论数
}
