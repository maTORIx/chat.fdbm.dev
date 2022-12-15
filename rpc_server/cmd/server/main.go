package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	_ "github.com/mattn/go-sqlite3"

	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
	"github.com/matorix/chat.fdbm.dev/gen/chat/v1/chatv1connect"
	"github.com/matorix/chat.fdbm.dev/internal/config"
	"github.com/matorix/chat.fdbm.dev/internal/models"
	"github.com/matorix/chat.fdbm.dev/internal/pubsub"
	"github.com/matorix/chat.fdbm.dev/internal/utils"
)

var chatsPubSub *pubsub.ChatsPubSub

func InitializeDB() error {
	chatModel := &models.ChatModel{}
	if err := chatModel.Init(); err != nil {
		return err
	}
	return nil
}

func AddChat(chat *chatv1.SendChatRequest, ip_address string) error {
	chatModel := &models.ChatModel{}
	newChat, err := chatModel.Create(chat.DiscussionInfo.Id, chat.Name, chat.Body, ip_address, chat.DiscussionInfo.LowPassword)
	if err != nil {
		return err
	}

	chatsPubSub.Publish(newChat.DiscussionId, newChat.Hash, &chatv1.Chat{
		Id:        newChat.Id,
		UserId:    newChat.Uid,
		Name:      newChat.Name,
		Body:      newChat.Body,
		CreatedAt: newChat.CreatedAt,
	})
	return nil
}

func ListChat(info *chatv1.GetChatsRequest) ([]*chatv1.Chat, error) {
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

type ChatServer struct{}

func (s *ChatServer) SendChat(
	ctx context.Context,
	req *connect.Request[chatv1.SendChatRequest],
) (*connect.Response[chatv1.SendChatResponse], error) {
	log.Println("Request headers: ", req.Header())
	err := AddChat(req.Msg, strings.Split(req.Peer().Addr, ":")[0])
	if err != nil {
		log.Fatal("error")
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	res := connect.NewResponse(&chatv1.SendChatResponse{})
	res.Header().Set("Chat-Version", "v1")
	return res, nil
}

func (s *ChatServer) GetChats(
	ctx context.Context,
	req *connect.Request[chatv1.GetChatsRequest],
) (*connect.Response[chatv1.GetChatsResponse], error) {
	log.Println("Request headers: ", req.Header())
	chats, err := ListChat(req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	res := connect.NewResponse(&chatv1.GetChatsResponse{Chats: chats})
	return res, nil
}

func (s *ChatServer) GetChatsStream(
	ctx context.Context,
	req *connect.Request[chatv1.GetChatsStreamRequest],
	streamRes *connect.ServerStream[chatv1.GetChatsStreamResponse],
) error {
	log.Println("Request headers: ", req.Header())

	discussionId := req.Msg.DiscussionInfo.Id
	hash := utils.GenerateHash(req.Msg.DiscussionInfo.LowPassword, config.SecretKey)
	chatsPubSub.Subscribe(discussionId, hash, streamRes)
	for {
		time.Sleep(time.Second * 30)
	}
}

func main() {
	err := InitializeDB()
	if err != nil {
		log.Fatal(err)
	}
	chatsPubSub = pubsub.NewChatsPubSub()

	server := &ChatServer{}
	mux := http.NewServeMux()
	path, handler := chatv1connect.NewChatServiceHandler(server)
	mux.Handle(path, handler)
	corsHandler := cors.Default().Handler(h2c.NewHandler(mux, &http2.Server{}))
	http.ListenAndServe(
		"localhost:8080",
		corsHandler,
	)
}
