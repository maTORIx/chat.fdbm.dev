package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	_ "github.com/mattn/go-sqlite3"

	chatv1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
	"github.com/matorix/chat.fdbm.dev/gen/chat/v1/chatv1connect"
)

var SECRET_KEY string = "SECRET_KEY"

type WatchInfo struct {
	Hash      string
	StreamRes *connect.ServerStream[chatv1.GetChatsStreamResponse]
}

var watchList map[string]map[string]*WatchInfo = map[string]map[string]*WatchInfo{}

// db
func InitializeDB(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	initializeSql := `
	    create table if not exists chats (
			id integer not null primary key autoincrement,
			discussion_id string not null,
			name string not null,
			message string not null,
			hash string not null,
			created_at string not null
		);
	`
	_, err = db.Exec(initializeSql)
	if err != nil {
		return err
	}
	return nil
}

func generateHash(s, salt string) string {
	r := sha256.Sum256([]byte(s + salt))
	return hex.EncodeToString(r[:])
}

func AddChat(chat *chatv1.SendChatRequest) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	utcTime := strconv.Itoa(int(time.Now().UnixMilli()))
	hash := generateHash(chat.LowPassword, SECRET_KEY)

	insertSql := `insert into chats(discussion_id, name, message, hash, created_at) values (?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(chat.DiscussionId, chat.Name, chat.Message, hash, utcTime)
	if err != nil {
		return err
	}
	chatId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	announceToWatchList(chat.DiscussionId, hash, &chatv1.Chat{
		Id:        strconv.Itoa(int(chatId)),
		Name:      chat.Name,
		Message:   chat.Message,
		CreatedAt: utcTime,
	})
	return nil
}

func ListChat(info *chatv1.GetChatsRequest) ([]*chatv1.Chat, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// listSql := `select id, name, message, created_at from chats`
	listSql := `select id, name, message, created_at from chats where discussion_id = ? and hash = ?`
	stmt, err := db.Prepare(listSql)
	if err != nil {
		return nil, err
	}
	// rows, err := stmt.Query()
	hash := generateHash(info.LowPassword, SECRET_KEY)
	rows, err := stmt.Query(info.DiscussionId, hash)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, err
	}
	chats := []*chatv1.Chat{}
	defer rows.Close()
	for rows.Next() {
		chat := &chatv1.Chat{}
		if err := rows.Scan(&chat.Id, &chat.Name, &chat.Message, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func addWatchList(discussion_id string, info *WatchInfo) string {
	id, _ := uuid.NewUUID()
	strid := id.String()
	if _, ok := watchList[discussion_id]; !ok {
		watchList[discussion_id] = map[string]*WatchInfo{}
	}
	watchList[discussion_id][strid] = info
	return strid
}

func announceToWatchList(discussion_id, hash string, newChat *chatv1.Chat) {
	if _, ok := watchList[discussion_id]; !ok {
		return
	}
	var chats []*chatv1.Chat
	chats = append(chats, newChat)

	for _, info := range watchList[discussion_id] {
		go func(info *WatchInfo) {
			if info != nil && info.Hash == hash {
				info.StreamRes.Send(&chatv1.GetChatsStreamResponse{Chats: chats})
			}
		}(info)
	}
}

// grpc server

type ChatServer struct{}

func (s *ChatServer) Greet(
	ctx context.Context,
	req *connect.Request[chatv1.GreetRequest],
) (*connect.Response[chatv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&chatv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Chat-Version", "v1")
	return res, nil
}

func (s *ChatServer) SendChat(
	ctx context.Context,
	req *connect.Request[chatv1.SendChatRequest],
) (*connect.Response[chatv1.SendChatResponse], error) {
	log.Println("Request headers: ", req.Header())
	err := AddChat(req.Msg)
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

	discussionId := req.Msg.DiscussionId
	lowPassword := req.Msg.LowPassword
	watchId := addWatchList(discussionId, &WatchInfo{
		Hash:      generateHash(lowPassword, SECRET_KEY),
		StreamRes: streamRes,
	})
	for {
		err := streamRes.Send(&chatv1.GetChatsStreamResponse{Chats: []*chatv1.Chat{}})
		if err != nil {
			watchList[discussionId][watchId] = nil
			return connect.NewError(connect.CodeInternal, err)
		}
		time.Sleep(time.Second * 30)
	}
}

func main() {
	err := InitializeDB("./database.db")
	if err != nil {
		log.Fatal(err)
	}

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
