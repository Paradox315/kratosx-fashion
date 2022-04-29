package server

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/gofiber/fiber/v2"
	v1 "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/conf"
	mw "kratosx-fashion/app/system/internal/middleware"
	"kratosx-fashion/app/system/internal/service"
)

// NewHTTPServer new a XHTTP server.
func NewHTTPServer(c *conf.Server,
	jwtMw *mw.JWTService,
	casbinMw *mw.CasbinAuth,
	cache *mw.Cache,
	limiter *mw.Limiter,
	globalMw *mw.GlobalMiddleware,
	publicSrv *service.PubService,
	userSrv *service.UserService,
	roleSrv *service.RoleService,
	resourceSrv *service.ResourceService,
	logger log.Logger) *xhttp.Server {
	var opts = []xhttp.ServerOption{
		xhttp.Logger(log.With(logger, "server", "xhttp")),
		xhttp.Middleware(
			globalMw.Install()...,
		),
		xhttp.FiberConfig(fiber.Config{
			JSONDecoder: encoding.GetCodec("json").Unmarshal,
			JSONEncoder: encoding.GetCodec("json").Marshal,
			// Override default error handler
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				// Default 500 statuscode
				code := fiber.StatusInternalServerError

				if e, ok := err.(*fiber.Error); ok {
					// Override status code if fiber.Error type
					code = e.Code
				}

				// Return statuscode with error message
				return apistate.Error[any]().WithError(err).WithCode(code).Send(c)
			},
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
	log.NewHelper(logger).Info("xhttp server middleware init", jwtMw.Name(), casbinMw.Name(), cache.Name(), limiter.Name())
	{
		v1.RegisterPubXHTTPServer(srv, publicSrv)
		v1.RegisterUserXHTTPServer(srv, userSrv)
		v1.RegisterRoleXHTTPServer(srv, roleSrv)
		v1.RegisterResourceXHTTPServer(srv, resourceSrv)
	}
	return srv
}
