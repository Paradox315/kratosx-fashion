package infra

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	consulAPI "github.com/hashicorp/consul/api"
	"kratosx-fashion/app/system/conf"
)

func NewDiscovery(conf *conf.Consul) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Address
	c.Scheme = conf.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func NewRegistrar(conf *conf.Consul) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Address
	c.Scheme = conf.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
