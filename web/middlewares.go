package web

import (
	"instancer/internal/env"
	"net/http"
)

const (
	// XCtfdAuth is the header that CTFd uses to authenticate requests
	XCtfdAuth = "X-Ctfd-Auth"
)

func FromCtfd(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := env.Get()
		token := r.Header.Get(XCtfdAuth)

		if token != cfg.Token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func FromAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := env.Get()
		token := r.Header.Get(XCtfdAuth)

		if token != cfg.GlobalToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
