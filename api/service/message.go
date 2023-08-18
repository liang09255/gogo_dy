package service

import (
	"fmt"
	"main/controller/ctlModel/messageCtlModel"
	"main/dal"
)

type messageService struct{}

var MessageService = &messageService{}

func (m *messageService) SendMessage(fromUsrID, toUserID int64, actionType int32, content string) (msg string, err error) {
	err = dal.MessageDal.NewMessage(fromUsrID, toUserID, actionType, content)
	if err != nil {
		return
	}
	// 如果没有错误，返回成功消息
	msg = fmt.Sprintf("Message sent to user with ID %d, content: %s", toUserID, content)
	return
}

func (m *messageService) GetChatMessages(userID, toUserID int64, preMsgTime int64) ([]messageCtlModel.Message, error) {
	messages, err := dal.MessageDal.GetMessages(userID, toUserID, preMsgTime)
	if err != nil {
		return []messageCtlModel.Message{}, fmt.Errorf("failed to get chat messages: %v", err)
	}
	return convertMessagesToChatResponse(messages), nil
}

func convertMessagesToChatResponse(messages []dal.Message) []messageCtlModel.Message {
	var chatMessages []messageCtlModel.Message
	for _, msg := range messages {
		chatMessage := messageCtlModel.Message{
			ID:         int64(msg.ID),
			Content:    msg.Content,
			CreateTime: msg.CreatedAt.UnixMicro(),
			FromUserID: msg.FromUserID,
			ToUserID:   msg.ToUserID,
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	return chatMessages
}
