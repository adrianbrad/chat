package channel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/adrianbrad/chat/messageProcessor"

	"github.com/adrianbrad/chat/message"
	"github.com/adrianbrad/chat/model"
	"github.com/adrianbrad/chat/repository"
	"github.com/adrianbrad/chat/trace"
	"github.com/gorilla/websocket"
)

//Channel extends http.Handler
type Channel interface {
	Run()
	ServeHTTP(http.ResponseWriter, *http.Request)
	JoinRoom() chan ClientRooms
	LeaveRoom() chan ClientRooms
}

type channel struct {
	//receivedMessages is a channel that holds incoming message
	//incoming messages should be broadcasted to the other channels
	messageQueue chan *message.ReceivedMessage
	//joinChannel is a channel for clients wishing to joinChannel the channel
	joinChannel chan Client
	//leaveChannel is a channel for clients withing to leaveChannel the channel
	leaveChannel chan Client
	//joinRoom is a channel for wsConnections wishing to subscribe to a room messages
	joinRoom chan ClientRooms
	//leaveRoom is a channel for wsConnections wishing to unsubscribe from a room messages
	leaveRoom chan ClientRooms
	// * the joinChannel, leaveChannel, joinRoom and leaveRoom channels exist simply to allow us to safely add and remove clients from the clients map

	//clients holds all current clients in this channel
	clients map[Client]bool
	//tracer will receive trace information of activity in the channel
	tracer trace.Tracer
	//repo persists the changes made to the channel
	usersChannelsRepo repository.UsersChannelsRepository
	usersRepo         repository.Repository
	//channel ID holds the current channel id
	channelID int
	//rooms hold references for all the connections in a room RoomID -> A client with a WebsocketConn
	rooms map[int]map[Client]bool
	//messageProcessor is used for handling incoming messages transforming them to broadcasted messages or doing the actions requested by the message
	messageProcessor messageProcessor.MessageProcessor
}

func New(usersChannelsRepo repository.UsersChannelsRepository, channelID int, usersRepo repository.Repository, messageProcessor messageProcessor.MessageProcessor, roomIDs []int) Channel {
	c := &channel{
		messageQueue:      make(chan *message.ReceivedMessage),
		joinChannel:       make(chan Client),
		leaveChannel:      make(chan Client),
		joinRoom:          make(chan ClientRooms),
		leaveRoom:         make(chan ClientRooms),
		clients:           make(map[Client]bool),
		tracer:            trace.New(os.Stdout),
		usersChannelsRepo: usersChannelsRepo,
		channelID:         channelID,
		usersRepo:         usersRepo,
		messageProcessor:  messageProcessor,
		rooms:             make(map[int]map[Client]bool),
	}
	for _, roomID := range roomIDs {
		c.rooms[roomID] = make(map[Client]bool)
	}
	return c
}

func (c *channel) Run() {
	for {
		select { //this select statement will run the code for a particular channel when a message is received on that channel, it will only run a code block at a time so we ensure syncronization for the r.clients map
		case client := <-c.joinChannel:
			userID := client.GetUserID()
			err := c.usersChannelsRepo.AddOrUpdateUserToChannel(userID, c.channelID)
			if err == nil {
				c.clients[client] = true
				c.tracer.Trace("New client joined")
			} else {
				log.Println("Channel.Run -for.select.case.joinChannel: ", err)
				client.Close()
			}
		case client := <-c.leaveChannel:
			//leaving
			delete(c.clients, client)
			close(client.ForwardMessage())
			c.tracer.Trace("Client left")
		case msg := <-c.messageQueue:
			c.tracer.Trace("Message received: ", msg)
			//TODO broadcast message to specified rooms
			//for clients in rooms[roomID] -> client.ForwardMessage() <- msg
			c.broadcastMessage(c.messageProcessor.ProcessMessage(msg))
		case clientRoom := <-c.joinRoom:
			err := c.addClientToRoom(clientRoom)
			if err != nil {
				clientRoom.Client.ForwardMessage() <- c.messageProcessor.ErrorMessage(err.Error())
			}
		case clientRoom := <-c.leaveRoom:
			log.Println(clientRoom)
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
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Channel.ServeHTTP:", "Invalid user ID passed in header")
	}

	user, err := c.usersRepo.GetOne(userID)
	if err != nil {
		log.Println("Canot find user with ID:", userID, "psql err:", err)
	}

	// For the connection to establish we have to send the subprotocol(token) back to the client - it's from websocket specifications
	socket, err := upgrader.Upgrade(w, req, http.Header{"Sec-WebSocket-Protocol": websocket.Subprotocols(req)})
	if err != nil {
		log.Println("Channel.ServeHTTP:", err)
		return
	}

	client := NewClient(socket, make(chan *message.BroadcastedMessage, messageBufferSize), c.messageQueue, user.(model.User), c.joinRoom, c.leaveRoom)
	c.joinChannel <- client
	defer func() {
		c.leaveChannel <- client
	}()

	go client.Write() //we run the write method in a different thread
	client.Read()     //we keep reading messages in this thread, thus blocking operations and keeping the connection alive
}

func (c channel) broadcastMessage(bm *message.BroadcastedMessage) {
	if len(bm.RoomIDs) == 1 && bm.RoomIDs[0] == -1 { //broadcast to all rooms
		for _, room := range c.rooms {
			for client := range room {
				client.ForwardMessage() <- bm
			}
		}
	}

	for _, roomID := range bm.RoomIDs {
		if room, ok := c.rooms[roomID]; ok {
			for clientInRoom := range room {
				clientInRoom.ForwardMessage() <- bm
			}
		}
	}
	c.tracer.Trace(" -- sent to client")
}

func (c *channel) addClientToRoom(clientRoom ClientRooms) (err error) {
	for _, roomID := range clientRoom.Rooms {
		if room, ok := c.rooms[roomID]; ok {
			if _, clientAlreadyInRoom := room[clientRoom.Client]; !clientAlreadyInRoom {
				room[clientRoom.Client] = true
			} else {
				err = fmt.Errorf("Client already in room")
			}
		} else {
			err = fmt.Errorf("Room does not exist %d", roomID)
		}
	}
	return
}

func (c *channel) removeClientFromRoom(client Client, roomID int) error {
	if room, ok := c.rooms[roomID]; ok {
		delete(room, client)
	} else {
		log.Println("Room does not exist")
	}
	return nil
}

func (c channel) JoinRoom() chan ClientRooms {
	return c.joinRoom
}

func (c channel) LeaveRoom() chan ClientRooms {
	return c.leaveRoom
}
