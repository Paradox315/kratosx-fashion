package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/conf"
	service2 "kratosx-fashion/app/system/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server,
	publicSrv *service2.PubService,
	userSrv *service2.UserService,
	roleSrv *service2.RoleService,
	menuSrv *service2.MenuService,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	{
		v1.RegisterPubServer(srv, publicSrv)
		v1.RegisterUserServer(srv, userSrv)
		v1.RegisterRoleServer(srv, roleSrv)
		v1.RegisterMenuServer(srv, menuSrv)
	}
	return srv
}
