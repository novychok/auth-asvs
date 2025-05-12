package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/novychok/authasvs/internal/entity"
	"github.com/novychok/authasvs/internal/pkg/jwts"
	"github.com/novychok/authasvs/internal/pkg/validator"
	"github.com/novychok/authasvs/internal/repository"
	"github.com/novychok/authasvs/internal/service"

	"golang.org/x/crypto/bcrypt"
)

const (
	tokenDuration        = 15 * time.Minute
	refreshTokenDuration = 24 * time.Hour
	issuer               = "authasvs"
)

type srv struct {
	l                *slog.Logger
	v                *validator.ValidatorASVS
	jwtSecretManager *jwts.SecretManager

	authRepo repository.Auth
}

const (
	TokenPayloadSize = 60
)

func (s *srv) ChangePassword(ctx context.Context,
	changePasswordRequest *entity.ChangePasswordRequest) error {
	l := s.l.With(slog.String("method", "ChangePassword"))

	login := &entity.UserLogin{
		Email:    changePasswordRequest.Email,
		Password: changePasswordRequest.CurrentPassword,
	}

	token, err := s.Login(ctx, login)
	if err != nil {
		l.Error("user login failed", slog.Any("error", err))

		return err
	}

	verifyToken := &entity.VerifyToken{
		Token: token.Token,
	}

	if err := s.VerifyToken(ctx, verifyToken); err != nil {
		l.Error("token verify failed", slog.Any("error", err))

		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword),
		bcrypt.DefaultCost)
	if err != nil {
		l.Error("generate password hash failed", slog.Any("error", err))

		return err
	}

	update := &entity.PasswordUpdate{
		Email:           login.Email,
		NewPasswordHash: string(hash),
	}

	if err := s.authRepo.UpdatePassword(ctx, update); err != nil {
		l.Error("failed to update password", slog.Any("error", err))

		return err
	}

	return nil
}

func (s *srv) Register(ctx context.Context, userCreate *entity.UserCreate) (*entity.UserToken, error) {
	l := s.l.With(slog.String("method", "Register"))

	// err := s.v.StructCtx(ctx, userCreate)
	// if err != nil {
	// 	l.Error("validation failed", slog.Any("error", err))

	// 	return nil, err
	// }

	hash, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Error("generate password hash failed", slog.Any("error", err))

		return nil, err
	}

	user := &entity.User{
		Name:         userCreate.Name,
		Email:        userCreate.Email,
		PasswordHash: string(hash),
	}

	err = s.authRepo.Create(ctx, user)
	if err != nil {
		l.Error("create user failed", slog.Any("error", err))

		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		l.Error("generate token failed", slog.Any("error", err))

		return nil, err
	}

	return token, nil
}

func (s *srv) Login(ctx context.Context, login *entity.UserLogin) (*entity.UserToken, error) {
	l := s.l.With(slog.String("method", "LogIn"))

	// err := s.v.StructCtx(ctx, login)
	// if err != nil {
	// 	l.Error("validation failed", slog.Any("error", err))

	// 	return nil, err
	// }

	user, err := s.authRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		l.Error("get user by email failed", slog.Any("error", err))

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		l.Error("passwords do not match", slog.Any("error", err))

		return nil, entity.ErrAuthorizationFailed
	}

	token, err := s.generateToken(user)
	if err != nil {
		l.Error("generate token failed", slog.Any("error", err))

		return nil, err
	}

	return token, nil
}

func (s *srv) RefreshToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.UserToken, error) {
	l := s.l.With(slog.String("method", "RefreshToken"))

	// err := s.v.StructCtx(ctx, req)
	// if err != nil {
	// 	l.Error("validation failed", slog.Any("error", err))

	// 	return nil, err
	// }

	claims := &entity.UserClaims{}

	_, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretManager.PublicKey(), nil
	})
	if err != nil {
		l.Error("parse refresh token failed", slog.Any("error", err))

		return nil, err
	}

	user, err := s.authRepo.Get(ctx, entity.UserID(claims.Subject))
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		l.Error("generate token failed", slog.Any("error", err))

		return nil, err
	}

	return token, nil
}

func (s *srv) VerifyToken(ctx context.Context, verifyRequest *entity.VerifyToken) error {
	l := s.l.With(slog.String("method", "VerifyToken"))

	// err := s.v.StructCtx(ctx, verifyRequest)
	// if err != nil {
	// 	l.Error("validation failed", slog.Any("error", err))

	// 	return err
	// }

	_, err := s.getClaims(verifyRequest.Token)
	if err != nil {
		l.Error("get claims failed", slog.Any("error", err))

		return err
	}

	return nil
}

func (s *srv) GetByToken(ctx context.Context, token string) (*entity.User, error) {
	l := s.l.With(slog.String("method", "GetUserByToken"))

	claims, err := s.getClaims(token)
	if err != nil {
		l.Error("get claims failed", slog.Any("error", err))

		return nil, err
	}

	user, err := s.authRepo.Get(ctx, entity.UserID(claims.Subject))
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	return user, nil
}

func (s *srv) Get(ctx context.Context, id entity.UserID) (*entity.User, error) {
	l := s.l.With(slog.String("method", "Get"))

	user, err := s.authRepo.Get(ctx, id)
	if err != nil {
		l.Error("get user by id failed", slog.Any("error", err))

		return nil, err
	}

	return user, nil
}

func (s *srv) generateToken(user *entity.User) (*entity.UserToken, error) {
	now := time.Now()
	tokenID := uuid.New().String()
	refreshTokenID := uuid.New().String()
	tokenExpiresAt := now.Add(tokenDuration)
	refreshTokenExpiresAt := now.Add(refreshTokenDuration)

	claims := entity.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(tokenExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	refreshClaims := entity.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshTokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	tokenString, err := token.SignedString(s.jwtSecretManager.PrivateKey())
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtSecretManager.PrivateKey())
	if err != nil {
		return nil, err
	}

	return &entity.UserToken{
		Token:            tokenString,
		ExpiresAt:        tokenExpiresAt,
		RefreshToken:     refreshTokenString,
		RefreshExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func (s *srv) getClaims(token string) (*entity.UserClaims, error) {
	claims := &entity.UserClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretManager.PublicKey(), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func New(
	l *slog.Logger,
	v *validator.ValidatorASVS,
	jwtSecretManager *jwts.SecretManager,

	authRepo repository.Auth,
) service.Auth {
	return &srv{
		l:                l.With(slog.String("service", "auth")),
		v:                v,
		jwtSecretManager: jwtSecretManager,

		authRepo: authRepo,
	}
}
