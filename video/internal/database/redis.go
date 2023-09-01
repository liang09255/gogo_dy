package database

import (
	"common/ggDB"
	"github.com/go-redis/redis/v8"
)

var _redis *redis.Client

func init() {
	_redis = ggDB.NewRedis()
}

func ManualInitRedis() {}

func GetRedis() *redis.Client {
	return _redis
}
