package pubsub

import (
	"log"
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
)

type ChatsPubSub struct {
	WatchList sync.Map
}

func NewChatsPubSub() *ChatsPubSub {
	return &ChatsPubSub{
		WatchList: sync.Map{},
	}
}

type MessageSubscriberInfo struct {
	Hash      string
	StreamRes *connect.ServerStream[chatv1.GetChatsStreamResponse]
}

func (ps *ChatsPubSub) Subscribe(discussionId, hash string, stream *connect.ServerStream[chatv1.GetChatsStreamResponse]) string {
	id, _ := uuid.NewUUID()
	watchId := id.String()
	v, _ := ps.WatchList.LoadOrStore(discussionId, &sync.Map{})
	sub, ok := v.(*sync.Map)
	if !ok {
		log.Fatal("fatal error: unexpected error occur")
	}
	sub.Store(watchId, &MessageSubscriberInfo{
		Hash:      hash,
		StreamRes: stream,
	})
	return watchId
}

func (ps *ChatsPubSub) Publish(discussionId, hash string, newChat *chatv1.Chat) {
	v, ok := ps.WatchList.Load(discussionId)
	if !ok {
		return
	}
	sub, _ := v.(*sync.Map)
	sub.Range(func(key any, v any) bool {
		info, _ := v.(*MessageSubscriberInfo)
		watchId, _ := key.(string)

		if info != nil && info.Hash == hash {
			err := info.StreamRes.Send(&chatv1.GetChatsStreamResponse{Chat: newChat})
			if err != nil {
				sub.Delete(watchId)
			}
		}
		return true
	})
}

func (ps *ChatsPubSub) IsDisconnected(discussionId, watchId string) bool {
	v, ok := ps.WatchList.Load(discussionId)
	if !ok {
		return false
	}
	sub, _ := v.(*sync.Map)
	w, sub_ok := sub.Load(watchId)
	if !sub_ok {
		return false
	}
	info, _ := w.(*MessageSubscriberInfo)
	return info == nil
}
