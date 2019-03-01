package messageProcessor

import (
	"github.com/adrianbrad/chat/message"
)

type MessageProcessor interface {
	ProcessMessage(*message.ReceivedMessage) *message.BroadcastedMessage
	ErrorMessage(string) *message.BroadcastedMessage
	HistoryMessage(interface{}) *message.BroadcastedMessage
}

type messageProcessor struct {
}

func New() MessageProcessor {
	return &messageProcessor{}
}

func (mp messageProcessor) ProcessMessage(rm *message.ReceivedMessage) *message.BroadcastedMessage {
	return &message.BroadcastedMessage{
		RoomIDs: rm.RoomIDs,
		UserID:  rm.UserID,
		Content: rm.Content,
	}
}

func (mp messageProcessor) ErrorMessage(err string) *message.BroadcastedMessage {
	return &message.BroadcastedMessage{
		Action:  "error",
		Content: err,
	}
}

func (mp messageProcessor) HistoryMessage(history interface{}) *message.BroadcastedMessage {
	return nil
}
