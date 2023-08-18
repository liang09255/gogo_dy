package dao

import (
	"context"
	"user/internal/data"
	"user/internal/database"
	"user/internal/database/gorms"
)

type UserDao struct {
	conn *gorms.GormConn
}

func NewUserDao() *UserDao {
	return &UserDao{
		conn: gorms.New(),
	}
}

func (ud *UserDao) GetUserCountByUsername(ctx context.Context, userName string) (count int64, err error) {
	// dao层必须只含有数据库相关操作，其他数据处理操作放在domain层
	err = ud.conn.Session(ctx).Model(&data.User{}).Where("username = ?", userName).Count(&count).Error
	return
}

func (ud *UserDao) AddUser(conn database.DbConn, ctx context.Context, mem *data.User) (err error) {
	ud.conn = conn.(*gorms.GormConn)
	err = ud.conn.Tx(ctx).Create(mem).Error
	return
}
