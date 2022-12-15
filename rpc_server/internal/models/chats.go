package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/matorix/chat.fdbm.dev/internal/config"
	"github.com/matorix/chat.fdbm.dev/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

type ChatModel struct{}

type Chat struct {
	Id           string
	DiscussionId string
	IpAddress    string
	Uid          string
	Name         string
	Body         string
	Hash         string
	CreatedAt    int32
}

type UserVisibleChat struct {
	Id           string
	DiscussionId string
	Uid          string
	Name         string
	Body         string
	CreatedAt    int32
}

func (model *ChatModel) Init() error {
	db, err := sql.Open("sqlite3", config.DatabaseName)
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

func (model *ChatModel) Create(discussionId, name, body, ipAddress, lowPassword string) (*Chat, error) {
	if len(body) > config.BodyLengthLimit {
		return nil, fmt.Errorf("error: %s", "Body is too long.")
	} else if len(name) > config.NameLengthLimit {
		return nil, fmt.Errorf("error: %s", "Name is too long.")
	}

	createdAt := int32(time.Now().UnixMilli())
	hash := utils.GenerateHash(lowPassword, config.SecretKey)
	uid := utils.GenerateUserID(ipAddress, config.SecretKey)
	id, err := model.ExecInsert(discussionId, ipAddress, uid, name, body, hash, createdAt)
	if err != nil {
		return nil, err
	}
	return &Chat{
		Id:           id,
		DiscussionId: discussionId,
		IpAddress:    ipAddress,
		Uid:          uid,
		Name:         name,
		Body:         body,
		Hash:         hash,
		CreatedAt:    createdAt,
	}, nil
}

func (model *ChatModel) ExecInsert(discussionId, ipAddress, uid, name, body, hash string, created_at int32) (string, error) {
	db, err := sql.Open("sqlite3", config.DatabaseName)
	if err != nil {
		return "", err
	}
	defer db.Close()

	sql := `insert into chats(
            discussion_id,
			ip_address,
			user_id,
			name,
			message,
			hash,
			created_at
		) values (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return "", err
	}
	result, err := stmt.Exec(
		discussionId,
		ipAddress,
		uid,
		name,
		body,
		hash,
		created_at,
	)
	if err != nil {
		return "", err
	}
	chatId, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(chatId)), nil
}

func (model *ChatModel) List(discussionId, lowPassword string, cursor *string, limit int, earlierAt int32) ([]*UserVisibleChat, error) {
	var rows *sql.Rows
	var err error

	if limit > config.DatabaseListLimit {
		return nil, fmt.Errorf("error: limit over")
	}

	hash := utils.GenerateHash(lowPassword, config.SecretKey)
	if cursor == nil || *cursor == "" {
		rows, err = model.ExecListQuery(discussionId, hash, limit)
	} else {
		rows, err = model.ExecListQueryWithOffset(discussionId, hash, limit, *cursor, earlierAt)
	}

	if err != nil {
		return nil, err
	}
	chats := []*UserVisibleChat{}
	defer rows.Close()
	for rows.Next() {
		chat := &UserVisibleChat{}
		if err := rows.Scan(&chat.Id, &chat.Uid, &chat.Name, &chat.Body, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chat.DiscussionId = discussionId
		chats = append(chats, chat)
	}
	return chats, nil
}

func (model *ChatModel) ExecListQuery(discussionId, hash string, limit int) (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", config.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("error: database connection refused")
	}
	defer db.Close()
	listSql := `select id, user_id, name, message, created_at from chats 
	where discussion_id = ? and hash = ? 
	order by created_at desc, id asc
	limit ?;`
	stmt, err := db.Prepare(listSql)
	if err != nil {
		return nil, err
	}
	return stmt.Query(discussionId, hash, limit)
}

func (model *ChatModel) ExecListQueryWithOffset(discussionId, hash string, limit int, cursor string, ealierAt int32) (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", config.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("error: database connection refused")
	}
	defer db.Close()
	listSql := `select id, user_id, name, message, created_at from chats 
	where discussion_id = ? and hash = ? and (created_at < ? or (created_at <= ? and id < ?)) 
	order by created_at desc, id desc
	limit ?;`
	stmt, err := db.Prepare(listSql)
	if err != nil {
		return nil, err
	}
	return stmt.Query(discussionId, hash, ealierAt, ealierAt, cursor, limit)
}
