package domain

import (
	"common/ggIDL/video"
	"context"
	"video/internal/dal"
	"video/internal/repo"
)

type VideoDomain struct {
	videoRepo repo.VideoRepo
	tranRepo  repo.TranRepo
}

func NewVideoDomain() *VideoDomain {
	return &VideoDomain{
		videoRepo: dal.NewVideoDao(),
		tranRepo:  dal.NewTranRepo(),
	}
}
func (vd *VideoDomain) FavoriteVideoAction(ctx context.Context, videoid int64, actionType video.ActionType) error {
	if actionType == video.ActionType_Add {
		return vd.videoRepo.AddFavoriteCount(ctx, videoid, 1)
	}
	return vd.videoRepo.CancelFavoriteCount(ctx, videoid, 1)
}
