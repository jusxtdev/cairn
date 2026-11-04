package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id uuid.UUID
	Name string
	Email string
	PasswordHash string
	CreatedAt time.Time
	UpdatedAt time.Time
}