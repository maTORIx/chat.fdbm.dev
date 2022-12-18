package pubsub

import (
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
)

type ChatsPubSub struct {
	Mu        *sync.RWMutex
	WatchList map[string]*SubWatchList
}

type SubWatchList struct {
	Mu        *sync.RWMutex
	WatchList map[string]*MessageSubscriberInfo
}

func NewChatsPubSub() *ChatsPubSub {
	return &ChatsPubSub{
		Mu:        &sync.RWMutex{},
		WatchList: map[string]*SubWatchList{},
	}
}

type MessageSubscriberInfo struct {
	Hash      string
	StreamRes *connect.ServerStream[chatv1.GetChatsStreamResponse]
}

func (ps *ChatsPubSub) Subscribe(discussionId, hash string, stream *connect.ServerStream[chatv1.GetChatsStreamResponse]) string {
	id, _ := uuid.NewUUID()
	watchId := id.String()
	ps.Mu.Lock()
	if _, ok := ps.WatchList[discussionId]; !ok {
		ps.WatchList[discussionId] = &SubWatchList{
			Mu:        &sync.RWMutex{},
			WatchList: map[string]*MessageSubscriberInfo{},
		}
	}
	ps.Mu.Unlock()

	ps.WatchList[discussionId].Mu.Lock()
	defer ps.WatchList[discussionId].Mu.Unlock()
	ps.WatchList[discussionId].WatchList[watchId] = &MessageSubscriberInfo{
		Hash:      hash,
		StreamRes: stream,
	}
	return watchId
}

func (ps *ChatsPubSub) Publish(discussion_id, hash string, newChat *chatv1.Chat) {
	if _, ok := ps.WatchList[discussion_id]; !ok {
		return
	}
	ps.WatchList[discussion_id].Mu.RLock()
	for watchId, info := range ps.WatchList[discussion_id].WatchList {
		go func(info *MessageSubscriberInfo, watchId string) {
			if info != nil && info.Hash == hash {
				err := info.StreamRes.Send(&chatv1.GetChatsStreamResponse{Chat: newChat})
				if err != nil {
					ps.WatchList[discussion_id].Mu.Lock()
					ps.WatchList[discussion_id].WatchList[watchId] = nil
					ps.WatchList[discussion_id].Mu.Unlock()
				}
			}
		}(info, watchId)
	}
	ps.WatchList[discussion_id].Mu.RUnlock()
}
