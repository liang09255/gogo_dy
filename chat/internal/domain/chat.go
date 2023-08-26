package domain

import (
	"chat/internal/dal"
	"chat/internal/model"
	"chat/internal/repo"
	"common/ggIDL/chat"
	"common/ggLog"
	"context"
	"github.com/bytedance/gopkg/util/gopool"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ChatDomain struct {
	chatRepo      repo.ChatRepo
	tranRepo      repo.TranRepo
	chatCacheRepo repo.ChatCacheRepo
}

func NewChatDomain() *ChatDomain {
	return &ChatDomain{
		chatRepo:      dal.NewChatDao(),
		tranRepo:      dal.NewTranRepo(),
		chatCacheRepo: dal.NewChatCacheDal(),
	}
}

func (cd *ChatDomain) List(ctx context.Context, fromId int64, toId int64, preMsgTime int64) (messageList []chat.ChatInfoModel, err error) {
	var messages []model.Message
	// 先尝试从缓存中获取
	messages, err = cd.chatCacheRepo.MGetChatInfo(ctx, fromId, toId)
	// 检查是否有未命中缓存的
	var cacheMIds []int64
	for _, u := range messages {
		cacheMIds = append(cacheMIds, u.FromUserID+u.ToUserID)
	}
	ggLog.Debugf("cacheMIds:%v", cacheMIds)
	// 从数据库中获取
	if len(cacheMIds) == 0 {
		msgList, err := cd.chatRepo.MGetList(ctx, fromId, toId, preMsgTime)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msgList...)
		// 异步写缓存
		gopool.Go(func() {
			if err := cd.chatCacheRepo.MSetChatInfo(context.Background(), msgList, 1*time.Second); err != nil {
				ggLog.Error("写入缓存失败:", err)
			}
			ggLog.Debugf("写入缓存成功, mid:%v", msgList)
		})
	}
	// 封装返回值
	for _, m := range messages {
		messageList = append(messageList, chat.ChatInfoModel{
			Id:         m.ID,
			Content:    m.Content,
			CreateTime: timestamppb.New(m.CreatedAt),
			FromUserId: m.FromUserID,
			ToUserId:   m.ToUserID,
			ActionType: m.ActionType,
		})
	}
	return messageList, nil

}
func (cd *ChatDomain) Action(ctx context.Context, chat *chat.ChatInfoModel) (msg string, err error) {
	// 写入数据库
	var message = &model.Message{
		ID:         chat.Id,
		Content:    chat.Content,
		FromUserID: chat.FromUserId,
		ToUserID:   chat.ToUserId,
		ActionType: chat.ActionType,
	}
	msg, err = cd.chatRepo.MAction(ctx, message)
	return
}
