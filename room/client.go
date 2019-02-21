package room

import (
	"github.com/adrianbrad/chat/message"
	"github.com/gorilla/websocket"
)

type Client interface {
	Read()
	Write()
	SendChannel() chan *message.Message
}

type client struct {
	//socket is the websocket connection for this client
	socket *websocket.Conn
	//send is the buffered channel on which messages are queued ready to be forwarded to the user browser
	send chan *message.Message
	//room is the room this client is chatting in, used to broadcast messages to everyone else in the room
	room Room
	//userData is used for storing information about the user
	userData map[string]interface{}
}

func (c *client) Read() {
	defer c.socket.Close()

	for {
		var msg *message.Message
		err := c.socket.ReadJSON(&msg)

		//if reading from socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
		msg.Name = c.userData["name"].(string)

		if c.userData["canSendMessages"].(bool) {
			c.room.ForwardChannel() <- msg
		}
	}
}

func (c *client) Write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		//if writing to socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
	}
}

func (c *client) SendChannel() chan *message.Message {
	return c.send
}
