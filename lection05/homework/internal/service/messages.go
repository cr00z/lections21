package service

import (
	"github.com/cr00z/chat/internal/domain"
	"github.com/cr00z/chat/internal/infrastructure/repository/memory"
)

type MessagesService struct {
	repo repository.Messages
}

func NewMessagesService(repo repository.Messages) *MessagesService {
	return &MessagesService{
		repo: repo,
	}
}

func (s MessagesService) GetMessages(userID int64) []domain.Message {
	return s.repo.GetMessages(userID)
}

func (s MessagesService) CreateMessage(message domain.Message) error {
	return s.repo.CreateMessage(message)
}
