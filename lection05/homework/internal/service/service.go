package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/cr00z/chat/internal/domain"
	"github.com/cr00z/chat/internal/infrastructure/repository/memory"
)

// TODO:
const salt = "replace_me"

type Service struct {
	repo repository.Repository
}

func New(r repository.Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s Service) CreateUser(user domain.User) (int64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.Authorization.CreateUser(user)
}

func (s Service) GetMessages() []domain.Message {
	return s.repo.Messages.GetMessages()
}

// utils

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
