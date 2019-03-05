package model

import (
	"time"

	"github.com/adrianbrad/chat/message"
	"github.com/lib/pq"
)

type Message struct {
	ID       int `json:"messageID"`
	Content  interface{}
	RoomIDs  pq.Int64Array
	UserID   int
	Username string
	SentAt   time.Time
	Action   string
}

func NewMessage(receivedMessage *message.ReceivedMessage) *Message {
	return &Message{
		Content: receivedMessage.Content,
		UserID:  receivedMessage.UserID,
		RoomIDs: receivedMessage.RoomIDs,
		Action:  "message",
	}
}
