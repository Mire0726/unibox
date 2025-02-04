package repository

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/model"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *model.Workspace) (*model.Workspace, error)
	GetByID(ctx context.Context, workspaceID, password string) (*model.Workspace, error)
}
