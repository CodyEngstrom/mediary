package middleware

import (
	"mediary/config"
	"net/http"
)

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement authentication
			next.ServeHTTP(w, r)
		})
	}
}
