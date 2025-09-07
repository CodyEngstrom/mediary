package httphandler

import "github.com/go-chi/chi/v5"

type Handler interface {
	Register(r chi.Router)
}
