package usecase

import (
	"context"
	"time"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
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
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := uc.messageRepo.Insert(ctx, message); err != nil {
		return err
	}

	return nil
}
