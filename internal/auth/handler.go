package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/jusxtdev/cairn/internal/config"
	"github.com/jusxtdev/cairn/internal/httpx"
	"github.com/jusxtdev/cairn/internal/users"
)

type Handler struct{
	cfg *config.Config
	UserRepository *users.UserRepository
}

func NewHandler(db *pgxpool.Pool, cfg *config.Config) *Handler{
	return &Handler{
		cfg,
		&users.UserRepository{Db: db},
	}
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
		fail(w, ErrBadRequest)
		return
	}

	// check if user exists (lower case email)
	body.Email = strings.ToLower(body.Email)
	exists, err := h.UserRepository.CheckUserExistsByEmail(r.Context(), body.Email)
	if err != nil {
		fail(w, fmt.Errorf("could not check if user exists : %w", err))
		return
	}
	if exists {
		fail(w, ErrEmailAlreadyExists)
		return
	}

	// hash password
	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		fail(w, fmt.Errorf("could not hash password : %w", err))
		return
	}
	// generate UUID
	uuid := uuid.New()
	// create user in db
	u := users.CreateUserInput{
		Id: uuid,
		Name: body.Name,
		Email: body.Email,
		HashedPassword: hashedPassword,
	}
	err = h.UserRepository.AddUser(r.Context(), u)
	if err != nil {
		fail(w, fmt.Errorf("could not add user : %w", err))
		return
	}
	// generate JWT token
	token, err := GenerateJWT(uuid.String(), []byte(h.cfg.JWTSecret), h.cfg.JWTExpiration)
	if err != nil {
		fail(w, fmt.Errorf("could not generate JWT token : %w", err))
		return
	}
	
	SetJWTCookie(w, token, h.cfg)
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "User registered successfully",
		Data: map[string]string{
			"token": token,
		},
	})
	slog.Info("user registered successfully", "email", body.Email)
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
func (h *Handler)Login(w http.ResponseWriter, r *http.Request) {
	// decode response body
	var body LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&body)
	if err != nil {
		fail(w, ErrBadRequest)
		return
	}

	// check if user exists (lower case email)
	body.Email = strings.ToLower(body.Email)
	var user users.User
	user, err = h.UserRepository.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			fail(w, ErrInvalidCredentials)
			return
		default:
			fail(w, fmt.Errorf("could not get user by email : %w", err))
		}
		return 
	}
	
	// compare password with hash
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(body.Password),
	)
	if err != nil {
		
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			fail(w, ErrInvalidCredentials)
			return
		}
		fail(w, fmt.Errorf("could not compare password hash : %w", err))
		return
	}
	
	// generate JWT token if password matches
	token, err := GenerateJWT(user.Id.String(), []byte(h.cfg.JWTSecret), h.cfg.JWTExpiration)
	if err != nil {
		fail(w, fmt.Errorf("could not generate JWT token : %w", err))
		return
	}

	SetJWTCookie(w, token, h.cfg)

	// return JWT token in response
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "User logged in successfully",
		Data: map[string]string{
			"token": token,
		},
	})
	slog.Info("user logged in successfully", "email", body.Email)
}