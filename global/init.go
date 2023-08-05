package global

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB

func Init() {
	var err error

	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err = v.ReadInConfig(); err != nil {
		hlog.CtxFatalf(nil, "read config failed: %v", err)
	}
	if err = v.Unmarshal(Config); err != nil {
		hlog.CtxFatalf(nil, "init config failed: %v", err)
	}
	InitMysqlDb(Config.Mysql.Dsn)
}
func InitMysqlDb(dsn string) {
	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		hlog.Fatalf("database initialize failed", err)
		return
	}
	hlog.Fatalf("数据库连接成功!")
}
