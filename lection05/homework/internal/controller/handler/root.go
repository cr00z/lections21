package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cr00z/goSimpleChat/internal/domain"
)

func (h Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input domain.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateUser(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	input.ID = id
	renderJSON(w, map[string]int64{"id": id})
}

func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input domain.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.service.Authorization.GenerateJWT(input)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	renderJSON(w, map[string]string{"token": token})
}
