package channel

// import (
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/adrianbrad/chat/messageProcessor"

// 	"github.com/adrianbrad/chat/trace"
// )

// func TestClientJoinsChannel(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()

// 	_ = initMockClients(3, ch)

// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.clients), 3) //All clients were found in the db

// 	cl4err := NewClientMock(1234, ch.joinRoom, ch.leaveRoom, ch.messageQueue)
// 	ch.joinChannel <- cl4err
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.clients), 3) //Client with id 4 is not in the db
// }

// func TestReceiveHistoryWhenJoiningRoom(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()

// 	cls := initMockClients(3, ch)

// 	cls[0].setNextMessageToReadAndRead(cls[0].joinRoomsMessage(5, 1))
// 	time.Sleep(50 * time.Millisecond)

// 	mesRepoMock := NewMessagesRepoMock(6, 5)
// 	var expected [][]interface{}
// 	expected = append(expected, mesRepoMock.GetAllWhere("RoomID", 1, 5))
// 	ch.messageProcessor.HistoryMessage(expected, []int{1})
// 	assertEqual(t, cls[0].messages[0].Content, fmt.Sprintf("%s", expected))

// 	cls[1].setNextMessageToReadAndRead(cls[1].joinRoomsMessage(5, 1, 2))
// 	time.Sleep(100 * time.Millisecond)

// 	expected = nil
// 	expected = append(expected, mesRepoMock.GetAllWhere("RoomID", 1, 5))
// 	expected = append(expected, mesRepoMock.GetAllWhere("RoomID", 2, 5))
// 	fmt.Println(append(expected, mesRepoMock.GetAllWhere("RoomID", 2, 5)))
// 	ch.messageProcessor.HistoryMessage(expected, []int{1, 2})

// 	assertEqual(t, cls[1].messages[0].Content, fmt.Sprintf("%s", expected))
// }

// func TestClientsSuccesfullyJoinRooms(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()

// 	cls := initMockClients(3, ch)

// 	cls[0].setNextMessageToRead(cls[0].joinRoomsMessage(0, 1))
// 	cls[0].Read()
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 1)

// 	cls[1].setNextMessageToRead(cls[1].joinRoomsMessage(0, 1))
// 	cls[1].Read()
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 2)

// 	//Client tries to join multiple rooms at once
// 	cls[2].setNextMessageToRead(cls[2].joinRoomsMessage(0, 1, 2))
// 	cls[2].Read()
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 3)
// 	assertEqual(t, len(ch.rooms[2]), 1)
// }

// func TestClientsSuccesfullyLeaveRooms(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()
// 	cls := initMockClients(3, ch)

// 	for _, client := range cls {
// 		client.setNextMessageToReadAndRead(client.joinRoomsMessage(0, 1, 2, 3))
// 	}
// 	time.Sleep(100 * time.Millisecond)

// 	assertEqual(t, len(ch.rooms[1]), 3)
// 	assertEqual(t, len(ch.rooms[2]), 3)
// 	assertEqual(t, len(ch.rooms[3]), 3)

// 	cls[0].setNextMessageToReadAndRead(cls[0].leaveRoomsMessage(1))
// 	time.Sleep(100 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 2)

// 	cls[1].setNextMessageToReadAndRead(cls[1].leaveRoomsMessage(1, 2))
// 	time.Sleep(100 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 1)
// 	assertEqual(t, len(ch.rooms[2]), 2)

// 	cls[2].setNextMessageToReadAndRead(cls[2].leaveRoomsMessage(1, 2, 3))
// 	time.Sleep(100 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 0)
// 	assertEqual(t, len(ch.rooms[2]), 1)
// 	assertEqual(t, len(ch.rooms[3]), 2)
// }

// func TestClientsUnsuccesfullyLeaveRooms(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()
// 	cls := initMockClients(3, ch)

