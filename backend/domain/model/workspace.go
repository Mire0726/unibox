package model

import "github.com/google/uuid"

type Workspace struct {
	ID       uuid.UUID
	Name     string
	Password string
}

func NewWorkspace(id uuid.UUID, name, password string) *Workspace {
	return &Workspace{
		ID:       id,
		Name:     name,
		Password: password,
	}
}
