package channel

type ClientJoinsRooms struct {
	ClientRooms
	HistoryLimit int
}

type ClientRooms struct {
	Client Client
	Rooms  []int
}
