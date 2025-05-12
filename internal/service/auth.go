package service

import (
	"context"

	"github.com/novychok/authasvs/internal/entity"
)

type Auth interface {
	ChangePassword(ctx context.Context, changePasswordRequest *entity.ChangePasswordRequest) error
	Register(ctx context.Context, user *entity.UserCreate) (*entity.UserToken, error)
	Login(ctx context.Context, login *entity.UserLogin) (*entity.UserToken, error)
	RefreshToken(ctx context.Context, refreshRequest *entity.RefreshTokenRequest) (*entity.UserToken, error)
	VerifyToken(ctx context.Context, verifyRequest *entity.VerifyToken) error
	GetByToken(ctx context.Context, token string) (*entity.User, error)
	Get(ctx context.Context, id entity.UserID) (*entity.User, error)
}
