package server

import (
	"github.com/gofiber/contrib/fiberzap"
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/app/system/internal/service"
	"kratosx-fashion/pkg/logutil"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	v1 "kratosx-fashion/api/system/v1"
)

// NewHTTPServer new a XHTTP server.
func NewHTTPServer(c *conf.Server,
	publicSrv *service.PubService,
	userSrv *service.UserService,
	roleSrv *service.RoleService,
	resourceSrv *service.ResourceService,
	logger log.Logger) *xhttp.Server {
	var opts = []xhttp.ServerOption{
		xhttp.Middleware(
			recover.New(),
			fiberzap.New(fiberzap.Config{
				Logger: logger.(*logutil.Logger).GetZap(),
			}),
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
	{
		v1.RegisterPubXHTTPServer(srv, publicSrv)
		v1.RegisterUserXHTTPServer(srv, userSrv)
		v1.RegisterRoleXHTTPServer(srv, roleSrv)
		v1.RegisterResourceXHTTPServer(srv, resourceSrv)
	}

	return srv
}
