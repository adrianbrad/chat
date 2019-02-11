package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	//socket is the websocket connection for this client
	socket *websocket.Conn
	//send is the buffered channel on which messages are queued ready to be forwarded to the user browser
	send chan []byte
	//room is the room this client is chatting in, used to broadcast messages to everyone else in the room
	room *room
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		//if reading from socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		//if writing to socket fails the for loop is brocken and the socket is closed
		if err != nil {
			return
		}
	}
}
