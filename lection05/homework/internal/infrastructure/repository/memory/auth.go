package repository

import "github.com/cr00z/chat/internal/domain"

type AuthMemory struct {
	Users []domain.User
}

func NewAuthMemory() *AuthMemory {
	return &AuthMemory{Users: make([]domain.User, 0)}
}

func (r *AuthMemory) CreateUser(user domain.User) (int64, error) {
	// INSERT ... RETURNING id analogue
	for _, u := range r.Users {
		if u.Username == user.Username {
			return 0, domain.ErrorIncorrectUsername
		}
	}
	id := int64(len(r.Users))
	user.ID = id
	r.Users = append(r.Users, user)
	return id, nil
}
