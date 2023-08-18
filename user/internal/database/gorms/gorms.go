package gorms

import (
	"context"
	"gorm.io/gorm"
)

var _db *gorm.DB

//func init() {
//	//配置MySQL连接参数
//	username := config.C.MC.Username //账号
//	password := config.C.MC.Password //密码
//	host := config.C.MC.Host         //数据库地址，可以是Ip或者域名
//	port := config.C.MC.Port         //数据库端口
//	Dbname := config.C.MC.Db         //数据库名
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
//	var err error
//	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//	if err != nil {
//		panic("连接数据库失败, error=" + err.Error())
//	}
//}

func GetDB() *gorm.DB {
	return _db
}

// 实现了 DbConn 接口
type GormConn struct {
	db *gorm.DB
	tx *gorm.DB
}

func New() *GormConn {
	return &GormConn{db: GetDB()}
}

func NewTran() *GormConn {
	return &GormConn{db: GetDB(), tx: GetDB()}
}

func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) Rollback() {
	g.tx.Rollback()
}

func (g *GormConn) Commit() {
	g.tx.Commit()
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}

func (g *GormConn) Begin() {
	g.tx = GetDB().Begin()
}
