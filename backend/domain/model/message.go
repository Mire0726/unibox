// domain/model/message.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID `json:"id"`
	ChannelID   string    `json:"channelID"`
	WorkspaceID string    `json:"workspaceID"`
	UserID      string    `json:"userID"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewMessage(channelID, userID, content string) *Message {
	return &Message{
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
	}
}
