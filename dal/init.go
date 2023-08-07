package dal

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/global"
)

func Init() {
	err := global.MysqlDB.AutoMigrate(
		&Favorite{},
		&Comment{},
		&User{},
		&Message{},
		&Video{},
		&Relation{},
	)
	if err != nil {
		hlog.Fatalf("auto migrate failed: %v", err)
	}
}
