package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/novychok/authasvs/internal/database/dao"
	"github.com/novychok/authasvs/internal/database/pqmodels"
	"github.com/novychok/authasvs/internal/entity"
	"github.com/novychok/authasvs/internal/pkg/postgres"
	"github.com/novychok/authasvs/internal/repository"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) UpdatePassword(ctx context.Context,
	update *entity.PasswordUpdate) error {
	user, err := pqmodels.Users(
		pqmodels.UserWhere.Email.EQ(update.Email),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	user.PasswordHash = update.NewPasswordHash
	user.UpdatedAt = time.Now()

	_, err = user.Update(ctx, r.db, boil.Whitelist("password_hash", "updated_at"))
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) Create(ctx context.Context,
	user *entity.User) error {
	userDB := &pqmodels.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	err := userDB.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}

	user.ID = entity.UserID(userDB.ID)
	user.CreatedAt = userDB.CreatedAt
	user.UpdatedAt = userDB.UpdatedAt

	return nil
}

func (r *postgresRepository) Get(ctx context.Context,
	id entity.UserID) (*entity.User, error) {
	userDB, err := pqmodels.Users(
		pqmodels.UserWhere.ID.EQ(id.String()),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &entity.User{}
	dao.UserTo(userDB, user)

	return user, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context,
	email string) (*entity.User, error) {
	userDB, err := pqmodels.Users(
		pqmodels.UserWhere.Email.EQ(email),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	user := &entity.User{}
	dao.UserTo(userDB, user)

	return user, nil
}

func NewPostgres(db postgres.Connection) repository.Auth {
	return &postgresRepository{
		db: db,
	}
}
