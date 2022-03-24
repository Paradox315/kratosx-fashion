// Code generated by protoc-gen-go-xhttp. DO NOT EDIT.
// versions:
// protoc-gen-go-xhttp v1.0.0

package v1

import (
	context "context"
	middleware "github.com/go-kratos/kratos/v2/middleware"
	xhttp "github.com/go-kratos/kratos/v2/transport/xhttp"
	apistate "github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	binding "github.com/go-kratos/kratos/v2/transport/xhttp/binding"
	"github.com/gofiber/fiber/v2"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.BindBody

const _ = xhttp.SupportPackageIsVersion1
const _ = middleware.SupportPackageIsVersion1

var _ = new(apistate.Resp)

// 公共接口
type PubXHTTPServer interface {
	Generate(context.Context, *EmptyRequest) (*CaptchaReply, error)
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Logout(context.Context, *EmptyRequest) (*EmptyReply, error)
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	RetrievePwd(context.Context, *RetrieveRequest) (*EmptyReply, error)
}

func RegisterPubXHTTPServer(s *xhttp.Server, srv PubXHTTPServer) {
	s.Route(func(r fiber.Router) {
		api := r.Group("api/system/v1/pub")
		// Register all service annotation
		{
		}
		api.Get("/captcha", _Pub_Generate0_XHTTP_Handler(srv))
		api.Post("/register", _Pub_Register0_XHTTP_Handler(srv))
		api.Post("/login", _Pub_Login0_XHTTP_Handler(srv))
		api.Post("/logout", _Pub_Logout0_XHTTP_Handler(srv))
		api.Post("/retrieve", _Pub_RetrievePwd0_XHTTP_Handler(srv))
	})
}

//
func _Pub_Generate0_XHTTP_Handler(srv PubXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in EmptyRequest
		if err := binding.BindQuery(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.Generate(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Pub_Register0_XHTTP_Handler(srv PubXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in RegisterRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.Register(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Pub_Login0_XHTTP_Handler(srv PubXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in LoginRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.Login(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Pub_Logout0_XHTTP_Handler(srv PubXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in EmptyRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.Logout(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Pub_RetrievePwd0_XHTTP_Handler(srv PubXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in RetrieveRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.RetrievePwd(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}
