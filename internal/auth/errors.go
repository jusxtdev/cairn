package auth

import "errors"

var (
	ErrBadRequest = errors.New("could not decode request body")
	ErrEmailAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)