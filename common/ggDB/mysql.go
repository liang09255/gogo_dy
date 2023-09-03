package ggDB

import (
	"common/ggConfig"
	"common/ggLog"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL() *gorm.DB {
	config := ggConfig.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Db,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		ggLog.Fatalf("database initialize failed", err)
	}

	ggLog.Infof("数据库连接成功!")
	return db
}
