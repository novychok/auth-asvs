package repository

import (
	"github.com/google/wire"
	"github.com/novychok/authasvs/internal/repository/auth"
)

var Set = wire.NewSet(
	auth.NewPostgres,
)
