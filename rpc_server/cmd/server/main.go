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
	"strings"
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

const SECRET_KEY string = "SECRET_KEY"
const BODY_CHARS_LIMIT int = 10000
const NAME_CHARS_LIMIT int = 10000

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
			ip_address string not null,
			user_id string not null,
			name string not null,
			message string not null,
			hash string not null,
			created_at integer not null
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

func generateUserID(ip_address, salt string) string {
	return generateHash(ip_address, salt)
}

func AddChat(chat *chatv1.SendChatRequest, ip_address string) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	if len(chat.Body) > BODY_CHARS_LIMIT {
		return fmt.Errorf("error: %s", "Body is too long.")
	} else if len(chat.Name) > NAME_CHARS_LIMIT {
		return fmt.Errorf("error: %s", "Name is too long.")
	}

	utcTime := int32(time.Now().UnixMilli())
	hash := generateHash(chat.DiscussionInfo.LowPassword, SECRET_KEY)
	user_id := generateUserID(ip_address, SECRET_KEY)

	insertSql := `insert into chats(
            discussion_id,
			ip_address,
			user_id,
			name,
			message,
			hash,
			created_at
		) values (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(
		chat.DiscussionInfo.Id,
		ip_address,
		user_id,
		chat.Name,
		chat.Body,
		hash,
		utcTime,
	)
	if err != nil {
		return err
	}
	chatId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	announceToWatchList(chat.DiscussionInfo.Id, hash, &chatv1.Chat{
		Id:        strconv.Itoa(int(chatId)),
		UserId:    user_id,
		Name:      chat.Name,
		Body:      chat.Body,
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
	listSql := `select id, user_id, name, message, created_at from chats where discussion_id = ? and hash = ?`
	stmt, err := db.Prepare(listSql)
	if err != nil {
		return nil, err
	}
	// rows, err := stmt.Query()
	hash := generateHash(info.DiscussionInfo.LowPassword, SECRET_KEY)
	rows, err := stmt.Query(info.DiscussionInfo.Id, hash)
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
		if err := rows.Scan(&chat.Id, &chat.UserId, &chat.Name, &chat.Body, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func addWatchList(discussion_id string, info *WatchInfo) string {
	id, _ := uuid.NewUUID()
	watchId := id.String()
	if _, ok := watchList[discussion_id]; !ok {
		watchList[discussion_id] = map[string]*WatchInfo{}
	}
	watchList[discussion_id][watchId] = info
	return watchId
}

func announceToWatchList(discussion_id, hash string, newChat *chatv1.Chat) {
	if _, ok := watchList[discussion_id]; !ok {
		return
	}
	for watchId, info := range watchList[discussion_id] {
		go func(info *WatchInfo, watchId string) {
			if info != nil && info.Hash == hash {
				err := info.StreamRes.Send(&chatv1.GetChatsStreamResponse{Chat: newChat})
				if err != nil {
					watchList[discussion_id][watchId] = nil
				}
			}
		}(info, watchId)
	}
}

// grpc server

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
	lowPassword := req.Msg.DiscussionInfo.LowPassword
	addWatchList(discussionId, &WatchInfo{
		Hash:      generateHash(lowPassword, SECRET_KEY),
		StreamRes: streamRes,
	})
	for {

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
