package repository

import "github.com/cr00z/goSimpleChat/internal/domain"

type AuthMemory struct {
	Users []domain.User
}

func NewAuthMemory() *AuthMemory {
	return &AuthMemory{Users: make([]domain.User, 0)}
}

func (r *AuthMemory) CreateUser(user domain.User) (int64, error) {
	// INSERT ... RETURNING id
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

func (r *AuthMemory) GetUser(user domain.User) (domain.User, error) {
	// SELECT ... WHERE username = '...' AND password = '...'
	for _, u := range r.Users {
		if u.Username == user.Username && u.Password == user.Password {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrorUnknownUser
}
