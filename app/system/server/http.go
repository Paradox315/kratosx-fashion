package server

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"kratosx-fashion/app/system/conf"

	fashion_v1 "kratosx-fashion/api/fashion/v1"
	sys_v1 "kratosx-fashion/api/system/v1"
	fashion_srv "kratosx-fashion/app/fashion/service"
	mw "kratosx-fashion/app/system/middleware"
	sys_srv "kratosx-fashion/app/system/service"
)

// NewHTTPServer new a XHTTP server.
func NewHTTPServer(c *conf.Server,
	jwtMw *mw.JWTService,
	casbinMw *mw.CasbinAuth,
	cache *mw.Cache,
	limiter *mw.Limiter,
	logHook *mw.Logger,
	globalMw *mw.GlobalMiddleware,
	publicSrv *sys_srv.PubService,
	userSrv *sys_srv.UserService,
	roleSrv *sys_srv.RoleService,
	resourceSrv *sys_srv.ResourceService,
	clothesSrv *fashion_srv.ClothesService,
	recommendSrv *fashion_srv.RecommendService,
	tryOnSrv *fashion_srv.TryOnService,
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
		r.Get("/csrf", func(c *fiber.Ctx) error {
			return c.SendString("Welcome to KratosX-Fashion!")
		})
		r.Get("/monitor", monitor.New())
	})
	log.NewHelper(logger).Infof("xhttp server middleware init: %s,%s,%s,%s,%s",
		jwtMw.Name(),
		casbinMw.Name(),
		cache.Name(),
		limiter.Name(),
		logHook.Name(),
	)
	{
		sys_v1.RegisterPubXHTTPServer(srv, publicSrv)
		sys_v1.RegisterUserXHTTPServer(srv, userSrv)
		sys_v1.RegisterRoleXHTTPServer(srv, roleSrv)
		sys_v1.RegisterResourceXHTTPServer(srv, resourceSrv)
	}
	{
		fashion_v1.RegisterClothesXHTTPServer(srv, clothesSrv)
		fashion_v1.RegisterRecommendXHTTPServer(srv, recommendSrv)
		fashion_v1.RegisterTryOnXHTTPServer(srv, tryOnSrv)
	}
	return srv
}
