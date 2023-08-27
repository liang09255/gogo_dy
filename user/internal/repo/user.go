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
  GetRelation(ctx context.Context, myid, taruserid int64) bool
  TransactionExample(ctx context.Context, conn database.DbConn, otherData string) error
	TransactionExample2(ctx context.Context, conn database.DbConn, otherData string) error
}
