package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cr00z/goSimpleChat/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (h Handler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var message domain.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message.FromUserID = userID

	if err := h.service.CreateMessage(message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	renderJSON(w, map[string]string{"status": "ok"})
}

func (h Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := getUserID(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	renderJSON(w, h.service.GetMessages(0))
}

func (h Handler) PostPrivateMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	toUserID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || toUserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message domain.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	message.FromUserID = userID
	message.ToUserID = toUserID

	if err := h.service.CreateMessage(message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	renderJSON(w, map[string]string{"status": "ok"})
}

func (h Handler) GetPrivateMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	renderJSON(w, h.service.GetMessages(userID))
}

// service

func getUserID(r *http.Request) (int64, error) {
	id := r.Context().Value(contextKey("ID"))
	if id == nil {
		return 0, domain.ErrorUserIDNotFound
	}

	userID, ok := id.(int64)
	if !ok {
		return 0, domain.ErrorUserIDInvalidType
	}

	return userID, nil
}
