package authapiv1

import "github.com/google/wire"

var Set = wire.NewSet(
	NewServer,
	NewHandler,
)
