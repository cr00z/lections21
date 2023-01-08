package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cr00z/goSimpleChat/internal/service"
	"github.com/go-chi/chi/v5"
)

type contextKey string

type Handler struct {
	service service.Service
}

func New(s service.Service) *Handler {
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
	restricted.Get("/users/me/messages", h.GetPrivateMessagesHandler)

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