// 	for _, client := range cls {
// 		client.setNextMessageToReadAndRead(client.joinRoomsMessage(0, 1, 2))
// 	}
// 	time.Sleep(50 * time.Millisecond)
// 	for _, client := range cls {
// 		client.clearMessages()
// 	}

// 	assertEqual(t, len(ch.rooms[1]), 3)
// 	assertEqual(t, len(ch.rooms[2]), 3)
// 	assertEqual(t, len(ch.rooms[3]), 0)

// 	cls[2].setNextMessageToReadAndRead(cls[2].leaveRoomsMessage(712, 8312))
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 3)
// 	assertEqual(t, len(ch.rooms[2]), 3)
// 	assertEqual(t, len(ch.rooms[3]), 0)

// 	assertEqual(t, cls[2].messages[0].Content, fmt.Sprintf("Room does not exist %d\nRoom does not exist %d\n", 712, 8312))

// 	cls[2].clearMessages()
// 	cls[2].setNextMessageToReadAndRead(cls[2].leaveRoomsMessage(1, 3))
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 2)
// 	assertEqual(t, len(ch.rooms[2]), 3)
// 	assertEqual(t, len(ch.rooms[3]), 0)
// 	assertEqual(t, cls[2].messages[0].Content, fmt.Sprintf("Client is not in the room %d\n", 3))
// }

// func TestClientsFailToJoinRooms(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()

// 	cls := initMockClients(1, ch)
// 	//Client wants to join multiple rooms but one does not exist. it joins the possible ones and receives an error about the non existing room
// 	cls[0].setNextMessageToRead(cls[0].joinRoomsMessage(0, 1, 2, 4))
// 	cls[0].Read()
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 1)
// 	assertEqual(t, len(ch.rooms[2]), 1)
// 	assertEqual(t, len(cls[0].messages), 1)
// 	assertEqual(t, cls[0].messages[0].Content, "Room does not exist 4\n")

// 	//Client wants to join a room but it is already in it
// 	cls[0].clearMessages()
// 	assertEqual(t, len(cls[0].messages), 0)
// 	cls[0].setNextMessageToRead(cls[0].joinRoomsMessage(0, 1))
// 	cls[0].Read()
// 	time.Sleep(50 * time.Millisecond)
// 	assertEqual(t, len(ch.rooms[1]), 1)
// 	assertEqual(t, len(cls[0].messages), 1)
// 	assertEqual(t, cls[0].messages[0].Content, "Client already in room 1\n")
// }

// func TestClientsSendMessages(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()

// 	cls := initMockClients(5, ch)
// 	cls[0].setNextMessageToReadAndRead(cls[0].joinRoomsMessage(0, 1))
// 	cls[1].setNextMessageToReadAndRead(cls[1].joinRoomsMessage(0, 2))
// 	cls[2].setNextMessageToReadAndRead(cls[2].joinRoomsMessage(0, 3))
// 	cls[3].setNextMessageToReadAndRead(cls[3].joinRoomsMessage(0, 1))
// 	cls[4].setNextMessageToReadAndRead(cls[4].joinRoomsMessage(0, 2))

// 	time.Sleep(50 * time.Millisecond)
// 	for _, client := range cls {
// 		client.clearMessages()
// 	}

// 	//0 in room 1, 1 in room 2, 2 in room 3, 3 in room 1, 4 in room 2
// 	cls[0].setNextMessageToReadAndRead(cls[0].sendMessageToRooms("To all", -1))
// 	time.Sleep(100 * time.Millisecond)

// 	for _, client := range cls {
// 		assertEqual(t, len(client.messages), 1)
// 		assertEqual(t, client.messages[0].Content, "To all")
// 		assertEqual(t, client.messages[0].UserID, 1)
// 	}

// 	cls = append(cls, NewClientMock(5, ch.joinRoom, ch.leaveRoom, ch.messageQueue))
// 	ch.joinChannel <- cls[5]
// 	//we created a new client but he is not in any room but he can actually send messages to other rooms
// 	cls[5].setNextMessageToReadAndRead(cls[5].sendMessageToRooms("Sending to room 1 and 2", 1, 2))
// 	time.Sleep(100 * time.Millisecond)

