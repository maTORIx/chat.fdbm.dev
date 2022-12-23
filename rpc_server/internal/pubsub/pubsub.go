package pubsub

import (
	"fmt"
	"log"
	"sync"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
)

type PubSub[T any] struct {
	WatchList sync.Map
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		WatchList: sync.Map{},
	}
}

type SubscriberInfo[T any] struct {
	Hash      string
	IgnoreId  *string
	StreamRes *connect.ServerStream[T]
}

func (ps *PubSub[T]) Subscribe(discussionId, hash string, ignoreId *string, stream *connect.ServerStream[T]) string {
	id, _ := uuid.NewUUID()
	watchId := id.String()
	v, _ := ps.WatchList.LoadOrStore(discussionId, &sync.Map{})
	sub, ok := v.(*sync.Map)
	if !ok {
		log.Fatal("fatal error: unexpected error occur")
	}
	sub.Store(watchId, &SubscriberInfo[T]{
		Hash:      hash,
		StreamRes: stream,
		IgnoreId:  ignoreId,
	})
	return watchId
}

func (ps *PubSub[T]) Unsubscribe(discussionId, watchId string) {
	fmt.Println("Unsubscribe!")
	v, ok := ps.WatchList.Load(discussionId)
	if !ok {
		return
	}
	sub, _ := v.(*sync.Map)
	sub.Delete(watchId)
}

func (ps *PubSub[T]) Publish(discussionId, hash string, ignoreId *string, item *T) {
	v, ok := ps.WatchList.Load(discussionId)
	if !ok || item == nil {
		return
	}
	sub, _ := v.(*sync.Map)
	sub.Range(func(key any, v any) bool {
		info, _ := v.(*SubscriberInfo[T])
		if info != nil && info.Hash == hash && (info.IgnoreId == nil || info.IgnoreId != ignoreId) {
			go info.StreamRes.Send(item)
		}
		return true
	})
}
