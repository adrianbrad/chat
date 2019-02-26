package channel

import (
	"strconv"

	"github.com/adrianbrad/chat/message"
	"github.com/gorilla/websocket"
)

type Client interface {
	Read()
	Write()
	GetSendChannel() chan *message.Message
	GetUserID() int
}

type client struct {
	//socket is the websocket connection for this client
	socket *websocket.Conn
	//send is the buffered channel on which messages are queued ready to be forwarded to the user browser
	send chan *message.Message
	//channel is the channel this client is chatting in, used to broadcast messages to everyone else in the channel
	channel Channel
	//userData is used for storing information about the user
	userData map[string]interface{}
}

func (client *client) Read() {
	defer client.socket.Close()

	for {
		var msg *message.Message
		err := client.socket.ReadJSON(&msg)

		//if reading from socket fails the for loop is broken and the socket is closed
		if err != nil {
			return
		}
		msg.Name = strconv.Itoa(client.userData["UserID"].(int))

		// if client.userData["canSendMessages"].(bool) {
		client.channel.GetForwardChannel() <- msg
		// }
	}
}

func (client *client) Write() {
	defer client.socket.Close()

	for msg := range client.send {
		err := client.socket.WriteJSON(msg)
		//if writing to socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
	}
}

func (c *client) GetSendChannel() chan *message.Message {
	return c.send
}

func (c *client) GetUserID() int {
	return c.userData["UserID"].(int)
}