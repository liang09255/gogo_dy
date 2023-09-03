package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
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

func (c *CommentCacheDal) getCommentListKey(vid int64) string {
	return fmt.Sprintf(constant.CommentListKey, vid)
}

//func (f *FavoriteCacheDal) getUserGetFavoriteCountKey(uid int64) string {
//	return fmt.Sprintf(constant.UserGetFavoriteCountKey, uid)
//}
//
//func (f *FavoriteCacheDal) getVideoFavoriteCountKey(vid int64) string {
//	return fmt.Sprintf(constant.VideoFavoriteCountKey, vid)
//}
