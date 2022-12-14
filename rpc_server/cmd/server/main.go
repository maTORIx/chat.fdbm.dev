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
	"github.com/matorix/chat.fdbm.dev/internal/controllers"
	"github.com/matorix/chat.fdbm.dev/internal/models"
	"github.com/matorix/chat.fdbm.dev/internal/utils"
)

var chatsController controllers.ChatsController
var bytesController controllers.BytesController

func InitializeDB() error {
	chatModel := &models.ChatModel{}
	if err := chatModel.Init(); err != nil {
		return err
	}
	return nil
}

type ChatServer struct{}

func (s *ChatServer) SendChat(
	ctx context.Context,
	req *connect.Request[chatv1.SendChatRequest],
) (*connect.Response[chatv1.SendChatResponse], error) {
	log.Println("Request headers: ", req.Header())
	err := chatsController.AddChat(req.Msg, strings.Split(req.Peer().Addr, ":")[0])
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
	chats, err := chatsController.ListChat(req.Msg)
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
	watchId := chatsController.Pubsub.Subscribe(discussionId, hash, nil, streamRes)
	defer chatsController.Pubsub.Unsubscribe(discussionId, watchId)
	for {
		time.Sleep(config.ConnectionCheckInterval)
		err := streamRes.Send(&chatv1.GetChatsStreamResponse{Chats: []*chatv1.Chat{}})
		if err != nil {
			return nil
		}
	}
}

func (s *ChatServer) SendBytesStream(
	ctx context.Context,
	stream *connect.ClientStream[chatv1.SendBytesStreamRequest],
) (*connect.Response[chatv1.SendBytesStreamResponse], error) {
	log.Println("Request headers: ", stream.RequestHeader())
	result, err := bytesController.Send(stream)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	res := connect.NewResponse(result)
	return res, nil
}

func (s *ChatServer) ListenBytesStream(
	ctx context.Context,
	req *connect.Request[chatv1.ListenBytesStreamRequest],
	streamRes *connect.ServerStream[chatv1.ListenBytesStreamResponse],
) error {
	log.Println("Request headers: ", req.Header())
	bytesController.Listen(req.Msg, streamRes)
	return nil
}

func main() {
	err := InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	chatsController = controllers.NewChatsController()
	bytesController = controllers.NewBytesController()

	server := &ChatServer{}
	mux := http.NewServeMux()
	path, handler := chatv1connect.NewChatServiceHandler(server)
	mux.Handle(path, handler)
	mux.Handle("/", http.FileServer(http.Dir("../../gen/webapp")))
	corsHandler := cors.Default().Handler(h2c.NewHandler(mux, &http2.Server{}))
	tmp := &http.Server{
		Addr:        "localhost:8080",
		Handler:     corsHandler,
		ReadTimeout: 5 * time.Second,
	}
	tmp.ListenAndServe()
}
