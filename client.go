package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	//socket is the websocket connection for this client
	socket *websocket.Conn
	//send is the buffered channel on which messages are queued ready to be forwarded to the user browser
	send chan *message
	//room is the room this client is chatting in, used to broadcast messages to everyone else in the room
	room *room
	//userData is used for storing information about the user
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)

		//if reading from socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
		// msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		// avatarURL, avatarIsPresent := c.userData["avatar_url"]
		// if avatarIsPresent {
		// 	msg.AvatarURL = avatarURL.(string)
		// }
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		//if writing to socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
	}
}
