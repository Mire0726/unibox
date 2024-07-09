package repository

import (
	"context"
	"time"

	"github.com/Mire0726/unibox/backend/domain/model"
)

type MessageRepository interface {
	Insert(ctx context.Context, message *model.Message) error
	ListByWorkspaceID(ctx context.Context, channelID, workspaceID string) ([]*model.Message, error)
	ListRecentMessages(ctx context.Context, channelID, workspaceID string, limit int) ([]*model.Message, error)
	GetMessagesSince(ctx context.Context, channelID, workspaceID string, since time.Time) ([]*model.Message, error)
}

type MessageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}
