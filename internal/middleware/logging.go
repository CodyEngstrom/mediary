package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			traceID, _ := r.Context().Value(TraceContextKey).(string)

			logFunc := logger.Info
			if rw.statusCode >= 500 {
				logFunc = logger.Error
			} else if rw.statusCode >= 400 {
				logFunc = logger.Warn
			}

			logFunc("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.String("remote_ip", r.RemoteAddr),
				zap.Duration("duration", time.Since(start)),
				zap.String("trace_id", traceID),
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
