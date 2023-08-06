package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB
var AliOSSBucket *oss.Bucket

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
	InitAliOSSBucket()
}

func InitMysqlDb(dsn string) {
	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		hlog.Fatalf("database initialize failed", err)
		return
	}
	hlog.Infof("数据库连接成功!")
}

func InitAliOSSBucket() {
	AliOSSClient, err := oss.New(Config.AliOSS.Endpoint, Config.AliOSS.AccessKeyId, Config.AliOSS.AccessKeySecret)
	if err != nil {
		hlog.Fatalf("ali oss client init failed, err: %v", err)
	}
	AliOSSBucket, err = AliOSSClient.Bucket(Config.AliOSS.Bucket)
	if err != nil {
		hlog.Fatalf("ali oss bucket init failed, err: %v", err)
	}

	hlog.Infof("init ali oss bucket success, +%v", Config.AliOSS)
}
