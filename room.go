package main

import (
	"log"
	"net/http"
	"time"

	"github.com/adrianbrad/chat/trace"
	"github.com/gorilla/websocket"
	cache "github.com/patrickmn/go-cache"
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
	//auth will allow connections to the room
	auth *authenticator
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		auth:    &authenticator{cache.New(10*time.Second, 10*time.Second)},
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
	switch req.Method {
	case http.MethodGet:
		r.handleWS(w, req)
	case http.MethodPost:
		r.auth.authenticate(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (r *room) handleWS(w http.ResponseWriter, req *http.Request) {
	subprotocols := websocket.Subprotocols(req)
	if len(subprotocols) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := subprotocols[0]
	userID, ok := r.auth.tokens.Get(token)

	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	r.auth.tokens.Delete(token)

	socket, err := upgrader.Upgrade(w, req, http.Header{"Sec-WebSocket-Protocol": []string{token}})
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan *message, messageBufferSize),
		room:   r,
		userData: map[string]interface{}{
			"name": users[userID.(string)].Name,
		},
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write() //we run the write method in a different thread
	client.read()     //we keep reading messages in this thread, thus blocking operations and keeping the connection alive
}
