package repo

import "context"

type RelationCacheRepo interface {
	ADD(ctx context.Context, uid, target int64) error
	MADD(ctx context.Context, uid int64, target []int64) error
	EXISTS(ctx context.Context, uid, target int64) (bool, error)
	MEXISTS(ctx context.Context, uid int64, target []int64) ([]bool, error)
	KeyExist(ctx context.Context, uid int64) bool
	Delete(ctx context.Context, uid int64)
}
