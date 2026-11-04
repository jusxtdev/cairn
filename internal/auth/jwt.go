package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserId string `json:"userId"` 
	jwt.RegisteredClaims
}

func GenerateJWT(userId string, jwtSecret []byte, jwtExpiration int) (string, error) {
	jwtExpirationTime := time.Hour * time.Duration(jwtExpiration) // 72 hours
	claims := jwt.RegisteredClaims{
		Subject: userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpirationTime)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}