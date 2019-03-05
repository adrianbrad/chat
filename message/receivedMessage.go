package message

type ReceivedMessage struct {
	UserID  int
	Content string
	Action  string
	RoomIDs []int64
}
