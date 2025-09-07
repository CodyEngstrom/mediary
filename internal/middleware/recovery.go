package middleware

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {

					traceID, _ := r.Context().Value(TraceContextKey).(string)

					logger.Error("panic occurred",
						zap.Any("error", rec),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("trace_id", traceID),
					)

					if span := trace.SpanFromContext(r.Context()); span != nil {
						span.RecordError(rec.(error))
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":    "internal server error",
						"trace_id": traceID,
					})
					// http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
