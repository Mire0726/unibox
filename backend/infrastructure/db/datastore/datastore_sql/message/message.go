package message

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Mire0726/unibox/backend/domain/model"
	"gorm.io/gorm"
)

type message struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewMessage(db *gorm.DB, logger *log.Logger) *message {
	return &message{
		db:     db,
		logger: logger,
	}
}

func (m *message) Create(ctx context.Context, message *model.Message) (*model.Message, error) {
	if err := m.db.Create(message).Error; err != nil {
		m.logger.Printf("failed to create message: %v", err)

		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

func (m *message) ListByWorkspaceID(ctx context.Context, channelID, workspaceID string) ([]*model.Message, error) {
	var messages []*model.Message
	if err := m.db.Where("channel_id = ? AND workspace_id = ?", channelID, workspaceID).Find(&messages).Error; err != nil {
		m.logger.Printf("failed to list messages: %v", err)

		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return messages, nil
}

func (m *message) ListRecentMessages(ctx context.Context, channelID, workspaceID string, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	if err := m.db.Where("channel_id = ? AND workspace_id = ?", channelID, workspaceID).Order("timestamp DESC").Limit(limit).Find(&messages).Error; err != nil {
		m.logger.Printf("failed to list recent messages: %v", err)

		return nil, fmt.Errorf("failed to list recent messages: %w", err)
	}

	return messages, nil
}

func (m *message) GetMessagesSince(ctx context.Context, channelID, workspaceID string, timestamp time.Time) ([]*model.Message, error) {
	var messages []*model.Message
	if err := m.db.Where("channel_id = ? AND workspace_id = ? AND timestamp > ?", channelID, workspaceID, timestamp).Find(&messages).Error; err != nil {
		m.logger.Printf("failed to list messages since: %v", err)

		return nil, fmt.Errorf("failed to list messages since: %w", err)
	}

	return messages, nil
}
