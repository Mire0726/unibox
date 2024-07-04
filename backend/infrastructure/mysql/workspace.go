package mysql

import (
	"context"
	"database/sql"

	"github.com/Mire0726/unibox/backend/domain/model"
	"github.com/Mire0726/unibox/backend/domain/repository"
)

type WorkspaceRepositoryMySQL struct {
	DB *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) repository.WorkspaceRepository {
	return &WorkspaceRepositoryMySQL{DB: db}
}

func (repo *WorkspaceRepositoryMySQL) Insert(ctx context.Context, userID string, workspace *model.Workspace) error {
	const query = `
		INSERT INTO workspaces (id, name, password,createdBy)
		VALUES (?, ?, ?,?)
	`
	_, err := repo.DB.Exec(query, workspace.ID, workspace.Name, workspace.Password, userID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *WorkspaceRepositoryMySQL) FindByID(ctx context.Context, workspaceID, password string) (*model.Workspace, error) {
	const query = `
		SELECT id, name, password
		FROM workspaces
		WHERE id = ? AND password = ?
	`
	row := repo.DB.QueryRow(query, workspaceID, password)

	workspace := &model.Workspace{}
	if err := row.Scan(&workspace.ID, &workspace.Name, &workspace.Password); err != nil {
		return nil, err
	}

	return workspace, nil
}	
