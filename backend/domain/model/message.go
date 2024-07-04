// domain/model/message.go
package model

import "github.com/google/uuid"

type Message struct {
    ID        uuid.UUID
    ChannelID string
    WorkspaceID string
    UserID    string
    Content   string
    Timestamp string
}

func NewMessage(channelID, userID, content string) *Message {
    return &Message{
        ChannelID: channelID,
        UserID:    userID,
        Content:   content,
    }
}