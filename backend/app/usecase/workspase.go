package usecase

import (
	"context"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/infrastructure/db/datastore"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, workspace *model.Workspace) error
}

type WorkspaceUsecase struct {
	data datastore.Data
	Auth AuthUsecase
}

func NewWorkspaceUsecase(data datastore.Data, authUsecase AuthUsecase) *WorkspaceUsecase {
	return &WorkspaceUsecase{
		data: data,
		Auth: authUsecase,
	}
}

func (uc *WorkspaceUsecase) CreateWorkspace(ctx context.Context, workspace *model.Workspace) error {

	workspace, err := uc.data.ReadWriteStore().Workspace().Create(ctx, workspace)
	if err != nil {
		return err
	}

	return nil
}

func (uc *WorkspaceUsecase) SignInWorkspace(ctx context.Context,workspaceID, password string) (*model.Workspace, error) {
	workspace, err := uc.data.ReadWriteStore().Workspace().GetByID(ctx, workspaceID, password)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}
