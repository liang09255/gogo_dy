package repo

import (
	"context"
)

type HotCacheRepo interface {
	AddVidsToHotCache(ctx context.Context, uid int64, vids []int64) error
	GetAndResetHotCache(ctx context.Context) ([]int64, error)
}
