package middleware

import "net/http"

func MetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Record request metrics
			next.ServeHTTP(w, r)
		})
	}
}
