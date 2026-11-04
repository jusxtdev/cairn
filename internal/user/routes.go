package user

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes(r chi.Router) {
	r.Get("/me", h.GetMe)
}