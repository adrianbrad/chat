package channel

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
}

type channel struct {

	//channel ID holds the current channel id
	channelID int

	//messageQueue is a channel that holds incoming message
	//incoming messages should be broadcasted to the other clients
	messageQueue chan ClientMessage
	//joinChannel is a channel for clients wishing to joinChannel the channel
	joinChannel chan Client
	//leaveChannel is a channel for clients withing to leaveChannel the channel
	leaveChannel chan Client
	//joinRoom is a channel for wsConnections wishing to subscribe to a room messages
	joinRoom chan ClientJoinsRooms
	//leaveRoom is a channel for wsConnections wishing to unsubscribe from a room messages
	leaveRoom chan ClientRooms
	// * the joinChannel, leaveChannel, joinRoom and leaveRoom channels exist simply to allow us to safely add and remove clients from the clients map

	//clients holds all current clients in this channel
	clients map[Client]bool
	//rooms hold references for all the connections in a room RoomID -> A client with a WebsocketConn
	rooms map[int]map[Client]bool

	//repo persists the changes made to the channel
	usersChannelsRepo repository.UsersChannelsRepository
	usersRepo         repository.Repository
	messagesRepo      repository.Repository

	//messageProcessor is used for handling incoming messages transforming them to broadcasted messages or doing the actions requested by the message
	messageProcessor messageProcessor.MessageProcessor

	//tracer will receive trace information of activity in the channel
	tracer trace.Tracer
}

func New(
	usersChannelsRepo repository.UsersChannelsRepository,
	channelID int,
	usersRepo repository.Repository,
	messageProcessor messageProcessor.MessageProcessor,
	roomIDs []int,
	messagesRepo repository.Repository) Channel {
	c := &channel{
		messageQueue: make(chan ClientMessage),
		joinChannel:  make(chan Client),
		leaveChannel: make(chan Client),
		joinRoom:     make(chan ClientJoinsRooms),
		leaveRoom:    make(chan ClientRooms),
		clients:      make(map[Client]bool),
		tracer:       trace.New(os.Stdout),
		rooms:        make(map[int]map[Client]bool),

		usersChannelsRepo: usersChannelsRepo,
		channelID:         channelID,
		usersRepo:         usersRepo,
		messagesRepo:      messagesRepo,
		messageProcessor:  messageProcessor,
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
			} else {
				log.Println("Channel.Run -for.select.case.joinChannel: ", err)
				client.Close()
			}
		case client := <-c.leaveChannel:
			delete(c.clients, client)
			close(client.ForwardMessage())
		case clientMessage := <-c.messageQueue:
			err := c.broadcastMessage(c.messageProcessor.ProcessMessage(clientMessage.Message))
			if err != nil {
				clientMessage.Client.ForwardMessage() <- c.messageProcessor.ErrorMessage(err.Error())
			}
		case clientRoom := <-c.joinRoom:
			err := c.addClientToRoom(clientRoom)
			if err == nil {
				history := c.getHistory(clientRoom.Rooms, clientRoom.HistoryLimit)
				clientRoom.Client.ForwardMessage() <- c.messageProcessor.HistoryMessage(history, clientRoom.Rooms)
			} else {
				clientRoom.Client.ForwardMessage() <- c.messageProcessor.ErrorMessage(err.Error())

			}
		case clientRoom := <-c.leaveRoom:
			err := c.removeClientFromRoom(clientRoom)
			if err != nil {
				clientRoom.Client.ForwardMessage() <- c.messageProcessor.ErrorMessage(err.Error())
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

func (c channel) broadcastMessage(bm *message.BroadcastedMessage) (err error) {
	var errorMessage strings.Builder

	if len(bm.RoomIDs) == 1 && bm.RoomIDs[0] == -1 { //broadcast to all rooms
		for _, room := range c.rooms {
			for client := range room {
				client.ForwardMessage() <- bm
			}
		}
	} else {
		for _, roomID := range bm.RoomIDs {
			if room, ok := c.rooms[roomID]; ok {
				for clientInRoom := range room {
					clientInRoom.ForwardMessage() <- bm
				}
			} else {
				_, _ = fmt.Fprintf(&errorMessage, "Room does not exist %d", roomID)
			}
		}
	}

	if errorMessage.String() != "" {
		err = fmt.Errorf(errorMessage.String())
	}
	return
}

func (c *channel) addClientToRoom(clientRoom ClientJoinsRooms) (err error) {
	var errorMessage strings.Builder

	for _, roomID := range clientRoom.Rooms {
		if room, ok := c.rooms[roomID]; ok {
			if _, clientAlreadyInRoom := room[clientRoom.Client]; !clientAlreadyInRoom {
				room[clientRoom.Client] = true
			} else {
				_, _ = fmt.Fprintf(&errorMessage, "Client already in room %d\n", roomID)
			}
		} else {
			_, _ = fmt.Fprintf(&errorMessage, "Room does not exist %d\n", roomID)
		}
	}
	if errorMessage.String() != "" {
		err = fmt.Errorf(errorMessage.String())
	}
	return
}

func (c *channel) removeClientFromRoom(clientRoom ClientRooms) (err error) {
	var errorMessage strings.Builder

	for _, roomID := range clientRoom.Rooms {
		if room, ok := c.rooms[roomID]; ok {
			if _, clientInRoom := room[clientRoom.Client]; clientInRoom {
				delete(room, clientRoom.Client)
			} else {
				_, _ = fmt.Fprintf(&errorMessage, "Client is not in the room %d\n", roomID)
			}
		} else {
			_, _ = fmt.Fprintf(&errorMessage, "Room does not exist %d\n", roomID)
		}
	}
	if errorMessage.String() != "" {
		err = fmt.Errorf(errorMessage.String())
	}
	return
}

type History [][]interface{}

func (c channel) getHistory(RoomIDS []int, numberOfMessages int) (history History) {
	for _, roomID := range RoomIDS {
		history = append(history, c.messagesRepo.GetAllWhere("RoomID", roomID, numberOfMessages))
	}
	return
}
