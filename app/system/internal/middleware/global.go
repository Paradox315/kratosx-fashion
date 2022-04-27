package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"kratosx-fashion/pkg/logutil"
)

type GlobalMiddleware struct {
	log log.Logger
}

func NewGlobalMiddleware(logger log.Logger) *GlobalMiddleware {
	return &GlobalMiddleware{
		log: logger,
	}
}

func (m GlobalMiddleware) Get() []fiber.Handler {
	return []fiber.Handler{
		recover.New(),
		cors.New(),
		compress.New(),
		fiberzap.New(fiberzap.Config{
			Logger: m.log.(*logutil.Logger).GetZap(),
		}),
	}
}
