package repo

import (
	"context"
	"user/internal/data"
	"user/internal/database"
)

type UserRepo interface {
	GetUserCountByUsername(ctx context.Context, userName string) (total int64, err error)
	AddUser(conn database.DbConn, ctx context.Context, mem *data.User) (err error)
}
