package repository

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/model"
)

type ChannelRepository interface {
	Create(ctx context.Context, channel *model.Channel) (*model.Channel, error)
}
