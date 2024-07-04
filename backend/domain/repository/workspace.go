package repository

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/model"
)

type WorkspaceRepository interface {
	Insert(ctx context.Context, userID string, workspace *model.Workspace) error
	FindByID(ctx context.Context, workspaceID, password string) (*model.Workspace, error)
}

type WorkspaceService struct {
	repo WorkspaceRepository
}

func NewWorkspaceService(repo WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{repo: repo}
}
