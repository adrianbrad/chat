package message

type BroadcastedMessage struct {
	Action   string
	RoomIDs  []int
	UserID   int
	Nickname string
	Content  interface{}
}
