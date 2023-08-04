package service

import (
	"fmt"
	"main/dal"
)


type messageService struct{}

var MessageService = &messageService{}


type ChatMessageResponse struct {
	ID         uint  `json:"id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
}


func (m *messageService) SendMessage(token string, toUserID int64, actionType int32, content string) (msg string, err error) {
	err = dal.MessageDal.NewMessage(token, toUserID, actionType, content)
	if err != nil {
		return
	}
	// 如果没有错误，返回成功消息
	msg = fmt.Sprintf("Message sent to user with ID %d, content: %s", toUserID, content)
	return
}




func (m *messageService) GetChatMessages(token string, toUserID int64, preMsgTime int64) ([]ChatMessageResponse, error) {
	messages, err := dal.MessageDal.GetMessages(token, toUserID, preMsgTime)
	if err != nil {
		return []ChatMessageResponse{}, fmt.Errorf("failed to get chat messages: %v", err)
	}

	return convertMessagesToChatResponse(messages), nil
}




func convertMessagesToChatResponse(messages []dal.Message) []ChatMessageResponse {
	chatMessages := []ChatMessageResponse{}
	for _, msg := range messages {
		chatMessage := ChatMessageResponse{
			ID:         msg.ID,
			Content:    msg.Content,
			CreateTime: msg.CreatedAt.Unix(),
			FromUserID: msg.FromUserID,
			ToUserID:   msg.ToUserID,
		}
		chatMessages = append(chatMessages, chatMessage)
	}

	return chatMessages
	// chatResponseBytes, err := json.Marshal(chatMessages)
    // if err != nil {
    //     return ""
    // }

    // return string(chatResponseBytes)
}

