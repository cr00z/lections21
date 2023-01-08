package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/cr00z/goSimpleChat/internal/domain"
	repository "github.com/cr00z/goSimpleChat/internal/infrastructure/repository/memory"
	"github.com/dgrijalva/jwt-go"
)

// TODO:
const (
	salt      = "replace_me"
	tokenTTL  = time.Hour
	tokenSign = "replace_me"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s AuthService) CreateUser(user domain.User) (int64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

type claims struct {
	jwt.StandardClaims
	UserID int64
}

func (s AuthService) GenerateJWT(user domain.User) (string, error) {
	var err error
	user.Password = generatePasswordHash(user.Password)
	if user, err = s.repo.GetUser(user); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})
	return token.SignedString([]byte(tokenSign))
}

func (s AuthService) ParseJWT(tokenString string) (int64, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrorInvalidSigningMethod
		}
		return []byte(tokenSign), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims{}, keyFunc)
	if err != nil {
		return 0, err
	}

	clms, ok := token.Claims.(*claims)
	if !ok {
		return 0, domain.ErrorInvalidTokenClaims
	}

	return clms.UserID, nil
}

// utils

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
