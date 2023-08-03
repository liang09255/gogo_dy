package dal

import (
	"gorm.io/gorm"
	"main/global"
)

type Message struct {
	gorm.Model
	Token      string `gorm:"type:varchar(255);column:token"` // TODO 这边不需要存token，后面删掉
	ToUserID   int64  `gorm:"type:bigint;column:to_user_id"`
	FromUserID int64  `gorm:"type:bigint;column:from_user_id"` // 存储发送者的ID,可以从鉴权token中得到
	ActionType int32  `gorm:"type:int;column:action_type"`
	Content    string `gorm:"type:text;column:content"`
}

type messageDal struct{}

var MessageDal = &messageDal{}

func (m *messageDal) NewMessage(token string, toUserID int64, actionType int32, content string) error {
	t := global.MysqlDB.Create(&Message{
		Token:      token,
		ToUserID:   toUserID,
		// TODO:  From_User_ID，自己的id从token怎么得来？ 不着急现在
		// FromUserID: myselfID
		ActionType: actionType,
		Content:    content,
	})
	return t.Error
}

func (m *messageDal) GetMessages(token string, toUserID int64, preMsgTime int64) ([]Message, error) {
	// FIXME 这个方法感觉是要把我发给对方和对方发给我的都查出来，时间排序
	// 这里暂时只查了我发给对方的
	var messages []Message

	query := global.MysqlDB.Where("to_user_id = ?", toUserID)

	// 如果提供了 preMsgTime，你可能需要添加一个合适的时间过滤条件
	if preMsgTime != 0 {
		// query = query.Where("created_at > ?", preMsgTime) // 假设使用了gorm.Model的CreatedAt字段
		query = query.Where("created_at > ?", time.Unix(preMsgTime, 0))
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
