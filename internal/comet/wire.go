//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package comet

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"gnettest/internal/comet/conf"
	"gnettest/internal/comet/server"
)

// wireApp init kratos application.
func wireApp(*conf.Config) (*kratos.App, error) {
	panic(wire.Build(server.ProviderSet, newApp))
}
