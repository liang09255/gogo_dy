package domain

import (
	"video/internal/dal"
	"video/internal/repo"
)

// 作为一个中间表
type VideoDetailDomain struct {
	tranRepo        repo.TranRepo
	videoDetailRepo repo.VideoDetailRepo
}

func NewVideoDetailDomain() *VideoDetailDomain {
	return &VideoDetailDomain{
		tranRepo:        dal.NewTranRepo(),
		videoDetailRepo: dal.NewVideoDetailDao(),
	}
}
