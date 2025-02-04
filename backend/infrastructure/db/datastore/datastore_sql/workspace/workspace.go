package workspace

import (
	"context"
	"fmt"
	"log"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"gorm.io/gorm"
)

type workspace struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewWorkspace(db *gorm.DB, logger *log.Logger) repository.WorkspaceRepository {
	return &workspace{
		db:     db,
		logger: logger,
	}
}

func (w *workspace) Create(ctx context.Context, workspace *model.Workspace) (*model.Workspace, error) {
	if err := w.db.Create(workspace).Error; err != nil {
		w.logger.Printf("failed to create workspace: %v", err)

		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	return workspace, nil
}

func (w *workspace) GetByID(ctx context.Context, workspaceID, password string) (*model.Workspace, error) {
	workspace := &model.Workspace{}
	if err := w.db.Where("id = ? AND password = ?", workspaceID, password).First(workspace).Error; err != nil {
		w.logger.Printf("failed to find workspace: %v", err)

		return nil, fmt.Errorf("failed to find workspace: %w", err)
	}

	return workspace, nil
}
