package cache

import (
	"sync"
	"time"
)

type MessageCache struct {
	mu       sync.RWMutex
	messages map[string]cachedMessage
}

type cachedMessage struct {
	Message   string
	Timestamp time.Time
}

func NewMessageCache() *MessageCache {
	return &MessageCache{
		messages: make(map[string]cachedMessage),
	}
}

func (c *MessageCache) Set(key string, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages[key] = cachedMessage{Message: message, Timestamp: time.Now()}
}

func (c *MessageCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	msg, found := c.messages[key]
	if !found || time.Since(msg.Timestamp) > 5*time.Minute {
		return "", false
	}
	return msg.Message, true
}

func (c *MessageCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = make(map[string]cachedMessage)
}
