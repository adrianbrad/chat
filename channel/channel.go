package channel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/adrianbrad/chat/message"
	"github.com/adrianbrad/chat/trace"
	"github.com/gorilla/websocket"
)

//Channel implements http.Handler
type Channel interface {
	Run()
	GetForwardChannel() chan *message.Message
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
	clients map[int]Client
	//tracer will receive trace information of activity in the channel
	tracer trace.Tracer
	//channelRepo persists the changes made to the channel
	// channelRepo repository.
}

func New() Channel {
	return &channel{
		forward: make(chan *message.Message),
		join:    make(chan Client),
		leave:   make(chan Client),
		clients: make(map[int]Client),
		tracer:  trace.New(os.Stdout),
	}
}

func (c *channel) Run() {
	for {
		select { //this select statement will run the code for a particular channel when a message is received on that channel, it will only run a code block at a time so we ensure syncronization for the r.clients map
		case client := <-c.join:
			//joining
			c.clients[client.GetUserID()] = client
			c.tracer.Trace("New cient joined")
		case client := <-c.leave:
			//leaving
			delete(c.clients, client.GetUserID())
			close(client.GetSendChannel())
			c.tracer.Trace("Client left")
		case msg := <-c.forward:
			c.tracer.Trace("Message received: ", string(msg.Message))
			//broadcast message to all clients
			for id := range c.clients {
				c.clients[id].GetSendChannel() <- msg
				c.tracer.Trace(" -- sent to client")
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
func (c *channel) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// User-Id header was passed before by the auth or client
	userIDstr := req.Header.Get("User-Id")
	if userIDstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Channel.ServeHTTP:", "No user ID found in the request header")
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Channel.ServeHTTP:", "Invalid user ID passed in header")
	}

	// For the connection to establish we have to send the subprotocol(token) back to the client - it's from websocket specifications
	socket, err := upgrader.Upgrade(w, req, http.Header{"Sec-WebSocket-Protocol": websocket.Subprotocols(req)})
	if err != nil {
		log.Println("Channel.ServeHTTP:", err)
		return
	}

	client := &client{
		socket:  socket,
		send:    make(chan *message.Message, messageBufferSize),
		channel: c,
		userData: map[string]interface{}{
			"UserID": userID,
		},
	}
	c.join <- client
	defer func() {
		c.leave <- client
	}()

	go client.Write() //we run the write method in a different thread
	client.Read()     //we keep reading messages in this thread, thus blocking operations and keeping the connection alive
}

func (c *channel) GetForwardChannel() chan *message.Message {
	return c.forward
}
