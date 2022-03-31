package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/pkg/ctxutil"
	"kratosx-fashion/pkg/xcast"
	"strings"
)

type ResourceService struct {
	pb.UnimplementedResourceServer

	uc  *biz.ResourceUsecase
	log *log.Helper
}

func NewResourceService(uc *biz.ResourceUsecase, logger log.Logger) *ResourceService {
	return &ResourceService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "service", "resource")),
	}
}

func (s *ResourceService) CreateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.IDReply, error) {
	id, err := s.uc.SaveMenu(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *ResourceService) UpdateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.IDReply, error) {
	id, err := s.uc.EditMenu(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *ResourceService) DeleteMenu(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	ids := strings.Split(req.Ids, ",")
	if err := s.uc.RemoveMenu(ctx, xcast.ToUintSlice(ids)); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *ResourceService) GetMenuTree(ctx context.Context, req *pb.IDRequest) (*pb.MenuReply, error) {
	roleIDs := ctxutil.GetRoleIDs(ctx)
	tree, err := s.uc.RoleMenuTree(ctx, xcast.ToUintSlice(roleIDs)...)
	if err != nil {
		return nil, err
	}
	resp := &pb.MenuReply{}
	for _, menu := range tree {
		var menuResp *pb.Menu
		_ = copier.Copy(menuResp, &menu)
		resp.Tree = append(resp.Tree, menuResp)
	}
	return resp, nil
}
func (s *ResourceService) GetMenuTreeByRole(ctx context.Context, req *pb.IDRequest) (*pb.MenuReply, error) {
	tree, err := s.uc.RoleMenuTree(ctx, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	resp := &pb.MenuReply{}
	for _, menu := range tree {
		var menuResp *pb.Menu
		_ = copier.Copy(menuResp, &menu)
		resp.Tree = append(resp.Tree, menuResp)
	}
	return &pb.MenuReply{}, nil
}
func (s *ResourceService) GetRouteTree(ctx context.Context, req *pb.IDRequest) (*pb.RouterReply, error) {
	roleIDs := ctxutil.GetRoleIDs(ctx)
	tree, err := s.uc.RoleRouterTree(ctx, roleIDs...)
	if err != nil {
		return nil, err
	}
	resp := &pb.RouterReply{}
	for _, route := range tree {
		var routeResp *pb.RouterGroup
		_ = copier.Copy(routeResp, &route)
		resp.Routers = append(resp.Routers, routeResp)
	}
	return resp, nil
}
func (s *ResourceService) GetRouteTreeByRole(ctx context.Context, req *pb.IDRequest) (*pb.RouterReply, error) {
	tree, err := s.uc.RoleRouterTree(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp := &pb.RouterReply{}
	for _, route := range tree {
		var routeResp *pb.RouterGroup
		_ = copier.Copy(routeResp, &route)
		resp.Routers = append(resp.Routers, routeResp)
	}
	return resp, nil
}
func (s *ResourceService) EditRoutePolicy(ctx context.Context, req *pb.RoutePolicyRequest) (*pb.EmptyReply, error) {
	var routers []biz.RouterGroup
	for _, r := range req.Routers {
		var router biz.RouterGroup
		_ = copier.Copy(&router, r)
		routers = append(routers, router)
	}

	if err := s.uc.EditRouters(ctx, req.Id, routers); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
