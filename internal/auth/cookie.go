package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/jusxtdev/cairn/internal/config"
)

var TokenName = "token"

func SetJWTCookie(w http.ResponseWriter, token string, cfg *config.Config){
	secure := false
	if cfg.Env == "production" {
		secure = true
	}
	
	cookie := &http.Cookie{
		Name: TokenName,
		Value: token,
		Expires: time.Now().Add(time.Hour * time.Duration(cfg.JWTExpiration)),
		HttpOnly: true,
		Secure : secure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func GetJWTCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(TokenName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}
	return cookie.Value, nil
}