package model

import "github.com/google/uuid"

type Channel struct {
	OrganizationID string
	ID             uuid.UUID
	Name           string
}

func NewChannel(organizationID, id, name string) *Channel {
	return &Channel{
		Name:	name,
	}
}
