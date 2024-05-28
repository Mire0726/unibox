// domain/model/message.go
package model

type Message struct {
    ID        string
    ChannelID string
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