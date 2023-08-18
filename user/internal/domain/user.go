package domain

import (
	"context"
	"user/internal/dao"
	"user/internal/data"
	"user/internal/database"
	"user/internal/repo"
)

type UserDomain struct {
	userRepo repo.UserRepo
}

func NewUserDomain() *UserDomain {
	return &UserDomain{
		userRepo: dao.NewUserDao(),
	}
}

func (ud *UserDomain) UserIsNotExist(ctx context.Context, userName string) (bool, error) {
	// 这里不需要对参数进行处理，所以可以直接调用repo层的接口
	count, err := ud.userRepo.GetUserCountByUsername(ctx, userName)
	// 对dao层返回的数据进行处理、判断，再返回给service层想要的数据
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (ud *UserDomain) AddUser(conn database.DbConn, ctx context.Context, userName string, passWord string) error {
	user := &data.User{
		UserName: userName,
		Password: passWord,
	}
	err := ud.userRepo.AddUser(conn, ctx, user)
	if err != nil {
		return err
	}
	return nil
}
