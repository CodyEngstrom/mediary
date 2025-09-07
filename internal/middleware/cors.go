package middleware

import "net/http"

func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Add CORS headers here
			next.ServeHTTP(w, r)
		})
	}
}
