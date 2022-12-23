package controllers

import (
	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
	"github.com/matorix/chat.fdbm.dev/internal/models"
	"github.com/matorix/chat.fdbm.dev/internal/pubsub"
)

type ChatsController struct {
	Pubsub pubsub.PubSub[chatv1.GetChatsStreamResponse]
}

func NewChatsController() ChatsController {
	return ChatsController{
		Pubsub: *pubsub.NewPubSub[chatv1.GetChatsStreamResponse](),
	}
}

func (c *ChatsController) AddChat(chat *chatv1.SendChatRequest, ip_address string) error {
	chatModel := &models.ChatModel{}
	newChat, err := chatModel.Create(chat.DiscussionInfo.Id, chat.Name, chat.Body, ip_address, chat.DiscussionInfo.LowPassword)
	if err != nil {
		return err
	}
	tmp := &chatv1.Chat{
		Id:        newChat.Id,
		UserId:    newChat.Uid,
		Name:      newChat.Name,
		Body:      newChat.Body,
		CreatedAt: newChat.CreatedAt,
	}
	chats := []*chatv1.Chat{tmp}
	resp := &chatv1.GetChatsStreamResponse{Chats: chats}
	c.Pubsub.Publish(newChat.DiscussionId, newChat.Hash, nil, resp)
	return nil
}

func (c *ChatsController) ListChat(info *chatv1.GetChatsRequest) ([]*chatv1.Chat, error) {
	chatModel := &models.ChatModel{}
	chats, err := chatModel.List(
		info.DiscussionInfo.Id,
		info.DiscussionInfo.LowPassword,
		&info.PageingInfo.Cursor,
		int(info.PageingInfo.Limit),
		info.PageingInfo.EarlierAt,
	)
	if err != nil {
		return nil, err
	}
	result := []*chatv1.Chat{}
	for _, v := range chats {
		result = append(result, &chatv1.Chat{
			Id:        v.Id,
			UserId:    v.Uid,
			Name:      v.Name,
			Body:      v.Body,
			CreatedAt: v.CreatedAt,
		})
	}

	return result, nil
}
