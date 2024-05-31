package mysql

import (
	"context"
	"database/sql"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/Mire0726/unibox/backend/internal/cerror"
)

type ChannelRepositoryMySQL struct {
	DB *sql.DB
}

func NewChannelRepository(db *sql.DB) repository.ChannelRepository {
	return &ChannelRepositoryMySQL{DB: db}
}

func (repo *ChannelRepositoryMySQL) Insert(ctx context.Context, channel *model.Channel) error {
	const query = `
        INSERT INTO channels (organization_id, channel_id, name)
        VALUES (?, ?, ?)
    `
	_, err := repo.DB.ExecContext(ctx, query, channel.OrganizationID, channel.ID, channel.Name)
	if err != nil {
		return cerror.Wrap(err, "mysql", cerror.WithInternalCode())
	}

	return err
}
