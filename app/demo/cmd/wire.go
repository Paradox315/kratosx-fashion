//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratosx-fashion/app/demo/internal/biz"
	"kratosx-fashion/app/demo/internal/conf"
	"kratosx-fashion/app/demo/internal/data"
	"kratosx-fashion/app/demo/internal/server"
	"kratosx-fashion/app/demo/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Registry, *conf.Data, *conf.Logger, *conf.Storage) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
