package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/Mire0726/unibox/backend/internal/cerror"
)

type MessageRepositoryMySQL struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &MessageRepositoryMySQL{DB: db}
}

func (repo *MessageRepositoryMySQL) Insert(ctx context.Context, message *model.Message) error {
	const query = `
        INSERT INTO messages (message_id, channel_id, user_id, content, timestamp)
        VALUES (?, ?, ?, ?, ?)
    `
	_, err := repo.DB.ExecContext(ctx, query, message.ID, message.ChannelID, message.UserID, message.Content, message.Timestamp)
	if err != nil {
		return cerror.Wrap(err, "mysql", cerror.WithInternalCode())
	}
	fmt.Println("inserted")
	return err
}
