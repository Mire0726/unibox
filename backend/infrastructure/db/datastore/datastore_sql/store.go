package datastoresql

import (
	"context"
	"log"

	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/Mire0726/unibox/backend/infrastructure/db/datastore"
	"github.com/Mire0726/unibox/backend/infrastructure/db/datastore/datastore_sql/channel"
	"github.com/Mire0726/unibox/backend/infrastructure/db/datastore/datastore_sql/message"
	"github.com/Mire0726/unibox/backend/infrastructure/db/datastore/datastore_sql/workspace"

	"gorm.io/gorm"
)

type Store struct {
	db     *gorm.DB
	store  datastore.ReadWriteStore
	logger *log.Logger
}

func NewStore(db *gorm.DB, logger *log.Logger) *Store {
	return &Store{
		db: db,
		store: &nonTransactionalReadWriteStore{
			db:     db,
			logger: logger,
		},
		logger: logger,
	}
}

func (s *Store) ReadWrite() datastore.ReadWriteStore {
	return s.store
}

func (s *Store) ReadWriteStore() datastore.ReadWriteStore {
	return &nonTransactionalReadWriteStore{
		db:     s.db,
		logger: s.logger,
	}
}

func (s *Store) ReadWriteTransaction(ctx context.Context, f func(context.Context, datastore.ReadWriteStore) error) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	rw := &transactionalReadWriteStore{
		tx:     tx,
		logger: s.logger,
	}

	if err := f(ctx, rw); err != nil {
		s.logger.Printf("failed to execute transaction: %v", err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		s.logger.Printf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}

type nonTransactionalReadWriteStore struct {
	db     *gorm.DB
	logger *log.Logger
}

type transactionalReadWriteStore struct {
	tx     *gorm.DB
	logger *log.Logger
}

func (s *nonTransactionalReadWriteStore) WorkspaceRepository() repository.WorkspaceRepository {
	return workspace.NewWorkspace(s.db, s.logger)
}

func (s *transactionalReadWriteStore) WorkspaceRepository() repository.WorkspaceRepository {
	return workspace.NewWorkspace(s.tx, s.logger)
}

func (s *nonTransactionalReadWriteStore) ChannelRepository() repository.ChannelRepository {
	return channel.NewChannel(s.db, s.logger)
}

func (s *transactionalReadWriteStore) ChannelRepository() repository.ChannelRepository {
	return channel.NewChannel(s.tx, s.logger)
}

func (s *nonTransactionalReadWriteStore) MessageRepository() repository.MessageRepository {
	return message.NewMessage(s.db, s.logger)
}

func (s *transactionalReadWriteStore) MessageRepository() repository.MessageRepository {
	return message.NewMessage(s.tx, s.logger)
}

// func (s *nonTransactionalReadWriteStore) UserRepository() repository.UserRepository {
// 	return workspace.NewWorkspace(s.db, s.logger)
// }

// func (s *transactionalReadWriteStore) UserRepository() repository.UserRepository {
// 	return workspace.NewWorkspace(s.tx, s.logger)
// }
