package dal

import (
	"common/ggLog"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"user/internal/database"
	"user/internal/repo"
	"user/pkg/constant"
)

type RelationCacheDal struct {
	conn *redis.Client
}

var _ repo.RelationCacheRepo = (*RelationCacheDal)(nil)

var expireTime = 5 * time.Minute

func NewRelationCacheDal() *RelationCacheDal {
	return &RelationCacheDal{
		conn: database.GetRedis(),
	}
}

func (r *RelationCacheDal) ADD(ctx context.Context, uid, target int64) error {
	cacheKey := getBloomCacheKey(uid)
	return r.conn.Do(ctx, "BF.ADD", cacheKey, target).Err()
}

func (r *RelationCacheDal) MADD(ctx context.Context, uid int64, target []int64) error {
	cacheKey := getBloomCacheKey(uid)
	args := []interface{}{"BF.MADD", cacheKey}
	for _, t := range target {
		args = append(args, t)
	}
	err := r.conn.Do(ctx, args...).Err()
	if err != nil {
		return err
	}
	r.conn.Expire(ctx, cacheKey, expireTime)
	return nil
}

func (r *RelationCacheDal) EXISTS(ctx context.Context, uid, target int64) (bool, error) {
	cacheKey := getBloomCacheKey(uid)
	return r.conn.Do(ctx, "BF.EXISTS", cacheKey, target).Bool()
}

func (r *RelationCacheDal) MEXISTS(ctx context.Context, uid int64, target []int64) ([]bool, error) {
	cacheKey := getBloomCacheKey(uid)
	args := []interface{}{"BF.MEXISTS", cacheKey}
	for _, t := range target {
		args = append(args, t)
	}
	return r.conn.Do(ctx, args...).BoolSlice()
}

func (r *RelationCacheDal) KeyExist(ctx context.Context, uid int64) bool {
	cacheKey := getBloomCacheKey(uid)
	return r.conn.Exists(ctx, cacheKey).Val() == 1
}

func (r *RelationCacheDal) Delete(ctx context.Context, uid int64) {
	cacheKey := getBloomCacheKey(uid)
	if err := r.conn.Del(ctx, cacheKey).Err(); err != nil {
		ggLog.Error("redis delete error: %v", err)
	}
}

func getBloomCacheKey(uid int64) string {
	return fmt.Sprintf(constant.CacheRelationBloom, uid)
}
