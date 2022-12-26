package repository

import "github.com/cr00z/chat/internal/domain"

type MessagesMemory struct {
	Messages []domain.Message
}

func NewMessagesMemory() *MessagesMemory {
	return &MessagesMemory{Messages: make([]domain.Message, 0)}
}

func (r *MessagesMemory) CreateMessage(message domain.Message) error {
	r.Messages = append(r.Messages, message)
	return nil
}

func (r MessagesMemory) GetMessages() []domain.Message {
	return r.Messages
}
