package middleware

import (
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"sync"
)

var _ kmw.FiberMiddleware = (*Cache)(nil)

type Cache struct {
	once sync.Once
}

func NewCache() *Cache {
	c := &Cache{}
	c.once.Do(func() {
		kmw.RegisterMiddleware(c)
	})
	return c
}

func (c *Cache) MiddlewareFunc() fiber.Handler {
	cacheHandler := cache.New()
	return func(ctx *fiber.Ctx) error {
		if ctx.Query("refresh") == "true" {
			return ctx.Next()
		}
		return cacheHandler(ctx)
	}
}

func (c *Cache) Name() string {
	return kmw.CacheCfg
}
