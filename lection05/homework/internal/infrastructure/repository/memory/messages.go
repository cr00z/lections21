package repository

import "github.com/cr00z/goSimpleChat/internal/domain"

type MessagesMemory struct {
	Messages []domain.Message
}

func NewMessagesMemory() *MessagesMemory {
	return &MessagesMemory{Messages: make([]domain.Message, 0)}
}

func (r *MessagesMemory) CreateMessage(message domain.Message) error {
	// INSERT ...
	// TODO: check user id
	message.ID = int64(len(r.Messages))
	r.Messages = append(r.Messages, message)
	return nil
}

func (r MessagesMemory) GetMessages(userID int64) []domain.Message {
	// SELECT ... WHERE to_user = ...
	var result []domain.Message
	for _, msg := range r.Messages {
		if msg.ToUserID == userID {
			result = append(result, msg)
		}
	}
	return result
}
