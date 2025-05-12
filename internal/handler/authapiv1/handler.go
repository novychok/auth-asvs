package authapiv1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/novychok/authasvs/internal/entity"
	"github.com/novychok/authasvs/internal/service"

	authapiv1 "github.com/novychok/authasvs/pkg/authApi/v1"
)

const (
	tokenKey        = "token"
	refreshTokenKey = "refreshToken"
)

//go:generate oapi-codegen --config=./oapi-codegen.yaml ../../../api/authApi/openapi/v1.yaml
type handler struct {
	authService service.Auth
}

func (h *handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req authapiv1.ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errBadRequest(w, r, err.Error())

		return
	}

	chanePasswordReq := &entity.ChangePasswordRequest{
		Email:              string(req.Email),
		CurrentPassword:    req.CurrentPassword,
		NewPassword:        req.NewPassword,
		NewPasswordConfirm: req.NewPasswordConfirm,
	}

	h.authService.ChangePassword(r.Context(), chanePasswordReq)

	response(w, http.StatusOK, nil)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authapiv1.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errBadRequest(w, r, err.Error())

		return
	}

	loginRequest := &entity.UserLogin{
		Email:    string(req.Email),
		Password: req.Password,
	}

	userToken, err := h.authService.Login(r.Context(), loginRequest)
	if err != nil {
		if errors.Is(err, entity.ErrAuthorizationFailed) || errors.Is(err, entity.ErrUserNotFound) {
			errUnauthorized(w, r, "invalid credentials")

			return
		}

		errInternal(w, r, "internal server error")

		return
	}

	response(w, http.StatusOK, authapiv1.LoginResponse{
		AccessToken:           userToken.Token,
		AccessTokenExpiresAt:  userToken.ExpiresAt,
		RefreshToken:          userToken.RefreshToken,
		RefreshTokenExpiresAt: userToken.RefreshExpiresAt,
	})
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req authapiv1.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errBadRequest(w, r, err.Error())

		return
	}

	registerRequest := &entity.UserCreate{
		Email:           string(req.Email),
		Password:        req.Password,
		Name:            req.Name,
		PasswordConfirm: req.PasswordConfirm,
	}

	userToken, err := h.authService.Register(r.Context(), registerRequest)
	if err != nil {
		errInternal(w, r, err.Error())

		return
	}

	response(w, http.StatusOK, authapiv1.LoginResponse{
		AccessToken:           userToken.Token,
		AccessTokenExpiresAt:  userToken.ExpiresAt,
		RefreshToken:          userToken.RefreshToken,
		RefreshTokenExpiresAt: userToken.RefreshExpiresAt,
	})
}

func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	var req authapiv1.RefreshTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errBadRequest(w, r, err.Error())

		return
	}

	token, err := h.authService.RefreshToken(r.Context(), &entity.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		errInternal(w, r, err.Error())

		return
	}

	response(w, http.StatusOK, authapiv1.LoginResponse{
		AccessToken:           token.Token,
		AccessTokenExpiresAt:  token.ExpiresAt,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresAt: token.RefreshExpiresAt,
	})
}

func NewHandler(
	authService service.Auth,
) authapiv1.ServerInterface {
	h := &handler{
		authService: authService,
	}

	return h
}
