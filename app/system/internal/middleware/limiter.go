package middleware

import (
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"sync"
	"time"
)

var _ kmw.FiberMiddleware = (*Limiter)(nil)

type Limiter struct {
	once sync.Once
}

func NewLimiter() *Limiter {
	l := &Limiter{}
	l.once.Do(func() {
		kmw.RegisterMiddleware(l)
	})
	return l
}

func (l *Limiter) MiddlewareFunc() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1" || c.IP() == "::1"
		},
		Max:        10000,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
	})
}

func (l *Limiter) Name() string {
	return kmw.LimiterCfg
}
