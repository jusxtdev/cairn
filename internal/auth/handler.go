package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jusxtdev/cairn/internal/httpx"
	"github.com/jusxtdev/cairn/internal/config"
)

type Handler struct{
	db *pgxpool.Pool
	cfg *config.Config
}

func NewHandler(db *pgxpool.Pool, cfg *config.Config) *Handler{
	return &Handler{db, cfg}
}

type RegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}
func (h *Handler)Register(w http.ResponseWriter, r *http.Request) {
	// decode response body
	var body RegisterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		httpx.WriteJSONResponse(w, http.StatusBadRequest, httpx.Response{})
		return
	}

	// check if user exists (lower case email)
	body.Email = strings.ToLower(body.Email)
	var email string
	err = h.db.QueryRow(
		r.Context(),
		"SELECT email FROM users WHERE email = $1",
		body.Email,
	).Scan(&email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// user does not exist, proceed with registration
		} else {
			httpx.WriteJSONResponse(w, http.StatusInternalServerError, httpx.Response{Status: false, Message: fmt.Sprintf("error : %s", err.Error())})
			return
		}
	}
	if email == body.Email {
		httpx.WriteJSONResponse(w, http.StatusConflict, httpx.Response{Status: false, Message: "email already exists"})
		return
	}

	// hash password
	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		httpx.WriteJSONResponse(w, http.StatusInternalServerError, httpx.Response{Status: false, Message: fmt.Sprintf("error : %s", err.Error())})
		return
	}
	// generate UUID
	uuid := uuid.New()
	// create user in db
	query := `
	INSERT INTO users (id, name, email, password)
	VALUES ($1, $2, $3, $4) 
	`
	_, err = h.db.Exec(r.Context(), query, uuid, body.Name, body.Email, hashedPassword)
	if err != nil {
		httpx.WriteJSONResponse(w, http.StatusInternalServerError, httpx.Response{
			Status: false,
			Message: fmt.Sprintf("error : %s", err.Error()),
		})
		return
	}
	// generate JWT token
	token, err := GenerateJWT(uuid.String(), []byte(h.cfg.JWTSecret), h.cfg.JWTExpiration)
	if err != nil {
		httpx.WriteJSONResponse(w, http.StatusInternalServerError, httpx.Response{
			Status: false,
			Message: fmt.Sprintf("error : %s", err.Error()),
		})
		return
	}

	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "User registered successfully",
		Data: map[string]string{
			"token": token,
		},
	})
}

func (h *Handler)Login(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "User logged in successfully",
	})
}