package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"kratosx-fashion/pkg/logutil"
	"strings"
)

type GlobalMiddleware struct {
	log log.Logger
}

func NewGlobalMiddleware(logger log.Logger) *GlobalMiddleware {
	return &GlobalMiddleware{
		log: logger,
	}
}

func (m *GlobalMiddleware) Install() []fiber.Handler {
	return []fiber.Handler{
		recover.New(),
		//csrf.New(),
		cors.New(cors.Config{
			Next:         nil,
			AllowOrigins: "*",
			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodHead,
				fiber.MethodPut,
				fiber.MethodDelete,
				fiber.MethodPatch,
			}, ","),
			AllowHeaders:     "",
			AllowCredentials: true,
			ExposeHeaders:    "",
			MaxAge:           0,
		}),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		requestid.New(),
		fiberzap.New(fiberzap.Config{
			Logger: m.log.(*logutil.Logger).GetZap(),
			Fields: []string{"latency", "status", "method", "url", "requestId"},
		}),
	}
}
