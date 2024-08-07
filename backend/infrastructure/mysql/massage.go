package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
        INSERT INTO messages (message_id, channel_id, user_id, content, workspace_id, timestamp)
        VALUES (?, ?, ?, ?, ?,?)
    `
	_, err := repo.DB.ExecContext(ctx, query, message.ID, message.ChannelID, message.UserID, message.Content, message.WorkspaceID, message.Timestamp)
	if err != nil {
		fmt.Printf("Error inserting message: %v\n", err)

		return cerror.Wrap(err, "mysql", cerror.WithInternalCode())
	}
	return err
}

func (repo *MessageRepositoryMySQL) ListByWorkspaceID(ctx context.Context, channelID, workspaceID string) ([]*model.Message, error) {
	const query = `
		SELECT message_id, channel_id, user_id, content, workspace_id, timestamp
		FROM messages
		WHERE channel_id = ? AND workspace_id = ?
		ORDER BY timestamp ASC
	`
	rows, err := repo.DB.QueryContext(ctx, query, channelID, workspaceID)
	if err != nil {
		return nil, cerror.Wrap(err, "mysql", cerror.WithInternalCode())
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var message model.Message
		if err := rows.Scan(&message.ID, &message.ChannelID, &message.UserID, &message.Content, &message.WorkspaceID, &message.Timestamp); err != nil {
			return nil, cerror.Wrap(err, "mysql", cerror.WithInternalCode())
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (repo *MessageRepositoryMySQL) ListRecentMessages(ctx context.Context, channelID, workspaceID string, limit int) ([]*model.Message, error) {
	query := `
        SELECT message_id, channel_id, user_id, content, workspace_id, timestamp
        FROM messages
        WHERE channel_id = ? AND workspace_id = ?
        ORDER BY timestamp DESC
        LIMIT ?
    `
	rows, err := repo.DB.QueryContext(ctx, query, channelID, workspaceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Content, &msg.WorkspaceID, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (repo *MessageRepositoryMySQL) GetMessagesSince(ctx context.Context, channelID, workspaceID string, since time.Time) ([]*model.Message, error) {
	query := `
        SELECT message_id, channel_id, user_id, content, workspace_id, timestamp
        FROM messages
        WHERE channel_id = ? AND workspace_id = ? AND timestamp > ?
        ORDER BY timestamp ASC
    `
	rows, err := repo.DB.QueryContext(ctx, query, channelID, workspaceID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Content, &msg.WorkspaceID, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}
