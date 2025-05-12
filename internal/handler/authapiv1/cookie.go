package authapiv1

import (
	"net/http"

	"github.com/novychok/authasvs/internal/entity"
)

func setCookies(token *entity.UserToken, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenKey,
		Value:    token.Token,
		Expires:  token.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    token.RefreshToken,
		Expires:  token.RefreshExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
