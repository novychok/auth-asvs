package config

import "github.com/google/wire"

var Set = wire.NewSet(
	New,
	GetPostgres,
	GetPlatfromAPIV1,
	GetJWT,
)
