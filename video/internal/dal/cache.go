package dal

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// 缓存实现
var Rc *RedisCache

//func init() {
//	rdb := redis.NewClient(config.C.InitRedisOptions())
//	Rc = &RedisCache{
//		rdb: rdb,
//	}
//}

type RedisCache struct {
	Rdb *redis.Client
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := rc.Rdb.Set(ctx, key, value, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//r:= rc.Rdb.Get(ctx,key)
	result, err := rc.Rdb.Get(ctx, key).Result()
	// 如果键不存在
	if err == redis.Nil {
		return result, false, nil
	}
	if err != nil {
		return result, false, err
	}
	return result, true, nil
}

func (rc *RedisCache) HSet(ctx context.Context, key, field, value string) error {
	_, err := rc.Rdb.HSet(ctx, key, field, value).Result()
	return err
}

func (rc *RedisCache) HKeys(ctx context.Context, key string) ([]string, error) {
	result, err := rc.Rdb.HKeys(ctx, key).Result()
	return result, err
}

func (rc *RedisCache) Delete(ctx context.Context, keys []string) error {
	_, err := rc.Rdb.Del(ctx, keys...).Result()
	return err
}
