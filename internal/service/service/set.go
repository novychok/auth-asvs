package service

import (
	"github.com/google/wire"
	"github.com/novychok/authasvs/internal/service/auth"
)

var Set = wire.NewSet(
	auth.New,
)
