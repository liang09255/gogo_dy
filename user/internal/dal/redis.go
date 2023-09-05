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

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result, err := rc.Rdb.Get(ctx, key).Result()
	if err != nil {
		return result, err
	}
	return result, nil
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

func (rc *RedisCache) ListenForExpirationEvents() {
	pubsub := rc.Rdb.Subscribe(context.Background(), "__keyevent@0__:expired") // 注意：这里的0是数据库编号，您可能需要根据实际的Redis配置进行调整
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	if _, err := pubsub.Receive(context.Background()); err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Consume messages.
	for msg := range ch {
		if msg.Channel == "__keyevent@0__:expired" {
			// 这里处理key过期的逻辑
			keyExpired := msg.Payload
			go handleKeyExpiration(keyExpired) // 这是一个示例函数，您需要定义它以处理key的过期事件
		}
	}
}

func handleKeyExpiration(key string) {
	// 这里是处理key过期的具体逻辑，比如同步关注统计表和relation表
}
