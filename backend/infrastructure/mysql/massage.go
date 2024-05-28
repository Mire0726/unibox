package mysql

import (
	"context"
	"database/sql"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
)

// MessageRepositoryMySQL は SQL データベースを使用してメッセージデータを管理するためのリポジトリです。
type MessageRepositoryMySQL struct {
	DB *sql.DB
}

// NewMessageRepository は新しい MessageRepositoryMySQL インスタンスを返します。
func NewMessageRepository(db *sql.DB) repository.MessageRepository {
    return &MessageRepositoryMySQL{DB: db}
}


// SaveMessage はデータベースに新しいメッセージを保存します。
func (repo *MessageRepositoryMySQL) Insert(ctx context.Context, message *model.Message) error {
	const query = `
        INSERT INTO messages (channel_id, user_id, content, timestamp)
        VALUES (?, ?, ?, ?)
    `
	_, err := repo.DB.ExecContext(ctx, query, message.ChannelID, message.UserID, message.Content, message.Timestamp)
	return err
}
