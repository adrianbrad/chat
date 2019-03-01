package channel

import "github.com/adrianbrad/chat/message"

type ClientMessage struct {
	Client  Client
	Message *message.ReceivedMessage
}
