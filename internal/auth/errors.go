package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jusxtdev/cairn/internal/httpx"
)

var (
	ErrBadRequest = errors.New("could not decode request body")
	ErrEmailAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrMissingCredentials = errors.New("authentication credentials are required")
    ErrInvalidAuthCredentials = errors.New("invalid authentication credentials")
)

func fail(w http.ResponseWriter, err error){
	switch {
	case errors.Is(err, ErrBadRequest):
		slog.Warn(err.Error())
		httpx.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrEmailAlreadyExists):
		slog.Warn(err.Error())
		httpx.WriteErrorResponse(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrInvalidCredentials):
		slog.Warn(err.Error())
		httpx.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrMissingCredentials):
		slog.Warn(err.Error())
		httpx.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrInvalidAuthCredentials):
		slog.Warn(err.Error())
		httpx.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
	default:
		slog.Error(err.Error())
		httpx.WriteErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}
}