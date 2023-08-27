package repo

import (
	"chat/internal/model"
	"context"
)

type ChatRepo interface {
	MAction(ctx context.Context, message *model.Message) (msg string, err error)
	MGetListByMIds(ctx context.Context, mIds []int64) (messages []model.Message, err error)
	MGetList(ctx context.Context, fromId int64, toId int64, preMsgTime int64) (messages []model.Message, err error)
}
