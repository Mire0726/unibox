package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/infrastructure/db/datastore"
)

type Channel interface {
	Post(ctx context.Context, userID, organizationID, name string) error
}

type ChannelUsecase struct {
	data datastore.Data
	auth AuthUsecase
}

func NewChannelUsecase(data datastore.Data, authUsecase AuthUsecase) *ChannelUsecase {
	return &ChannelUsecase{
		data: data,
		auth: authUsecase,
	}
}

func (uc *ChannelUsecase) Post(ctx context.Context, userID, organizationID, name string) error {
	channel := &model.Channel{
		OrganizationID: organizationID,
		ID:             uuid.New(),
		Name:           name,
	}
	if err := uc.data.ReadWriteStore.Channel().Create(ctx, channel); err != nil {
		return err
	}

	return nil
}
