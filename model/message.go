package model

import "time"

type Message struct {
	ID        int
	Content   string
	RoomID    int
	UserID    int
	CreatedAt time.Time
}
