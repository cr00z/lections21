package handler

import "net/http"

func (h Handler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("post message"))
}

func (h Handler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, h.service.GetMessages())
}

func (h Handler) PostPrivateMessageHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("post private message"))
}

func (h Handler) GetPrivateMessagesHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("get private messages"))
}
