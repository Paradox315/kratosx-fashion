package ctxutil

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
)

var ErrorContextParse = errors.InternalServer("CONTEXT_PARSE_ERROR", "context parse error")

func SetUid(ctx context.Context, uid uint) error {
	if c, ok := transport.FromFiberContext(ctx); ok {
		c.Locals("uid", uid)
		return nil
	}
	return ErrorContextParse
}

func SetUsername(ctx context.Context, username string) error {
	if c, ok := transport.FromFiberContext(ctx); ok {
		c.Locals("username", username)
		return nil
	}
	return ErrorContextParse
}
func SetMobile(ctx context.Context, mobile string) error {
	if c, ok := transport.FromFiberContext(ctx); ok {
		c.Locals("mobile", mobile)
		return nil
	}
	return ErrorContextParse
}
func SetEmail(ctx context.Context, email string) error {
	if c, ok := transport.FromFiberContext(ctx); ok {
		c.Locals("email", email)
		return nil
	}
	return ErrorContextParse
}
func SetRoleIDs(ctx context.Context, rids []uint) error {
	if c, ok := transport.FromFiberContext(ctx); ok {
		c.Locals("roles", rids)
		return nil
	}
	return ErrorContextParse
}
func GetUid(ctx context.Context) (uid uint) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("uid").(uint)
	}
	return
}

func GetUsername(ctx context.Context) (username string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("username").(string)
	}
	return
}

func GetMobile(ctx context.Context) (mobile string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("mobile").(string)
	}
	return
}

func GetEmail(ctx context.Context) (mobile string) {
	if c, ok := transport.FromFiberContext(ctx); ok {
		return c.Locals("email").(string)
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
