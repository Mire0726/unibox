package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/Mire0726/unibox/backend/pkg/log"
	"github.com/google/uuid"
)

type Message interface {
	CreateMessagee(ctx context.Context, rawMsg []byte) error
	ListMessages(ctx context.Context, channelID string) ([]*model.Message, error)
}

type MessageUsecase struct {
	MessageRepo repository.MessageRepository
	Auth        AuthUsecase
	Hub         *model.Hub
}

func NewMessageUsecase(messageRepo repository.MessageRepository, authUsecase AuthUsecase, hub *model.Hub) *MessageUsecase {
	return &MessageUsecase{
		MessageRepo: messageRepo,
		Auth:        authUsecase,
		Hub:         hub,
	}
}

func (uc *MessageUsecase) CreateMessage(ctx context.Context, userID, channelID, workspaceID string, content string) error {
	if uc.MessageRepo == nil {
		log.Error("MessageRepo is not initialized")
		return errors.New("internal server error")
	}

	if uc.Hub == nil {
		log.Error("Hub is not initialized")
		return errors.New("internal server error")
	}

	message := &model.Message{
		ID:          uuid.New(),
		ChannelID:   channelID,
		WorkspaceID: workspaceID,
		UserID:      userID,
		Content:     content,
		Timestamp:   time.Now(),
	}

	if err := uc.MessageRepo.Insert(ctx, message); err != nil {
		return err
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	uc.Hub.Broadcast <- messageData

	return nil
}

func (uc *MessageUsecase) ListMessages(ctx context.Context, channelID, workspaceID string) ([]*model.Message, error) {
	if uc.MessageRepo == nil {
		log.Error("MessageRepo is not initialized")
		return nil, errors.New("internal server error")
	}

	if uc.Hub == nil {
		log.Error("Hub is not initialized")
		return nil, errors.New("internal server error")
	}

	messages, err := uc.MessageRepo.ListByWorkspaceID(ctx, channelID, workspaceID)
	if err != nil {
		return nil, err
	}

	messageData, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	uc.Hub.Broadcast <- messageData

	return messages, nil
}

func (uc *MessageUsecase) StartRealtimeUpdates(channelID, workspaceID string, interval time.Duration) {
	lastCheck := time.Now()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			messages, err := uc.MessageRepo.GetMessagesSince(context.Background(), channelID, workspaceID, lastCheck)
			if err != nil {
				// エラーハンドリング
				continue
			}
			if len(messages) > 0 {
				for _, msg := range messages {
					messageData, err := json.Marshal(msg)
					if err != nil {
						// エラーハンドリング
						continue
					}
					uc.Hub.Broadcast <- messageData
				}
				lastCheck = messages[len(messages)-1].Timestamp
			}
		}
	}
}
