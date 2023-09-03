package dal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"time"
	"video/internal/database"
	"video/internal/repo"
	"video/pkg/constant"
)

type FavoriteCacheDal struct {
	conn *redis.Client
}

var _ repo.FavoriteCacheRepo = (*FavoriteCacheDal)(nil)

func NewFavoriteCacheRepo() *FavoriteCacheDal {
	return &FavoriteCacheDal{
		conn: database.GetRedis(),
	}
}

// 获得视频点赞缓存
func (f *FavoriteCacheDal) GetVideoFavoriteCount(ctx context.Context, vid int64) (int64, bool, error) {
	key := f.getVideoFavoriteCountKey(vid)
	c := f.conn.Get(ctx, key)
	if c.Err() == redis.Nil {
		return -1, false, nil
	}
	if c.Err() != nil {
		return -1, false, c.Err()
	}

	// 转为整数
	count, err := strconv.ParseInt(c.Val(), 10, 64)
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return -1, true, err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return count, true, err
}

// 获得用户点赞数
func (f *FavoriteCacheDal) GetUserFavoriteCount(ctx context.Context, uid int64) (int64, bool, error) {
	key := f.getUserFavoriteCountKey(uid)
	c := f.conn.Get(ctx, key)
	if c.Err() == redis.Nil {
		return -1, false, nil
	}
	if c.Err() != nil {
		return -1, false, c.Err()
	}

	// 解析
	// 转为整数
	count, err := strconv.ParseInt(c.Val(), 10, 64)
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return -1, true, err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return count, true, err
}

// 获得用户获赞数
func (f *FavoriteCacheDal) GetUserGetFavoriteCount(ctx context.Context, uid int64) (int64, bool, error) {
	key := f.getUserGetFavoriteCountKey(uid)
	c := f.conn.Get(ctx, key)
	if c.Err() == redis.Nil {
		return -1, false, nil
	}
	if c.Err() != nil {
		return -1, false, c.Err()
	}

	// 转为整数
	count, err := strconv.ParseInt(c.Val(), 10, 64)

	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return -1, true, err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}

	return count, true, err
}

// 设置超时时间,如果中间被查询了还可以延长
func (f FavoriteCacheDal) SetVideoFavoriteCount(ctx context.Context, vid int64, value int64, expire time.Duration) error {
	key := f.getVideoFavoriteCountKey(vid)
	expireTime := expire + time.Duration(rand.Intn(10))*time.Second
	c := f.conn.Set(ctx, key, value, expireTime)

	return c.Err()
}

func (f FavoriteCacheDal) SetUserFavoriteCount(ctx context.Context, uid int64, value int64, expire time.Duration) error {
	key := f.getUserFavoriteCountKey(uid)
	expireTime := expire + time.Duration(rand.Intn(10))*time.Second
	c := f.conn.Set(ctx, key, value, expireTime)
	return c.Err()
}

func (f FavoriteCacheDal) SetUseGetFavoriteCount(ctx context.Context, uid int64, value int64, expire time.Duration) error {
	key := f.getUserGetFavoriteCountKey(uid)
	expireTime := expire + time.Duration(rand.Intn(10))*time.Second
	c := f.conn.Set(ctx, key, value, expireTime)
	return c.Err()
}

func (f *FavoriteCacheDal) IncrVideoFavoriteCount(ctx context.Context, vid int64) error {
	// 自增
	key := f.getVideoFavoriteCountKey(vid)
	_, err := f.conn.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return err
}

func (f *FavoriteCacheDal) IncrUserFavoriteCount(ctx context.Context, uid int64) error {
	// 自增
	key := f.getUserFavoriteCountKey(uid)
	_, err := f.conn.Incr(ctx, key).Result()
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return err
}

func (f *FavoriteCacheDal) IncrUserGetFavoriteCount(ctx context.Context, uid int64) error {
	// 自增
	key := f.getUserGetFavoriteCountKey(uid)
	_, err := f.conn.Incr(ctx, key).Result()
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return err
}

func (f *FavoriteCacheDal) DecrVideoFavoriteCount(ctx context.Context, vid int64) error {
	key := f.getVideoFavoriteCountKey(vid)
	_, err := f.conn.Decr(ctx, key).Result()
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return err
}

func (f *FavoriteCacheDal) DecrUserFavoriteCount(ctx context.Context, uid int64) error {
	// 自增
	key := f.getUserFavoriteCountKey(uid)
	_, err := f.conn.Decr(ctx, key).Result()
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return err
}

func (f *FavoriteCacheDal) DecrUserGetFavoriteCount(ctx context.Context, uid int64) error {
	// 自增
	key := f.getUserGetFavoriteCountKey(uid)
	_, err := f.conn.Decr(ctx, key).Result()
	ttlResult, err := f.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < 30*time.Second {
		f.conn.Expire(ctx, key, 30*time.Second)
	}
	return err
}

func (f *FavoriteCacheDal) getUserFavoriteCountKey(uid int64) string {
	return fmt.Sprintf(constant.UserFavoriteCountKey, uid)
}

func (f *FavoriteCacheDal) getUserGetFavoriteCountKey(uid int64) string {
	return fmt.Sprintf(constant.UserGetFavoriteCountKey, uid)
}

func (f *FavoriteCacheDal) getVideoFavoriteCountKey(vid int64) string {
	return fmt.Sprintf(constant.VideoFavoriteCountKey, vid)
}
