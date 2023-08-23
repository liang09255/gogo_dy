package repo

import (
	"context"
	"time"

	"user/internal/model"
)

type UserCacheRepo interface {
	MGetUserInfo(ctx context.Context, uids []int64) (users []model.User, err error)
	MSetUserInfo(ctx context.Context, users []model.User, expire time.Duration) error
}
