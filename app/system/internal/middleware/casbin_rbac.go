package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/gofiber/fiber/v2"
	"os"
	"sync"
)

var _ kmw.FiberMiddleware = (*CasbinAuth)(nil)

const casbinReason = "CASBIN_RBAC_ERROR"

type CasbinAuth struct {
	once sync.Once
	e    *casbin.SyncedEnforcer
	log  *log.Helper
}

func (c *CasbinAuth) MiddlewareFunc() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if os.Getenv("env") == "dev" {
			return ctx.Next()
		}
		// 获取请求的URI
		obj := ctx.Path()
		// 获取请求方法
		act := ctx.Method()
		// 获取用户的角色
		rids := ctx.Locals("roles").([]string)
		checked := false
		for _, rid := range rids {
			success, e := c.e.Enforce(rid, obj, act)
			if success && e == nil {
				checked = true
				break
			}
		}
		if !checked {
			return apistate.Error[any]().WithError(errors.Unauthorized(casbinReason, "权限不足")).Send(ctx)
		}
		return ctx.Next()
	}
}

func (c *CasbinAuth) Name() string {
	return kmw.AuthorizerCfg
}

func NewCasbinAuth(e *casbin.SyncedEnforcer, logger log.Logger) *CasbinAuth {
	c := &CasbinAuth{
		e:   e,
		log: log.NewHelper(log.With(logger, "middleware", "authorizer")),
	}
	c.once.Do(func() {
		kmw.RegisterMiddleware(c)
	})
	return c
}
