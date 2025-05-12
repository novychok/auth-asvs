package authapiv1

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/novychok/authasvs/internal/entity"
	authapiv1 "github.com/novychok/authasvs/pkg/authApi/v1"
)

func (s *Server) auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customCtx := ContextFromRequest(r)
		_, exists := customCtx.Get("authRequired")
		if !exists {
			h.ServeHTTP(w, r.WithContext(r.Context()))

			return
		}

		var tokenValue string

		tokenCookie, err := r.Cookie(tokenKey)
		if err != nil {
			refreshTokenCookie, err := r.Cookie(refreshTokenKey)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.DefaultResponder(w, r, authapiv1.Error{
					Message: "unauthorized",
				})

				return
			}

			token, err := s.authService.RefreshToken(r.Context(), &entity.RefreshTokenRequest{
				RefreshToken: refreshTokenCookie.Value,
			})
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.DefaultResponder(w, r, authapiv1.Error{
					Message: "unauthorized",
				})

				return
			}

			setCookies(token, w)

			tokenValue = token.Token

		} else {
			tokenValue = tokenCookie.Value
		}

		user, err := s.authService.GetByToken(r.Context(), tokenValue)
		if err != nil {
			h.ServeHTTP(w, r.WithContext(r.Context()))

			return
		}

		ctx := WithUser(r.Context(), user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
