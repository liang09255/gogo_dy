package service

import (
	"chat/internal/domain"
	"common/ggIDL/chat"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatService struct {
	chat.UnsafeChatServer
	chatDomain *domain.ChatDomain
}

var _ chat.ChatServer = (*ChatService)(nil)

func New() *ChatService {
	return &ChatService{
		chatDomain: domain.NewChatDomain(),
	}
}

func (cs *ChatService) Action(ctx context.Context, request *chat.ChatActionRequest) (*chat.ChatActionResponse, error) {
	var chatInfo = &chat.ChatInfoModel{
		Content:    request.Content,
		CreateTime: timestamppb.Now(),
		FromUserId: request.FromUserId,
		ToUserId:   request.ToUserId,
		ActionType: int32(request.ActionType),
	}
	msg, err := cs.chatDomain.Action(ctx, chatInfo)
	return &chat.ChatActionResponse{Msg: msg}, err
}

func (cs *ChatService) List(ctx context.Context, request *chat.ListRequest) (*chat.ListResponse, error) {
	list, err := cs.chatDomain.List(ctx, request.FromUserId, request.ToUserId, request.PreMsgTime)
	var models []*chat.ChatInfoModel
	for _, m := range list {
		models = append(models, &chat.ChatInfoModel{
			Id:         m.Id,
			Content:    m.Content,
			CreateTime: m.CreateTime,
			FromUserId: m.FromUserId,
			ToUserId:   m.ToUserId,
			ActionType: m.ActionType,
		})
	}
	if err != nil {
		return nil, err
	}
	return &chat.ListResponse{List: models}, nil
}
