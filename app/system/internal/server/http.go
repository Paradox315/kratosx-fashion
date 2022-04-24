package server

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp"
	"github.com/gofiber/fiber/v2"
	v1 "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/conf"
	mw "kratosx-fashion/app/system/internal/middleware"
	"kratosx-fashion/app/system/internal/service"
)

// NewHTTPServer new a XHTTP server.
func NewHTTPServer(c *conf.Server,
	jwtSrv *mw.JWTService,
	casbinSrv *mw.CasbinAuth,
	globalMw *mw.GlobalMiddleware,
	publicSrv *service.PubService,
	userSrv *service.UserService,
	roleSrv *service.RoleService,
	resourceSrv *service.ResourceService,
	logger log.Logger) *xhttp.Server {
	var opts = []xhttp.ServerOption{
		xhttp.Logger(log.With(logger, "server", "xhttp")),
		xhttp.Middleware(
			globalMw.Get()...,
		),
		xhttp.FiberConfig(fiber.Config{
			JSONDecoder: encoding.GetCodec("json").Unmarshal,
			JSONEncoder: encoding.GetCodec("json").Marshal,
		}),
	}
	if c.Http.Network != "" {
		opts = append(opts, xhttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, xhttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, xhttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := xhttp.NewServer(opts...)
	srv.Route(func(r fiber.Router) {
		r.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Welcome to KratosX-Fashion!")
		})
	})
	log.NewHelper(logger).Info("xhttp server middleware init", jwtSrv.Name(), casbinSrv.Name())
	{
		v1.RegisterPubXHTTPServer(srv, publicSrv)
		v1.RegisterUserXHTTPServer(srv, userSrv)
		v1.RegisterRoleXHTTPServer(srv, roleSrv)
		v1.RegisterResourceXHTTPServer(srv, resourceSrv)
	}
	return srv
}
