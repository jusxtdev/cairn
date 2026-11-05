package auth

import (
	"encoding/json"
	"errors"
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
	UserService *users.UserService
}

func NewHandler(db *pgxpool.Pool, cfg *config.Config) *Handler{
	return &Handler{
		cfg,
		&users.UserService{Db: db},
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
		httpx.WriteBadRequest(w)
		return
	}

	// check if user exists (lower case email)
	body.Email = strings.ToLower(body.Email)
	exists, err := h.UserService.CheckUserExistsByEmail(r.Context(), body.Email)
	if err != nil {
		httpx.WriteInternalServerErr(w)
		return
	}
	if exists {
		httpx.WriteJSONResponse(
			w,
			http.StatusConflict,
			httpx.Response{
				Status: false,
				Message: "email already exists",
			})
		return
	}

	// hash password
	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		httpx.WriteInternalServerErr(w)
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
	err = h.UserService.AddUser(r.Context(), u)
	if err != nil {
		httpx.WriteInternalServerErr(w)
		return
	}
	// generate JWT token
	token, err := GenerateJWT(uuid.String(), []byte(h.cfg.JWTSecret), h.cfg.JWTExpiration)
	if err != nil {
		httpx.WriteInternalServerErr(w)
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
		httpx.WriteBadRequest(w)
		return
	}

	// check if user exists (lower case email)
	body.Email = strings.ToLower(body.Email)
	var user users.User
	user, err = h.UserService.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			httpx.WriteJSONResponse(w, http.StatusNotFound, httpx.Response{
				Status: false,
				Message: err.Error(),
			})
		default:
			httpx.WriteInternalServerErr(w)
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
			httpx.WriteJSONResponse(w, http.StatusUnauthorized, httpx.Response{
				Status: false,
				Message: "wrong password",
			})
			return
		}

		httpx.WriteInternalServerErr(w)
		return
	}

	// generate JWT token if password matches
	token, err := GenerateJWT(user.Id.String(), []byte(h.cfg.JWTSecret), h.cfg.JWTExpiration)
	if err != nil {
		httpx.WriteInternalServerErr(w)
		return
	}

	// return JWT token in response
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "User logged in successfully",
		Data: map[string]string{
			"token": token,
		},
	})
}