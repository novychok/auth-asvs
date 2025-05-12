package repository

import (
	"context"

	"github.com/novychok/authasvs/internal/entity"
)

type Auth interface {
	UpdatePassword(ctx context.Context, update *entity.PasswordUpdate) error
	Create(ctx context.Context, user *entity.User) error
	Get(ctx context.Context, id entity.UserID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
