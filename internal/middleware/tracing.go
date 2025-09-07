package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type tctx string

const TraceContextKey tctx = "trace-id"

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := context.WithValue(r.Context(), TraceContextKey, reqID)

		tracer := otel.Tracer("mediary-server")
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
