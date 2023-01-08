package repository

import "github.com/cr00z/goSimpleChat/internal/domain"

type Authorization interface {
	CreateUser(domain.User) (int64, error)
	GetUser(domain.User) (domain.User, error)
}

type Messages interface {
	GetMessages(int64) []domain.Message
	CreateMessage(domain.Message) error
}

type Repository struct {
	Authorization
	Messages
}

func New() (Repository, error) {
	r := Repository{
		Authorization: NewAuthMemory(),
		Messages:      NewMessagesMemory(),
	}

	// create a global messaging space
	_, err := r.Authorization.CreateUser(domain.User{ID: 0, Username: "Global", Password: ""})
	if err != nil {
		return Repository{}, err
	}

	return r, nil
}

func (r Repository) MakeDemo() {
	_, _ = r.Authorization.CreateUser(domain.User{ID: 1, Username: "User1", Password: ""})
	_, _ = r.Authorization.CreateUser(domain.User{ID: 2, Username: "User2", Password: ""})
	_ = r.Messages.CreateMessage(domain.Message{ID: 0, FromUserID: 1, ToUserID: 0, Text: "hello all"})
	_ = r.Messages.CreateMessage(domain.Message{ID: 1, FromUserID: 1, ToUserID: 0, Text: "im User1"})
	_ = r.Messages.CreateMessage(domain.Message{ID: 2, FromUserID: 2, ToUserID: 0, Text: "hello, User1"})
	_ = r.Messages.CreateMessage(domain.Message{ID: 3, FromUserID: 2, ToUserID: 1, Text: "im User2"})
}