// 	assertEqual(t, len(cls[0].messages), 2)
// 	assertEqual(t, len(cls[1].messages), 2)
// 	assertEqual(t, len(cls[3].messages), 2)
// 	assertEqual(t, len(cls[4].messages), 2)

// 	assertEqual(t, cls[0].messages[1].Content, "Sending to room 1 and 2")
// 	assertEqual(t, cls[0].messages[1].UserID, 5)
// 	assertEqual(t, cls[1].messages[1].Content, "Sending to room 1 and 2")
// 	assertEqual(t, cls[1].messages[1].UserID, 5)
// 	assertEqual(t, cls[3].messages[1].Content, "Sending to room 1 and 2")
// 	assertEqual(t, cls[3].messages[1].UserID, 5)
// 	assertEqual(t, cls[4].messages[1].Content, "Sending to room 1 and 2")
// 	assertEqual(t, cls[4].messages[1].UserID, 5)
// }

// func TestClientUnsuccessfullySendMessage(t *testing.T) {
// 	ch := initChannelWithMocks()
// 	go ch.Run()

// 	cls := initMockClients(2, ch)

// 	cls[0].setNextMessageToReadAndRead(cls[0].joinRoomsMessage(0, 1))
// 	cls[1].setNextMessageToReadAndRead(cls[1].joinRoomsMessage(0, 2))
// 	cls[1].setNextMessageToReadAndRead(cls[1].joinRoomsMessage(0, 1))

// 	time.Sleep(50 * time.Millisecond)
// 	for _, client := range cls {
// 		client.clearMessages()
// 	}

// 	cls[0].setNextMessageToReadAndRead(cls[0].sendMessageToRooms("Sending to room 1 and 2 and 7(inexistent)", 1, 2, 7))
// 	time.Sleep(100 * time.Millisecond)
// 	assertEqual(t, len(cls[1].messages), 2) //two messages as he is in room 1 and 2
// 	assertEqual(t, len(cls[0].messages), 2) //two messages: one received by being in room 1 and one error message
// 	assertEqual(t, cls[0].messages[1].Action, "error")
// 	assertEqual(t, cls[0].messages[1].Content, "Room does not exist 7")
// }

// func initChannelWithMocks() *channel {
// 	channelID := 1
// 	return &channel{
// 		messageQueue:      make(chan ClientMessage),
// 		joinChannel:       make(chan Client),
// 		leaveChannel:      make(chan Client),
// 		joinRoom:          make(chan ClientJoinsRooms),
// 		leaveRoom:         make(chan ClientRooms),
// 		clients:           make(map[Client]bool),
// 		tracer:            trace.New(os.Stdout),
// 		usersChannelsRepo: NewUSersChannelsRepoMock(),
// 		channelID:         channelID,
// 		usersRepo:         NewUsersRepoMock(),
// 		messagesRepo:      NewMessagesRepoMock(6, 5),
// 		messageProcessor:  messageProcessor.New(),
// 		rooms: map[int]map[Client]bool{
// 			1: map[Client]bool{},
// 			2: map[Client]bool{},
// 			3: map[Client]bool{},
// 		},
// 	}
// }

// func initMockClients(nrOfClients int, ch *channel) (clients []*clientMock) {
// 	for i := 0; i < nrOfClients; i++ {
// 		cm := NewClientMock(i+1, ch.joinRoom, ch.leaveRoom, ch.messageQueue)
// 		clients = append(clients, cm)
// 		//this happens in channel.ServeHTTP
// 		ch.joinChannel <- cm
// 	}
// 	return
// }

// func assertEqual(t *testing.T, a interface{}, b interface{}) {
// 	if a == b {
// 		return
// 	}
// 	t.Errorf("\nReceived %v (type %v)\nExpected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
// }
