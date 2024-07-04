package usecase

import (
	"context"
	"errors"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
	"github.com/Mire0726/unibox/backend/pkg/log"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, workspace *model.Workspace) error
}

type WorkspaceUsecase struct {
	WorkspaceRepo repository.WorkspaceRepository
	Auth          AuthUsecase
}

func NewWorkspaceUsecase(workspaceRepo repository.WorkspaceRepository, authUsecase AuthUsecase) *WorkspaceUsecase {
	return &WorkspaceUsecase{
		WorkspaceRepo: workspaceRepo,
		Auth:          authUsecase,
	}
}

func (uc *WorkspaceUsecase) CreateWorkspace(ctx context.Context, userID string, workspace *model.Workspace) error {
	if uc.WorkspaceRepo == nil {
		log.Error("WorkspaceRepo is not initialized")
		return errors.New("internal server error")
	}

	if err := uc.WorkspaceRepo.Insert(ctx, userID, workspace); err != nil {
		return err
	}

	return nil
}

func (uc *WorkspaceUsecase) SighnInWorkspace(ctx context.Context, userID string, workspaceID, password string) (*model.Workspace, error) {
	if uc.WorkspaceRepo == nil {
		log.Error("WorkspaceRepo is not initialized")
		return nil, errors.New("internal server error")
	}

	workspace, err := uc.WorkspaceRepo.FindByID(ctx, workspaceID, password)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}
