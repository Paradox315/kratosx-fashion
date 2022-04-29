package middleware

import (
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/utils"
	"sync"
	"time"
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
	return cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   1 * time.Minute,
		CacheHeader:  "X-Cache",
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
		ExpirationGenerator:  nil,
		StoreResponseHeaders: false,
		Storage:              nil,
	})
}

func (c *Cache) Name() string {
	return kmw.CacheCfg
}
