package main

import (
	"log"
	"net/http"

	"github.com/stretchr/objx"

	"github.com/adrianbrad/chat/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	//forward is a channel that holds incoming message
	//incoming messages should be broadcasted to the other channels
	forward chan *message
	//join is a channel for clients wishing to join the room
	join chan *client
	//leave is a channel for clients withing to leave the room
	leave chan *client
	// * the join and leave channels exist simply to allow us to safely add and remove clients from the clients map

	//clients holds all current clients in this room
	clients map[*client]bool
	//tracer will receive trace information of activity in the room
	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select { //this select statement will run the code for a particular channel when a message is received on that channel, it will only run a code block at a time so we ensure syncronization for the r.clients map
		case client := <-r.join:
			//joining
			r.clients[client] = true
			r.tracer.Trace("New cient joined")
		case client := <-r.leave:
			//leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg.Message))
			//broadcast message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

//ServeHTTP is used for upgrading a HTTP connection to websocket, storing the connection,create the client and pass it to the join channel for the current room
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write() //we run the write method in a different thread
	client.read()     //we keep reading messages in this thread, thus blocking operations and keeping the connection alive
}
