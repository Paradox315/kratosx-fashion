//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/middleware"
	"kratosx-fashion/app/system/internal/server"
	"kratosx-fashion/app/system/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Registry, *conf.Data, *conf.JWT, *conf.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, middleware.ProviderSet, newApp))
}
