package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/google/uuid"
)

type Message interface {
	Post(ctx context.Context, userID, channelID, content string) error
	CreateMessagee(ctx context.Context, rawMsg []byte) error
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

func (uc *MessageUsecase) Post(ctx context.Context, userID, channelID, content string) error {
	message := &model.Message{
		ID:        uuid.New(),
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := uc.MessageRepo.Insert(ctx, message); err != nil {
		return err
	}

	return nil
}

func (uc *MessageUsecase) CreateMessage(ctx context.Context, rawMsg []byte) error {
	var msg model.Message

	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return err
	}

	message := &model.Message{
		ID:        uuid.New(),
		ChannelID: msg.ChannelID,
		UserID:    msg.UserID,
		Content:   msg.Content,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := uc.MessageRepo.Insert(ctx, message); err != nil {
		return err
	}

	uc.Hub.Broadcast<- rawMsg

	return nil
}
