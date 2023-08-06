package dal

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/global"
)

var DB = global.MysqlDB

func Init() {
	err := global.MysqlDB.AutoMigrate(
		&Favorite{},
		&Comment{},
		&User{},
		&Message{})
	if err != nil {
		hlog.Fatalf("auto migrate failed: %v", err)
	}
}
