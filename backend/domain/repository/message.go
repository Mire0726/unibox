package repository

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/model"
)

type MessageRepository interface {
	Insert(ctx context.Context, message *model.Message) error
}

type MessageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}
