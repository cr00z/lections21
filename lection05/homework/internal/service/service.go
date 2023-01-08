package service

import (
	"github.com/cr00z/goSimpleChat/internal/domain"
	repository "github.com/cr00z/goSimpleChat/internal/infrastructure/repository/memory"
)

type Authorization interface {
	CreateUser(domain.User) (int64, error)
	GenerateJWT(domain.User) (string, error)
	ParseJWT(string) (int64, error)
}

type Messages interface {
	GetMessages(int64) []domain.Message
	CreateMessage(domain.Message) error
}

type Service struct {
	Authorization
	Messages
}

func New(repo repository.Repository) Service {
	return Service{
		Authorization: NewAuthService(repo.Authorization),
		Messages:      NewMessagesService(repo.Messages),
	}
}
