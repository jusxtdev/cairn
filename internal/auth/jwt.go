package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(userID string, jwtSecret []byte, jwtExpiration int) (string, error) {
	jwtExpirationTime := time.Hour * time.Duration(jwtExpiration) // 72 hours
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(jwtExpirationTime).Unix(), // Token expires in 72 hours
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}