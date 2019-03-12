package model

import (
	"time"

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
