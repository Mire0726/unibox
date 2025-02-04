package channel

import (
	"context"
	"fmt"
	"log"

	"github.com/Mire0726/unibox/backend/domain/model"
	"gorm.io/gorm"
)

type channel struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewChannel(db *gorm.DB, logger *log.Logger) *channel {
	return &channel{
		db:     db,
		logger: logger,
	}
}

func (c *channel) Create(ctx context.Context, channel *model.Channel) (*model.Channel, error) {
	if err := c.db.Create(channel).Error; err != nil {
		c.logger.Printf("failed to create channel: %v", err)

		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return channel, nil
}
