package database

import (
	"common/ggDB"
	"context"
	"gorm.io/gorm"
)

var _db *gorm.DB

func init() {
	_db = ggDB.NewMySQL()
}

func GetDB() *gorm.DB {
	return _db
}

// GormConn 实现了 DbConn 接口
type GormConn struct {
	db *gorm.DB
}

func New() *GormConn {
	return &GormConn{db: GetDB()}
}

func (g *GormConn) Rollback() {
	g.db.Rollback()
}

func (g *GormConn) Commit() {
	g.db.Commit()
}

func (g *GormConn) WithContext(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx)
}

func (g *GormConn) Begin() {
	g.db = GetDB().Begin()
}
