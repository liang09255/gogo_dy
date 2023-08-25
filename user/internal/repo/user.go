package repo

import (
	"context"
	"user/internal/database"
	"user/internal/model"
)

type UserRepo interface {
	Exist(ctx context.Context, userName string) bool
	CreateUser(ctx context.Context, user *model.User) error
	CheckUser(ctx context.Context, user *model.User) error
	MGetUserInfo(ctx context.Context, uids []int64) (users []model.User, err error)
	TransactionExample(ctx context.Context, conn database.DbConn, otherData string) error
	TransactionExample2(ctx context.Context, conn database.DbConn, otherData string) error
	// 点赞相关
	AddTotalFavorited(ctx context.Context, userid int64, count int64) error
	SubTotalFavorited(ctx context.Context, userid int64, count int64) error
	AddFavoriteCount(ctx context.Context, userid int64, count int64) error
	SubFavoriteCount(ctx context.Context, userid int64, count int64) error
}
