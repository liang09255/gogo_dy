package repo

import (
	"context"
	"time"
)

// 缓存接口
type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, key string, field string, value string) error
	HKeys(background context.Context, key string) ([]string, error)
	Delete(background context.Context, keys []string) error
}
