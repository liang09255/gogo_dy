package dal

import (
	"main/global"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ToUserID   int64  `gorm:"type:bigint;column:to_user_id"`
	FromUserID int64  `gorm:"type:bigint;column:from_user_id"` // 存储发送者的ID,可以从鉴权token中得到
	ActionType int32  `gorm:"type:int;column:action_type"`
	Content    string `gorm:"type:text;column:content"`
}

type messageDal struct{}

var MessageDal = &messageDal{}

func (m *messageDal) NewMessage(FromUserID, toUserID int64, actionType int32, content string) error {
	t := global.MysqlDB.Create(&Message{
		ToUserID:   toUserID,
		FromUserID: FromUserID,
		ActionType: actionType,
		Content:    content,
	})
	return t.Error
}

func (m *messageDal) GetMessages(userID, toUserID int64, preMsgTime int64) ([]Message, error) {
	var messages []Message
	query := global.MysqlDB.
		Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?) AND (created_at > ?)",
			userID, toUserID, toUserID, userID, time.UnixMicro(preMsgTime)).
		Order("created_at").Find(&messages)
	return messages, query.Error
}
