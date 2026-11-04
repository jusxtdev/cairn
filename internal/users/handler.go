package users

import (
	"net/http"

	"github.com/jusxtdev/cairn/internal/httpx"
)

type Handler struct {}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler)GetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{Status: true, Data: userId})
}