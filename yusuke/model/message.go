package model

import (
	"database/sql"
	"fmt"
)

// Message はメッセージの構造体です
type Message struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
	UserName string `json:"username"`
	// 1-1. ユーザー名を表示しよう
}

// MessagesAll は全てのメッセージを返します
func MessagesAll(db *sql.DB) ([]*Message, error) {

	// 1-1. ユーザー名を表示しよう
	rows, err := db.Query(`select id, body, username from message`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []*Message
	for rows.Next() {
		m := &Message{}
		// 1-1. ユーザー名を表示しよう
		if err := rows.Scan(&m.ID, &m.Body, &m.UserName); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

// MessageByID は指定されたIDのメッセージを1つ返します
func MessageByID(db *sql.DB, id string) (*Message, error) {
	m := &Message{}

	// 1-1. ユーザー名を表示しよう
	if err := db.QueryRow(`select id, body from message where id = ?`, id).Scan(&m.ID, &m.Body, &m.UserName); err != nil {
		return nil, err
	}

	return m, nil
}

// Insert はmessageテーブルに新規データを1件追加します
func (m *Message) Insert(db *sql.DB) (*Message, error) {
	// 1-2. ユーザー名を追加しよう
	//INSERT INTO テーブル名(カラム1, カラム2, ...) VALUES(値1, 値2, ...);
	res, err := db.Exec(`insert into message (body, username) values (?, ?);`, m.Body, m.UserName)
	fmt.Printf("%#v", m)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:   id,
		Body: m.Body,
		// 1-2. ユーザー名を追加しよう
		UserName: m.UserName,
	}, nil
}

// 1-3. メッセージを編集しよう
// ...
func (m *Message) Update(db *sql.DB, id string) (*Message, error) {
	//UPDATE テーブル名 SET カラム名1 = 値1, カラム名2 = 値2, ...
	res, err := db.Exec(`UPDATE message SET body = ?, username = ? WHERE id = ?`, m.Body, m.UserName, id)
	if err != nil {
		return nil, err
	}
	res_id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Message{
		ID:   res_id,
		Body: m.Body,
		UserName: m.UserName,
	}, nil
}

// 1-4. メッセージを削除しよう
// ...
func Delete(db *sql.DB, id string) ( error) {
	_, err := db.Exec(`DELETE  from message where id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}
