package dal

import (
	"chat/internal/database"
	"chat/internal/model"
	"chat/internal/repo"
	"context"
	"time"
)

type ChatDal struct {
	conn *database.GormConn
}

var _ repo.ChatRepo = (*ChatDal)(nil)

func NewChatDao() *ChatDal {
	return &ChatDal{
		conn: database.New(),
	}
}

func (md *ChatDal) MAction(ctx context.Context, message *model.Message) (msg string, err error) {
	return "发送成功", md.conn.WithContext(ctx).Create(message).First(message).Error
}

func (md *ChatDal) MGetListByMIds(ctx context.Context, mIds []int64) (messages []model.Message, err error) {
	t := md.conn.WithContext(ctx).Where("id in ?", mIds).Find(&messages)
	return messages, t.Error
}
func (md *ChatDal) MGetList(ctx context.Context, fromId int64, toId int64, preMsgTime int64) (messages []model.Message, err error) {
	t := md.conn.WithContext(ctx).Where("from_user_id = ? AND to_user_id = ? AND created_at > ?", fromId, toId, time.UnixMicro(preMsgTime)).Order("created_at").Find(&messages)
	return messages, t.Error
}
