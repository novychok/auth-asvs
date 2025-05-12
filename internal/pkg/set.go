package pkg

import (
	"github.com/google/wire"
	"github.com/novychok/authasvs/internal/pkg/context"
	"github.com/novychok/authasvs/internal/pkg/jwts"
	"github.com/novychok/authasvs/internal/pkg/postgres"
	"github.com/novychok/authasvs/internal/pkg/slog"
	"github.com/novychok/authasvs/internal/pkg/validator"
	"github.com/novychok/authasvs/internal/pkg/vault"
)

var Set = wire.NewSet(
	context.New,
	jwts.New,
	postgres.New,
	slog.New,
	validator.New,
	vault.New,
)
