package users

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jusxtdev/cairn/internal/httpx"
)

type Handler struct {
	UserRepository *UserRepository
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		&UserRepository{
			Db: db,
		},
	}
}

func (h *Handler)GetMe(w http.ResponseWriter, r *http.Request) {
	userId, ok:= r.Context().Value("userId").(string)
	if !ok {
		fail(w, errors.New("user id missing from request context"))
		return
	}
	user, err := h.UserRepository.GetUserById(r.Context(), userId)
	if err != nil {
		fail(w, err)
		return
	}
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{Status: true, Data: user})
}