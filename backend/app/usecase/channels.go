package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
)

type Channel interface {
	Post(ctx context.Context, userID, organizationID, name string) error
}

type ChannelUsecase struct {
	channelRepo repository.ChannelRepository
	auth        AuthUsecase
}

func NewChannelUsecase(channelRepo repository.ChannelRepository, authUsecase AuthUsecase) *ChannelUsecase {
	return &ChannelUsecase{
		channelRepo: channelRepo,
		auth:        authUsecase,
	}
}

func (uc *ChannelUsecase) Post(ctx context.Context, userID, organizationID, name string) error {
	channel := &model.Channel{
		OrganizationID: organizationID,
		ID:             uuid.New(),
		Name:           name,
	}
	if err := uc.channelRepo.Insert(ctx, channel); err != nil {
		return err
	}

	return nil
}
