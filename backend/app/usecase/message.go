package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/google/uuid"
)

type Message interface {
	Post(ctx context.Context, userID, channelID, content string) error
}

type MessageUsecase struct {
	messageRepo repository.MessageRepository
	auth        AuthUsecase
}

func NewMessageUsecase(messageRepo repository.MessageRepository, authUsecase AuthUsecase) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
		auth:        authUsecase,
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
	fmt.Println(message.Content)

	if err := uc.messageRepo.Insert(ctx, message); err != nil {
		return err
	}

	return nil
}
