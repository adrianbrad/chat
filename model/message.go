package model

import "time"

type Message struct {
	ID      int
	Content string
	RoomIDs []int
	UserID  int
	SentAt  time.Time
}
