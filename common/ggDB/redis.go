package ggDB

import (
	"common/ggConfig"
	"common/ggLog"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func NewRedis() *redis.Client {
	config := ggConfig.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})
	ggLog.Infof("redis connect success")
	return rdb
}
