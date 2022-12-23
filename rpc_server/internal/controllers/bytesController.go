package controllers

import (
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
	"github.com/matorix/chat.fdbm.dev/internal/config"
	"github.com/matorix/chat.fdbm.dev/internal/pubsub"
	"github.com/matorix/chat.fdbm.dev/internal/utils"
)

type BytesController struct {
	Pubsub pubsub.PubSub[chatv1.ListenBytesStreamResponse]
}

func NewBytesController() BytesController {
	return BytesController{
		Pubsub: *pubsub.NewPubSub[chatv1.ListenBytesStreamResponse](),
	}
}

func (bc *BytesController) Send(stream *connect.ClientStream[chatv1.SendBytesStreamRequest]) (*chatv1.SendBytesStreamResponse, error) {
	var hash string
	var lowPassword string
	var msg *chatv1.SendBytesStreamRequest
	var id *string
	for stream.Receive() {
		msg = stream.Msg()
		if *id != msg.Data.Id && id != nil {
			return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("error: invalid id. Data Id must be same as the previous stream request"))
		} else if id == nil {
			id = &msg.Data.Id
		}

		if lowPassword != msg.DiscussionInfo.LowPassword {
			lowPassword = msg.DiscussionInfo.LowPassword
			hash = utils.GenerateHash(lowPassword, config.SecretKey)
		}

		bc.Pubsub.Publish(
			*id,
			hash,
			&msg.User.UnsafeId,
			&chatv1.ListenBytesStreamResponse{Data: msg.Data},
		)
	}
	bc.Pubsub.Publish(
		*id,
		hash,
		&msg.User.UnsafeId,
		&chatv1.ListenBytesStreamResponse{Data: &chatv1.BytesData{
			Data:      []byte{},
			Id:        msg.Data.Id,
			Type:      msg.Data.Type,
			Filename:  msg.Data.Filename,
			CreatedAt: int32(time.Now().UnixMilli()),
			Finished:  true,
		}},
	)
	return &chatv1.SendBytesStreamResponse{}, nil
}

func (bc *BytesController) Listen(req *chatv1.ListenBytesStreamRequest, streamRes *connect.ServerStream[chatv1.ListenBytesStreamResponse]) {
	hash := utils.GenerateHash(req.DiscussionInfo.LowPassword, config.SecretKey)
	watchId := bc.Pubsub.Subscribe(req.DiscussionInfo.Id, hash, &req.User.UnsafeId, streamRes)
	defer bc.Pubsub.Unsubscribe(req.DiscussionInfo.Id, watchId)
	for {
		time.Sleep(config.ConnectionCheckInterval)
		err := streamRes.Send(&chatv1.ListenBytesStreamResponse{
			Data: &chatv1.BytesData{
				Data:      []byte{},
				Id:        "none",
				Type:      "ConnectionCheck",
				Filename:  "ConnectionCheck",
				CreatedAt: int32(time.Now().UnixMilli()),
				Finished:  false,
			},
		})
		if err != nil {
			return
		}
	}
}
