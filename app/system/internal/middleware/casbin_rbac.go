package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/gofiber/fiber/v2"
	"os"
	"sync"
)

var _ kmw.FiberMiddleware = (*CasbinAuth)(nil)

type CasbinAuth struct {
	once sync.Once
	log  *log.Helper
}

func (c *CasbinAuth) MiddlewareFunc() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if os.Getenv("env") == "dev" {
			return c.Next()
		}
		return c.Next()
	}
}

func (c *CasbinAuth) Name() string {
	return kmw.AuthorizerCfg
}

func NewCasbinAuth(logger log.Logger) *CasbinAuth {
	c := &CasbinAuth{
		log: log.NewHelper(log.With(logger, "middleware", "authorizer")),
	}
	c.once.Do(func() {
		kmw.RegisterMiddleware(c)
	})
	return c
}
