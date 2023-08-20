package repo

import (
	"context"
	"user/internal/model"
)

type UserRepo interface {
	Exist(ctx context.Context, userName string) bool
	CreateUser(ctx context.Context, user *model.User) error
	CheckUser(ctx context.Context, user *model.User) error
	MGetUserInfo(ctx context.Context, uids []int64) (users []model.User, err error)
}
