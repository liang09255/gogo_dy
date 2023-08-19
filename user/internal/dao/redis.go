package dao

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	Rdb *redis.Client
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := rc.Rdb.Set(ctx, key, value, expire).Err()
	if err != nil {
		hlog.CtxErrorf(ctx, "Redis Put error: %v", err)
		return err
	}
	return nil
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result, err := rc.Rdb.Get(ctx, key).Result()
	if err != nil {
		hlog.CtxErrorf(ctx, "Redis Get error: %v", err)
		return result, err
	}
	return result, nil
}

func (rc *RedisCache) HSet(ctx context.Context, key, field, value string) error {
	_, err := rc.Rdb.HSet(ctx, key, field, value).Result()
	if err != nil {
		hlog.CtxErrorf(ctx, "Redis HSet error: %v", err)
		return err
	}
	return nil
}

func (rc *RedisCache) HKeys(ctx context.Context, key string) ([]string, error) {
	result, err := rc.Rdb.HKeys(ctx, key).Result()
	if err != nil {
		hlog.CtxErrorf(ctx, "Redis HKeys error: %v", err)
		return result, err
	}
	return result, nil
}

func (rc *RedisCache) Delete(ctx context.Context, keys []string) error {
	_, err := rc.Rdb.Del(ctx, keys...).Result()
	if err != nil {
		hlog.CtxErrorf(ctx, "Redis Delete error: %v", err)
		return err
	}
	return nil
}
