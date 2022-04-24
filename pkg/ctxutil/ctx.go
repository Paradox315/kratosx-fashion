package ctxutil

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/gofiber/fiber/v2"
)

func GetFiberCtx(ctx context.Context) (*fiber.Ctx, bool) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c, true
	}
	return nil, false
}

func GetUid(ctx context.Context) (uid string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("uid").(string)
	}
	return
}

func GetUsername(ctx context.Context) (username string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("username").(string)
	}
	return
}
func GetNickname(ctx context.Context) (nickname string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("nickname").(string)
	}
	return
}

func GetRoleIDs(ctx context.Context) (roleIDs []string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("roles").([]string)
	}
	return
}
