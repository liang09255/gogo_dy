package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"chat/internal/database"
	"chat/internal/model"
	"chat/internal/repo"
	"chat/pkg/constant"

	"github.com/go-redis/redis/v8"
)

type ChatCacheDal struct {
	conn *redis.Client
}

var _ repo.ChatCacheRepo = (*ChatCacheDal)(nil)

func NewChatCacheDal() *ChatCacheDal {
	return &ChatCacheDal{
		conn: database.GetRedis(),
	}
}

func (c *ChatCacheDal) MGetChatInfo(ctx context.Context, fromId int64, toId int64) (messages []model.Message, err error) {
	cacheKey := c.getMessageCacheKey(fromId + toId)
	result := c.conn.Get(ctx, cacheKey)
	if result.Err() != nil {
		return nil, result.Err()
	}
	err = json.Unmarshal([]byte(result.Val()), &messages)
	if err != nil {
		return nil, err
	}
	return
}

func (c *ChatCacheDal) MSetChatInfo(ctx context.Context, messages []model.Message, expire time.Duration) error {
	for _, message := range messages {
		msgJson, err := json.Marshal(message)
		if err != nil {
			return err
		}
		cacheKey := c.getMessageCacheKey(message.FromUserID + message.ToUserID)
		// 随机过期时间，防止缓存雪崩
		expireTime := expire + time.Duration(rand.Intn(10))*time.Second
		c := c.conn.Set(ctx, cacheKey, msgJson, expireTime)
		if c.Err() != nil {
			return c.Err()
		}
	}
	return nil
}

func (c *ChatCacheDal) getMessageCacheKey(mId int64) string {
	return fmt.Sprintf(constant.UserInfoCacheKey, mId)
}
