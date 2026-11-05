package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jusxtdev/cairn/internal/config"
	"github.com/jusxtdev/cairn/internal/httpx"
	"github.com/jusxtdev/cairn/internal/users"
)

func getJWTSecret(token *jwt.Token) (any, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return []byte(cfg.JWTSecret), err
}

func AuthenticationMiddleware(userService *users.UserService, next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		// read authorization header OR cookies
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.WriteJSONResponse(w, http.StatusUnauthorized, httpx.Response{
				Status: false,
				Message: "invalid authorization header",
			})
			return 
		}
		parts := strings.Split(authHeader, " ") 
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.WriteJSONResponse(w, http.StatusUnauthorized, httpx.Response{
				Status: false,
				Message: "invalid authorization header",
			})
			return 
		}
		// extract token
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			parts[1],
			claims,
			getJWTSecret,
			jwt.WithExpirationRequired(),
		)
		if err != nil || !token.Valid {
			httpx.WriteJSONResponse(w, http.StatusUnauthorized, httpx.Response{
				Status: false,
				Message: "Invalid token",
			})
			return
		}
		// extract userId
		userId := claims.Subject
		// check if user exists 
		exists, err := userService.CheckUserExistsById(r.Context(), userId)
		if err != nil {
			httpx.WriteInternalServerErr(w)
			return
		}
		if !exists {
			httpx.WriteJSONResponse(w, http.StatusNotFound, httpx.Response{Status: false, Message: "user not found"})
			return
		}

		// next
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}