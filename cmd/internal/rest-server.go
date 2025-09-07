package internal

import (
	"mediary/config"
	"mediary/internal/httphandler"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewServer(cfg *config.Config, handlers []httphandler.Handler, middlewares ...func(http.Handler) http.Handler) *http.Server {
	r := chi.NewRouter()

	for _, mw := range middlewares {
		r.Use(mw)
	}

	for _, h := range handlers {
		h.Register(r)
	}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("simulated crash")
	})
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
