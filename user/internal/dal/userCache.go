package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"user/internal/database"
	"user/internal/model"
	"user/internal/repo"
	"user/pkg/constant"

	"github.com/go-redis/redis/v8"
)

type UserCacheDal struct {
	conn *redis.Client
}

var _ repo.UserCacheRepo = (*UserCacheDal)(nil)

func NewUserCacheRepo() *UserCacheDal {
	return &UserCacheDal{
		conn: database.GetRedis(),
	}
}

func (u *UserCacheDal) MGetUserInfo(ctx context.Context, uids []int64) (users []model.User, err error) {
	for _, uid := range uids {
		cacheKey := u.getUserInfoCacheKey(uid)
		c := u.conn.Get(ctx, cacheKey)
		if c.Err() != nil {
			return nil, c.Err()
		}
		var user model.User
		err := json.Unmarshal([]byte(c.Val()), &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return
}

func (u *UserCacheDal) MSetUserInfo(ctx context.Context, users []model.User, expire time.Duration) error {
	for _, user := range users {
		userJson, err := json.Marshal(user)
		if err != nil {
			return err
		}
		cacheKey := u.getUserInfoCacheKey(user.ID)
		// 随机过期时间，防止缓存雪崩
		expireTime := expire + time.Duration(rand.Intn(10))*time.Second
		c := u.conn.Set(ctx, cacheKey, userJson, expireTime)
		if c.Err() != nil {
			return c.Err()
		}
	}
	return nil
}

func (u *UserCacheDal) getUserInfoCacheKey(uid int64) string {
	return fmt.Sprintf(constant.UserInfoCacheKey, uid)
}

func (u *UserCacheDal) SetFollowStats(ctx context.Context, userID int64, followerCount int, followingCount int, expire time.Duration) error {
	stats := map[string]int{
		"followerCount":  followerCount,
		"followingCount": followingCount,
	}
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	cacheKey := u.getFollowStatsCacheKey(userID)
	// 随机过期时间，防止缓存雪崩
	expireTime := expire + time.Duration(rand.Intn(10))*time.Second
	c := u.conn.Set(ctx, cacheKey, data, expireTime)
	if c.Err() != nil {
		return c.Err()
	}
	return nil
}

func (u *UserCacheDal) getFollowStatsCacheKey(uid int64) string {
	return fmt.Sprintf("followStats:%d", uid)
}
