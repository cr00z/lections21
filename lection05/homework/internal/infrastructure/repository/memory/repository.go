package repository

import "github.com/cr00z/chat/internal/domain"

type Authorization interface {
	CreateUser(domain.User) (int64, error)
}

type Messages interface {
	GetMessages() []domain.Message
	CreateMessage(domain.Message) error 
}

type Repository struct {
	Authorization
	Messages
}

func New() Repository {
	r := Repository{
		Authorization: NewAuthMemory(),
		Messages:      NewMessagesMemory(),
	}

	// create a global messaging space
	r.Authorization.CreateUser(domain.User{ID: 0, Username: "Global", Password: ""})

	// demo
	r.Authorization.CreateUser(domain.User{ID: 1, Username: "User1", Password: ""})
	r.Authorization.CreateUser(domain.User{ID: 2, Username: "User2", Password: ""})
	r.Messages.CreateMessage(domain.Message{ID: 0, FromUserID: 1, ToUserID: 0, Text: "hello all"})
	r.Messages.CreateMessage(domain.Message{ID: 1, FromUserID: 1, ToUserID: 0, Text: "im User1"})
	r.Messages.CreateMessage(domain.Message{ID: 2, FromUserID: 2, ToUserID: 0, Text: "hello, User1"})
	r.Messages.CreateMessage(domain.Message{ID: 3, FromUserID: 2, ToUserID: 1, Text: "im User2"})

	return r
}
