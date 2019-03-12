package channel

import "github.com/adrianbrad/chat/model"

type ClientMessage struct {
	Client  Client
	Message *model.Message
}
