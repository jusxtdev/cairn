package auth

import (
	"net/http"
	"time"

	"github.com/jusxtdev/cairn/internal/config"
)

var TokenName = "token"

func SetJWTCookie(w http.ResponseWriter, token string, cfg *config.Config){
	cookie := &http.Cookie{
		Name: TokenName,
		Value: token,
		Expires: time.Now().Add(time.Hour * time.Duration(cfg.JWTExpiration)),
		HttpOnly: true,
		Secure : true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func GetJWTCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(TokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}
	return cookie.Value, nil
}