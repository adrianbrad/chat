package channel

import (
	"log"
	"strconv"

	"github.com/adrianbrad/chat/message"
	"github.com/adrianbrad/chat/model"
	"github.com/gorilla/websocket"
)

//Client is a Websocket client opened by a user, it also holds information about that user
type Client interface {
	Read()
	Write()
	ForwardMessage() chan *message.BroadcastedMessage
	GetUserID() int
	Close()
}

type client struct {
	//socket is the websocket connection for this client
	socket *websocket.Conn
	//forwardMessage is the buffered channel on which messages are queued ready to be forwarded to the user browser
	forwardMessage chan *message.BroadcastedMessage
	//channel is the channel this client is chatting in, used to broadcast messages to everyone else in the channel
	channelMessageQueue chan ClientMessage
	//join
	join chan ClientJoinsRooms
	//leave
	leave chan ClientRooms
	//userData is used for storing information about the user
	user model.User
}

func NewClient(
	socket *websocket.Conn,
	forwardMessage chan *message.BroadcastedMessage,
	incomingMessage chan ClientMessage, user model.User,
	join chan ClientJoinsRooms, leave chan ClientRooms) Client {
	return &client{
		socket:              socket,
		forwardMessage:      forwardMessage,
		channelMessageQueue: incomingMessage,
		user:                user,
		join:                join,
		leave:               leave,
	}
}

//We process the message here. This is the first place they reach
func (client *client) Read() {
	defer client.socket.Close()

	for {
		var receivedMessage *message.ReceivedMessage
		err := client.socket.ReadJSON(&receivedMessage)
		log.Println(err)
		//if reading from socket fails the for loop is broken and the socket is closed
		if err != nil {
			return
		}

		switch receivedMessage.Action {
		case "join":
			historyLimit := 30
			if givenHistoryLimit, err := strconv.Atoi(receivedMessage.Content); err == nil {
				historyLimit = givenHistoryLimit
			}

			client.join <- ClientJoinsRooms{
				ClientRooms: ClientRooms{
					Client: client,
					Rooms:  receivedMessage.RoomIDs},
				HistoryLimit: historyLimit}
		case "leave":
			client.leave <- ClientRooms{client, receivedMessage.RoomIDs}
		case "message":
			client.channelMessageQueue <- ClientMessage{client, receivedMessage}
		}
	}
}

func (client *client) Write() {
	defer client.socket.Close()

	for msg := range client.ForwardMessage() {
		err := client.socket.WriteJSON(msg)
		//if writing to socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
	}
}

func (client client) processMessage(rm *message.ReceivedMessage) *message.BroadcastedMessage {
	//TODO return a message to be broadcasted over the specified rooms
	return nil
}

func (client client) sendHistory(rm *message.ReceivedMessage) *message.BroadcastedMessage {
	//TODO get history from the channel for the room he asked to join based on the amount of messages given
	return nil
}

func (client *client) ForwardMessage() chan *message.BroadcastedMessage {
	return client.forwardMessage
}

func (client *client) GetUserID() int {
	return client.user.ID
}

func (client *client) Close() {
	client.socket.Close()
}
