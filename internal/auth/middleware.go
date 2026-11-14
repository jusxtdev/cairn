package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jusxtdev/cairn/internal/config"
	"github.com/jusxtdev/cairn/internal/requestctx"
	"github.com/jusxtdev/cairn/internal/users"
)

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == ""{
		return "", nil
	}

	parts := strings.Split(authHeader, " ") 
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrInvalidAuthCredentials
	}
	return parts[1], nil
}

func AuthenticationMiddleware(userRepository *users.UserRepository, cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		// read authorization header OR cookies
		var token string
		cookie, err := GetJWTCookie(r)
		if err != nil {
			fail(w, err)
			return
		}

		if cookie != "" {
			token = cookie
		} else {
			// if cookie is not provided, use auth header
			authtoken, err :=getTokenFromHeader(r)
			if err != nil{
				fail(w, err)
				return 
			} 
			if authtoken == "" {
				fail(w, ErrMissingCredentials)
				return
			}
			token = authtoken
		}

		// decode token
		claims := &jwt.RegisteredClaims{}
		decoded, err := jwt.ParseWithClaims(
			token,
			claims,
			func(t *jwt.Token) (any, error) {
				return []byte(cfg.JWTSecret), nil
			},
			jwt.WithExpirationRequired(),
		)
		if err != nil || !decoded.Valid {
			fail(w, err)
			return
		}
		// extract userId
		userId := claims.Subject
		// check if user exists 
		exists, err := userRepository.CheckUserExistsById(r.Context(), userId)
		if err != nil {
			fail(w, err)
			return
		}
		if !exists {
			fail(w, ErrUserNotFound)
			return
		}

		// next
		ctx := requestctx.WithUserId(r.Context(), userId) 
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}