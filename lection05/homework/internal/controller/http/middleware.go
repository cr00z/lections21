package handler

import "net/http"

func (h Handler) Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("auth"))
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
