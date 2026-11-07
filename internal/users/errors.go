package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jusxtdev/cairn/internal/httpx"
)

var ErrUserNotFound = errors.New("user not found")

func fail(w http.ResponseWriter, err error) {
	switch  {
	case errors.Is(err, ErrUserNotFound):
		slog.Warn(err.Error())
		httpx.WriteJSONResponse(w, http.StatusNotFound, httpx.Response{Status: false, Message: err.Error()})
	default:
		slog.Error(err.Error())
		httpx.WriteJSONResponse(w, http.StatusInternalServerError, httpx.Response{Status: false, Message: "internal server error"})
	}
}