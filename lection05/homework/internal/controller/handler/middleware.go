package handler

import (
	"context"
	"net/http"
	"strings"
)

func (h Handler) Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if headerParts[1] == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := h.service.Authorization.ParseJWT(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		idCtx := context.WithValue(r.Context(), contextKey("ID"), userID)
		handler.ServeHTTP(w, r.WithContext(idCtx))
	}
	return http.HandlerFunc(fn)
}
