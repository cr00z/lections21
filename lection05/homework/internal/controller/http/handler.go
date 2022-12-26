package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cr00z/chat/internal/domain"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	CreateUser(domain.User) (int64, error)
	GetMessages() []domain.Message
}

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h Handler) InitRoutes() *chi.Mux {
	root := chi.NewRouter()
	root.Post("/register", h.RegisterHandler)
	root.Post("/login", h.LoginHandler)

	restricted := chi.NewRouter()
	restricted.Use(h.Auth)
	restricted.Post("/messages", h.PostMessageHandler)
	restricted.Get("/messages", h.GetMessagesHandler)
	restricted.Post("/users/{id}/messages", h.PostPrivateMessageHandler)
	restricted.Get("/users/{id}/messages", h.GetPrivateMessagesHandler)

	root.Mount("/api", restricted)

	return root
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
