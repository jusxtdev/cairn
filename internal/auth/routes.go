package auth

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
}