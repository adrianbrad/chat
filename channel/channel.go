package channel

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/adrianbrad/chat/message"
	"github.com/adrianbrad/chat/repository"
	"github.com/adrianbrad/chat/trace"
	"github.com/gorilla/websocket"
)

//Channel implements http.Handler
type Channel interface {
	ChannelManager
	Run()
	ServeHTTP(http.ResponseWriter, *http.Request)

	MessageQueue() chan *message.BroadcastedMessage
	JoinRoom() chan Client
	LeaveRoom() chan Client
}

type ChannelManager interface {
	AddToRooms(int, Client)
	RemoveFromRooms(int, Client) error
}

type channel struct {
	//receivedMessages is a channel that holds incoming message
	//incoming messages should be broadcasted to the other channels
	messageQueue chan *message.BroadcastedMessage
	//joinChannel is a channel for clients wishing to joinChannel the channel
	joinChannel chan Client
	//leaveChannel is a channel for clients withing to leaveChannel the channel
	leaveChannel chan Client
	//joinRoom is a channel for wsConnections wishing to subscribe to a room messages
	joinRoom chan Client
	//leaveRoom is a channel for wsConnections wishing to unsubscribe from a room messages
	leaveRoom chan Client
	// * the joinChannel, leaveChannel, joinRoom and leaveRoom channels exist simply to allow us to safely add and remove clients from the clients map

	//clients holds all current clients in this channel
	clients map[Client]bool
	//tracer will receive trace information of activity in the channel
	tracer trace.Tracer
	//repo persists the changes made to the channel
	usersChannelsRepo repository.UsersChannelsRepository
	//channel ID holds the current channel id
	channelID int
	//rooms hold references for all the connections in a room RoomID -> A client with a WebsocketConn
	rooms map[int][]Client
}

func New(repo repository.UsersChannelsRepository, channelID int) Channel {
	return &channel{
		messageQueue:      make(chan *message.BroadcastedMessage),
		joinChannel:       make(chan Client),
		leaveChannel:      make(chan Client),
		clients:           make(map[Client]bool),
		tracer:            trace.New(os.Stdout),
		usersChannelsRepo: repo,
		channelID:         channelID,
	}
}

func (c *channel) Run() {
	for {
		select { //this select statement will run the code for a particular channel when a message is received on that channel, it will only run a code block at a time so we ensure syncronization for the r.clients map
		case client := <-c.joinChannel:
			userID := client.GetUserID()
			err := c.usersChannelsRepo.AddOrUpdateUserToChannel(userID, c.channelID)
			if err != nil {
				log.Println(err)
				return
			}
			c.clients[client] = true
			c.tracer.Trace("New client joined")
		case client := <-c.leaveChannel:
			//leaving
			delete(c.clients, client)
			close(client.ForwardMessage())
			c.tracer.Trace("Client left")
		case msg := <-c.messageQueue:
			c.tracer.Trace("Message received: ", msg)

			//TODO broadcast message to specified rooms
			//for clients in rooms[roomID] -> client.ForwardMessage() <- msg

			for client := range c.clients {
				client.ForwardMessage() <- msg
				c.tracer.Trace(" -- sent to client")
			}
		case joiningConn := <-c.joinRoom:
			log.Println(joiningConn)
		case leavingConn := <-c.leaveRoom:
			log.Println(leavingConn)
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
		//TODO accept connections only from the correct locator for the channel
		log.Println(r.URL.Host)
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
	_, err := strconv.Atoi(userIDstr)
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
		socket:         socket,
		forwardMessage: make(chan *message.BroadcastedMessage, messageBufferSize),
		channel:        c,
	}
	c.joinChannel <- client
	defer func() {
		c.leaveChannel <- client
	}()

	go client.Write() //we run the write method in a different thread
	client.Read()     //we keep reading messages in this thread, thus blocking operations and keeping the connection alive
}

func (c channel) MessageQueue() chan *message.BroadcastedMessage {
	return c.messageQueue
}

func (c *channel) AddToRooms(roomID int, client Client) {
	c.rooms[roomID] = append(c.rooms[roomID], client)
}

func (c *channel) RemoveFromRooms(roomID int, client Client) (err error) {
	return nil
}

func (c channel) JoinRoom() chan Client {
	return c.joinRoom
}

func (c channel) LeaveRoom() chan Client {
	return c.leaveRoom
}
