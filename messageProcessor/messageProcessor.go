package messageProcessor

import (
	"github.com/adrianbrad/chat/model"
)

type MessageProcessor interface {
	// ProcessMessage(*message.ReceivedMessage) *message.BroadcastedMessage
	ErrorMessage(string) *model.Message
	HistoryMessage(interface{}, []int64) *model.Message
}

type messageProcessor struct {
}

func New() MessageProcessor {
	return &messageProcessor{}
}

// func (mp messageProcessor) ProcessMessage(rm *message.ReceivedMessage) *message.BroadcastedMessage {
// 	return &message.BroadcastedMessage{
// 		RoomIDs: rm.RoomIDs,
// 		UserID:  rm.UserID,
// 		Content: rm.Content,
// 		Action:  "message",
// 	}
// }

func (mp messageProcessor) ErrorMessage(err string) *model.Message {
	return &model.Message{
		Action:  "error",
		Content: err,
	}
}

func (mp messageProcessor) HistoryMessage(history interface{}, roomIDs []int64) *model.Message {
	return &model.Message{
		Action:   "history",
		RoomIDs:  roomIDs,
		Content:  history,
		Username: "System",
	}
}
