package dal

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/global"
)

func Init() {
	err := global.MysqlDB.AutoMigrate(&Hello{}, &User{}, &Message{})
	if err != nil {
		hlog.Fatalf("auto migrate failed: %v", err)
	}
}
