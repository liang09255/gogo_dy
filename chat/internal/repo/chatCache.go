package repo

import (
	"chat/internal/model"
	"context"
	"time"
)

type ChatCacheRepo interface {
	MGetChatInfo(ctx context.Context, fromId int64, toId int64) (messages []model.Message, err error)
	MSetChatInfo(ctx context.Context, messages []model.Message, expire time.Duration) error
}
