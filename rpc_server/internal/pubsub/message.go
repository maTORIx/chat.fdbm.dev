package pubsub

import (
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
)

type ChatsPubSub struct {
	mu        sync.Mutex
	WatchList map[string]map[string]*MessageSubscriberInfo
}

func NewChatsPubSub() *ChatsPubSub {
	return &ChatsPubSub{
		mu:        sync.Mutex{},
		WatchList: map[string]map[string]*MessageSubscriberInfo{},
	}
}

type MessageSubscriberInfo struct {
	Hash      string
	StreamRes *connect.ServerStream[chatv1.GetChatsStreamResponse]
}

func (ps *ChatsPubSub) Subscribe(discussionId, hash string, stream *connect.ServerStream[chatv1.GetChatsStreamResponse]) string {
	id, _ := uuid.NewUUID()
	watchId := id.String()
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.WatchList[discussionId]; !ok {
		ps.WatchList[discussionId] = map[string]*MessageSubscriberInfo{}
	}
	ps.WatchList[discussionId][watchId] = &MessageSubscriberInfo{
		Hash:      hash,
		StreamRes: stream,
	}
	return watchId
}

func (ps *ChatsPubSub) Publish(discussion_id, hash string, newChat *chatv1.Chat) {
	if _, ok := ps.WatchList[discussion_id]; !ok {
		return
	}
	for watchId, info := range ps.WatchList[discussion_id] {
		go func(info *MessageSubscriberInfo, watchId string) {
			if info != nil && info.Hash == hash {
				err := info.StreamRes.Send(&chatv1.GetChatsStreamResponse{Chat: newChat})
				if err != nil {
					ps.mu.Lock()
					ps.WatchList[discussion_id][watchId] = nil
					ps.mu.Unlock()
				}
			}
		}(info, watchId)
	}
}
