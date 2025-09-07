package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {

					traceID, _ := r.Context().Value(TraceContextKey).(string)

					var err error
					switch x := rec.(type) {
					case error:
						err = x
					case string:
						err = errors.New(x)
					default:
						err = fmt.Errorf("unknown panic: %v", x)
					}

					logger.Error("panic recovered",
						zap.Any("panic", rec),          // raw panic value
						zap.Error(err),                 // normalized error
						zap.String("method", r.Method), // request context
						zap.String("path", r.URL.Path),
						zap.String("trace_id", traceID),
						zap.ByteString("stacktrace", debug.Stack()), // stacktrace
					)

					if span := trace.SpanFromContext(r.Context()); span != nil {
						span.RecordError(err)
						span.SetStatus(codes.Error, err.Error())
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":    "internal server error",
						"trace_id": traceID,
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
