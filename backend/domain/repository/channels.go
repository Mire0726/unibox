package repository

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/model"
)

type ChannelRepository interface {
	Insert(ctx context.Context, channel *model.Channel) error
}

type ChannelService struct {
	repo ChannelRepository
}

func NewChannelService(repo ChannelRepository) *ChannelService {
	return &ChannelService{repo: repo}
}