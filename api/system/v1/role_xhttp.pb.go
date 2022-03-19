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
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.BindBody

const _ = xhttp.SupportPackageIsVersion1
const _ = middleware.SupportPackageIsVersion1

var _ = new(apistate.Resp)

// 角色服务
type RoleXHTTPServer interface {
	CreateRole(context.Context, *RoleRequest) (*IDReply, error)
	DeleteRole(context.Context, *IDRequest) (*EmptyReply, error)
	GetRole(context.Context, *IDRequest) (*RoleReply, error)
	ListRole(context.Context, *ListRequest) (*ListRoleReply, error)
	UpdateRole(context.Context, *RoleRequest) (*IDReply, error)
	UpdateRoleStatus(context.Context, *IDRequest) (*IDReply, error)
}

func RegisterRoleXHTTPServer(s *xhttp.Server, srv RoleXHTTPServer) {
	s.Route(func(r fiber.Router) {
		api := r.Group("api/api/system/v1/role")
		// Register all service annotation
		{
			api.Use(middleware.Authenticator(), middleware.Authorizer())
		}
		api.Post("/", _Role_CreateRole0_XHTTP_Handler(srv))
		api.Put("/", _Role_UpdateRole0_XHTTP_Handler(srv))
		api.Put("/status", _Role_UpdateRoleStatus0_XHTTP_Handler(srv))
		api.Delete("/:id", _Role_DeleteRole0_XHTTP_Handler(srv))
		api.Get("/:id", _Role_GetRole0_XHTTP_Handler(srv))
		api.Post("/list", _Role_ListRole0_XHTTP_Handler(srv))
	})
}

//
func _Role_CreateRole0_XHTTP_Handler(srv RoleXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in RoleRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.CreateRole(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Role_UpdateRole0_XHTTP_Handler(srv RoleXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in RoleRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.UpdateRole(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Role_UpdateRoleStatus0_XHTTP_Handler(srv RoleXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in IDRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.UpdateRoleStatus(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Role_DeleteRole0_XHTTP_Handler(srv RoleXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in IDRequest
		if err := binding.BindParams(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.DeleteRole(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Role_GetRole0_XHTTP_Handler(srv RoleXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in IDRequest
		if err := binding.BindParams(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.GetRole(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}

//
func _Role_ListRole0_XHTTP_Handler(srv RoleXHTTPServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var in ListRequest
		if err := binding.BindBody(ctx, &in); err != nil {
			return err
		}
		reply, err := srv.ListRole(ctx.Context(), &in)
		if err != nil {
			return apistate.Error().WithError(err).Send(ctx)
		}
		return apistate.Success().WithData(reply).Send(ctx)
	}
}
