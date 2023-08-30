package dal

import (
	"context"
	"fmt"
	"user/internal/database"

	"user/internal/model"
	"user/internal/repo"
)

type UserDal struct {
	conn *database.GormConn
}

const MaxUsernameLength = 32

var _ repo.UserRepo = (*UserDal)(nil)

func NewUserDao() *UserDal {
	return &UserDal{
		conn: database.New(),
	}
}

func (ud *UserDal) Exist(ctx context.Context, userName string) bool {
	user := new(model.User)
	ud.conn.WithContext(ctx).Where("username = ?", userName).First(&user)
	return user.ID != 0
}

func (ud *UserDal) CreateUser(ctx context.Context, user *model.User) error {
	if len(user.Username) > MaxUsernameLength {
		return fmt.Errorf("用户名超出最大长度限制")
	}
	return ud.conn.WithContext(ctx).Create(user).First(user).Error
}

func (ud *UserDal) CheckUser(ctx context.Context, user *model.User) error {
	t := ud.conn.WithContext(ctx).Where("username = ? and password = ?", user.Username, user.Password).Find(user)
	//id为零值，说明sql执行失败
	if user.ID == 0 {
		return fmt.Errorf("用户验证失败")
	}
	return t.Error
}

func (ud *UserDal) MGetUserInfo(ctx context.Context, uids []int64) (users []model.User, err error) {

	t := ud.conn.WithContext(ctx).Where("id in ?", uids).Find(&users)

	return users, t.Error
}
