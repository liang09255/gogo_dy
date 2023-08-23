package global

import (
	"common/ggConfig"
	"common/ggDB"
	"common/ggIDL/relation"
	"common/ggIDL/user"
	"common/ggRPC"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB
var AliOSSBucket *oss.Bucket

var UserClient user.UserClient

var RelationClient relation.RelationClient

func Init() {
	InitMysqlDb()
	InitAliOSSBucket()
	InitUserClient()
}

func InitUserClient() {
	UserClient = ggRPC.GetUserClient()
	RelationClient = ggRPC.GetRelationClient()
}

func InitMysqlDb() {
	MysqlDB = ggDB.NewMySQL()
}

func InitAliOSSBucket() {
	AliOSSClient, err := oss.New(ggConfig.Config.AliOSS.Endpoint, ggConfig.Config.AliOSS.AccessKeyId, ggConfig.Config.AliOSS.AccessKeySecret)
	if err != nil {
		hlog.Fatalf("ali oss client init failed, err: %v", err)
	}
	AliOSSBucket, err = AliOSSClient.Bucket(ggConfig.Config.AliOSS.Bucket)
	if err != nil {
		hlog.Fatalf("ali oss bucket init failed, err: %v", err)
	}
	hlog.Infof("init ali oss bucket success")
}
