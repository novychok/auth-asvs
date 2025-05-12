//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/novychok/authasvs/internal/config"
	"github.com/novychok/authasvs/internal/handler/authapiv1"
	"github.com/novychok/authasvs/internal/pkg"
	"github.com/novychok/authasvs/internal/repository/repository"
	"github.com/novychok/authasvs/internal/service/service"
)

func Init() (*App, func(), error) {
	wire.Build(
		config.Set,
		pkg.Set,

		repository.Set,
		service.Set,
		authapiv1.Set,

		New,
	)

	return nil, nil, nil
}
