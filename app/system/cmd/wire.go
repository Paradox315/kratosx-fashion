//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"kratosx-fashion/app/system/conf"
	"kratosx-fashion/app/system/data"
	"kratosx-fashion/app/system/data/infra"
	"kratosx-fashion/app/system/middleware"
	"kratosx-fashion/app/system/server"

	sys_biz "kratosx-fashion/app/system/biz"
	sys_repo "kratosx-fashion/app/system/data/repo"
	sys_srv "kratosx-fashion/app/system/service"

	fashion_biz "kratosx-fashion/app/fashion/biz"
	fashion_repo "kratosx-fashion/app/fashion/repo"
	fashion_srv "kratosx-fashion/app/fashion/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Consul, *conf.Algorithm, *conf.Storage, *conf.Data, *conf.JWT, *conf.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			server.ProviderSet,
			data.ProviderSet,
			infra.ProviderSet,
			middleware.ProviderSet,

			sys_repo.ProviderSet,
			sys_biz.ProviderSet,
			sys_srv.ProviderSet,

			fashion_repo.ProviderSet,
			fashion_biz.ProviderSet,
			fashion_srv.ProviderSet,
			newApp,
		),
	)
}
