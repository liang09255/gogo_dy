package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"time"
	"video/internal/database"
	"video/internal/model"
	"video/internal/repo"
	"video/pkg/constant"
)

type CommentCacheDal struct {
	conn *redis.Client
}

var _ repo.CommentCacheRepo = (*CommentCacheDal)(nil)

func NewCommentCacheRepo() *CommentCacheDal {
	return &CommentCacheDal{
		conn: database.GetRedis(),
	}
}

func (c *CommentCacheDal) SetCommentList(ctx context.Context, vid int64, commentList []model.Comment) error {
	key := c.getCommentListKey(vid)
	// 将数据解析出来
	data, err := json.Marshal(commentList)
	if err != nil {
		return err
	}
	expireTime := constant.CommentExpireTime + time.Duration(rand.Intn(10))*time.Second
	s := c.conn.Set(ctx, key, data, expireTime)
	return s.Err()
}

func (c *CommentCacheDal) DelCommentList(ctx context.Context, vid int64) error {
	key := c.getCommentListKey(vid)
	s := c.conn.Del(ctx, key)
	return s.Err()
}

func (c *CommentCacheDal) DelDelayCommentList(ctx context.Context, vid int64) error {
	time.Sleep(500 * time.Millisecond)
	return c.DelCommentList(ctx, vid)
}

func (c *CommentCacheDal) GetCommentList(ctx context.Context, vid int64) ([]model.Comment, bool, error) {
	key := c.getCommentListKey(vid)
	s := c.conn.Get(ctx, key)
	if s.Err() == redis.Nil {
		return []model.Comment{}, false, nil
	}
	if s.Err() != nil {
		return []model.Comment{}, false, nil
	}
	var commentList []model.Comment
	err := json.Unmarshal([]byte(s.Val()), &commentList)
	if err != nil {
		return commentList, true, err
	}

	// 延长缓存
	ttlResult, err := c.conn.TTL(ctx, key).Result()
	if err != nil {
		return commentList, true, err
	}
	if ttlResult < constant.CommentExpireTime {
		c.conn.Expire(ctx, key, constant.CommentExpireTime)
	}
	return commentList, true, nil
}

func (c *CommentCacheDal) SetVideoCommentCount(ctx context.Context, vid int64, value int64) error {
	key := c.getVideoCommentCountKey(vid)
	expireTime := constant.CommentExpireTime + time.Duration(rand.Intn(10))*time.Second
	s := c.conn.Set(ctx, key, value, expireTime)
	return s.Err()
}

func (c *CommentCacheDal) GetVideoCommentCount(ctx context.Context, vid int64) (int64, bool, error) {
	key := c.getVideoCommentCountKey(vid)
	s := c.conn.Get(ctx, key)
	if s.Err() == redis.Nil {
		return -1, false, nil
	}
	if s.Err() != nil {
		return -1, false, s.Err()
	}

	// 转为整数
	count, err := strconv.ParseInt(s.Val(), 10, 64)
	ttlResult, err := c.conn.TTL(ctx, key).Result()
	if ttlResult < constant.CommentExpireTime {
		c.conn.Expire(ctx, key, constant.CommentExpireTime)
	}

	return count, true, err
}

func (c *CommentCacheDal) IncrVideoCommentCount(ctx context.Context, vid int64) error {
	key := c.getVideoCommentCountKey(vid)
	_, err := c.conn.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	ttlResult, err := c.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < constant.CommentExpireTime {
		c.conn.Expire(ctx, key, constant.CommentExpireTime)
	}
	return err
}

func (c *CommentCacheDal) DecrVideoCommentCount(ctx context.Context, vid int64) error {
	key := c.getVideoCommentCountKey(vid)
	_, err := c.conn.Decr(ctx, key).Result()
	if err != nil {
		return err
	}
	ttlResult, err := c.conn.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttlResult < constant.CommentExpireTime {
		c.conn.Expire(ctx, key, constant.CommentExpireTime)
	}
	return err
}

func (c *CommentCacheDal) getCommentListKey(vid int64) string {
	return fmt.Sprintf(constant.CommentListKey, vid)
}

func (c *CommentCacheDal) getVideoCommentCountKey(vid int64) string {
	return fmt.Sprintf(constant.VideoCommentCountKey, vid)
}

//func (f *FavoriteCacheDal) getUserGetFavoriteCountKey(uid int64) string {
//	return fmt.Sprintf(constant.UserGetFavoriteCountKey, uid)
//}
//
//func (f *FavoriteCacheDal) getVideoFavoriteCountKey(vid int64) string {
//	return fmt.Sprintf(constant.VideoFavoriteCountKey, vid)
//}
