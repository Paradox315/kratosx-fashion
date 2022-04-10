// Code generated by protoc-gen-go-xhttp. DO NOT EDIT.
// versions:
// protoc-gen-go-xhttp v1.0.0

package v1

import (
	context "context"
	middleware "github.com/go-kratos/kratos/v2/middleware"
	transport "github.com/go-kratos/kratos/v2/transport"
	xhttp "github.com/go-kratos/kratos/v2/transport/xhttp"
	apistate "github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	binding "github.com/go-kratos/kratos/v2/transport/xhttp/binding"
)

import fiber "github.com/gofiber/fiber/v2"

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.BindBody

const _ = xhttp.SupportPackageIsVersion1
const _ = middleware.SupportPackageIsVersion1
const _ = transport.KindXHTTP

var _ = new(apistate.Resp[any])

// 用户服务
type UserXHTTPServer interface {
	CreateUser(context.Context, *UserRequest) (*IDReply, error)
	DeleteUser(context.Context, *IDsRequest) (*EmptyReply, error)
	GetUser(context.Context, *IDRequest) (*UserReply, error)
	InitUserInfo(context.Context, *EmptyRequest) (*UserState, error)
	ListLoginLog(context.Context, *ListRequest) (*ListLoginLogReply, error)
	ListUser(context.Context, *ListSearchRequest) (*ListUserReply, error)
	UpdatePassword(context.Context, *PasswordRequest) (*IDReply, error)
	UpdateUser(context.Context, *UserRequest) (*IDReply, error)
	UpdateUserStatus(context.Context, *StatusRequest) (*IDReply, error)
}

func RegisterUserXHTTPServer(s *xhttp.Server, srv UserXHTTPServer) {
	s.Route(func(r fiber.Router) {
		api := r.Group("api/system/v1/user")
		// Register all service annotation
		{
			api.Name("User-XHTTPServer")
			api.Use(middleware.Authenticator(), middleware.Authorizer())
		}
		api.Post("/", _User_CreateUser0_XHTTP_Handler(srv)).Name("User-CreateUser.0-XHTTP_Handler")
		api.Put("/", _User_UpdateUser0_XHTTP_Handler(srv)).Name("User-UpdateUser.0-XHTTP_Handler")
		api.Put("/password", _User_UpdatePassword0_XHTTP_Handler(srv)).Name("User-UpdatePassword.0-XHTTP_Handler")
		api.Put("/status", _User_UpdateUserStatus0_XHTTP_Handler(srv)).Name("User-UpdateUserStatus.0-XHTTP_Handler")
		api.Delete("/:ids", _User_DeleteUser0_XHTTP_Handler(srv)).Name("User-DeleteUser.0-XHTTP_Handler")
		api.Get("/:id", _User_GetUser0_XHTTP_Handler(srv)).Name("User-GetUser.0-XHTTP_Handler")
		api.Get("/init/info", _User_InitUserInfo0_XHTTP_Handler(srv)).Name("User-InitUserInfo.0-XHTTP_Handler")
		api.Post("/list", _User_ListUser0_XHTTP_Handler(srv)).Name("User-ListUser.0-XHTTP_Handler")
		api.Get("/log/list/:page_num/:page_size", _User_ListLoginLog0_XHTTP_Handler(srv)).Name("User-ListLoginLog.0-XHTTP_Handler")
	})
}

func _User_CreateUser0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in UserRequest
		if err := binding.BindBody(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.CreateUser(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*IDReply]().WithData(reply).Send(c)
	}
}

func _User_UpdateUser0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in UserRequest
		if err := binding.BindBody(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.UpdateUser(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*IDReply]().WithData(reply).Send(c)
	}
}

func _User_UpdatePassword0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in PasswordRequest
		if err := binding.BindBody(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.UpdatePassword(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*IDReply]().WithData(reply).Send(c)
	}
}

func _User_UpdateUserStatus0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in StatusRequest
		if err := binding.BindBody(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.UpdateUserStatus(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*IDReply]().WithData(reply).Send(c)
	}
}

func _User_DeleteUser0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in IDsRequest
		if err := binding.BindParams(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.DeleteUser(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*EmptyReply]().WithData(reply).Send(c)
	}
}

func _User_GetUser0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in IDRequest
		if err := binding.BindParams(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.GetUser(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*UserReply]().WithData(reply).Send(c)
	}
}

func _User_InitUserInfo0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in EmptyRequest
		if err := binding.BindQuery(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.InitUserInfo(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*UserState]().WithData(reply).Send(c)
	}
}

func _User_ListUser0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in ListSearchRequest
		if err := binding.BindBody(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.ListUser(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*ListUserReply]().WithData(reply).Send(c)
	}
}

func _User_ListLoginLog0_XHTTP_Handler(srv UserXHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in ListRequest
		if err := binding.BindParams(c, &in); err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		ctx := transport.NewFiberContext(context.Background(), c)
		reply, err := srv.ListLoginLog(ctx, &in)
		if err != nil {
			return apistate.Error[any]().WithError(err).Send(c)
		}
		return apistate.Success[*ListLoginLogReply]().WithData(reply).Send(c)
	}
}
