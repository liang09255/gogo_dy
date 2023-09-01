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

func (md *ChatDal) MGetList(ctx context.Context, fromId int64, toId int64, preMsgTime int64) (messages []model.Message, err error) {
	query := md.conn.WithContext(ctx).
		Where("((to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)) AND (created_at > ?)",
			fromId, toId, toId, fromId, time.UnixMicro(preMsgTime)).
		Order("created_at").Find(&messages)
	return messages, query.Error
}
