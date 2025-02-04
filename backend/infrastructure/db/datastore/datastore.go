package datastore

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/repository"
)

type Data interface {
	ReadWriteStore() ReadWriteStore

	ReadWriteTransaction(ctx context.Context, f func(context.Context, ReadWriteStore) error) error
}
type ReadWriteStore interface {
	Channel() repository.ChannelRepository
	Workspace() repository.WorkspaceRepository
	Message() repository.MessageRepository
}
