package channel

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adrianbrad/chat/message"
	"github.com/adrianbrad/chat/trace"
	"github.com/gorilla/websocket"
)

//Channel implements http.Handler
type Channel interface {
	Run()
	ForwardChannel() chan *message.Message
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type channel struct {
	//forward is a channel that holds incoming message
	//incoming messages should be broadcasted to the other channels
	forward chan *message.Message
	//join is a channel for clients wishing to join the channel
	join chan Client
	//leave is a channel for clients withing to leave the channel
	leave chan Client
	// * the join and leave channels exist simply to allow us to safely add and remove clients from the clients map

	//clients holds all current clients in this channel
	clients map[Client]bool
	//tracer will receive trace information of activity in the channel
	tracer trace.Tracer
}

func New() Channel {
	return &channel{
		forward: make(chan *message.Message),
		join:    make(chan Client),
		leave:   make(chan Client),
		clients: make(map[Client]bool),
		tracer:  trace.New(os.Stdout),
	}
}

func (r *channel) Run() {
	for {
		select { //this select statement will run the code for a particular channel when a message is received on that channel, it will only run a code block at a time so we ensure syncronization for the r.clients map
		case client := <-r.join:
			//joining
			r.clients[client] = true
			r.tracer.Trace("New cient joined")
		case client := <-r.leave:
			//leaving
			delete(r.clients, client)
			close(client.SendChannel())
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg.Message))
			//broadcast message to all clients
			for client := range r.clients {
				client.SendChannel() <- msg
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
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println(r.URL.Host)
		fmt.Println(r.Host)
		return true
	},
}

//ServeHTTP is used for upgrading a HTTP connection to websocket, storing the connection,create the client and pass it to the join channel for the current channel
func (r *channel) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//User-Id header was passed before by the auth or client
	// userID := req.Header.Get("User-Id")
	// if userID == "" {
	// 	log.Fatal("Channel.ServeHTTP:", "No user ID found in the request header")
	// 	return
	// }

	//here we have to ensure that the token is in a valid form in the subprotocols header
	socket, err := upgrader.Upgrade(w, req, http.Header{"Sec-WebSocket-Protocol": websocket.Subprotocols(req)})
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket:  socket,
		send:    make(chan *message.Message, messageBufferSize),
		channel: r,
		userData: map[string]interface{}{
			"name":            "temp",
			"canSendMessages": true,
		},
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.Write() //we run the write method in a different thread
	client.Read()     //we keep reading messages in this thread, thus blocking operations and keeping the connection alive
}

func (r *channel) ForwardChannel() chan *message.Message {
	return r.forward
}
