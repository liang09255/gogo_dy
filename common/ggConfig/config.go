package ggConfig

import (
	"log"

	"github.com/spf13/viper"
)

var Config = &config{}

func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// 添加配置文件查找路径
	v.AddConfigPath("../common/ggConfig/")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// 解析到Config
	if err := v.Unmarshal(Config); err != nil {
		log.Fatal(err)
	}
}
