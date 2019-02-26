package model

type Room struct {
	ID          int
	Name        string
	Description string
	ChannelID   int
	UserIDs     []int
}
