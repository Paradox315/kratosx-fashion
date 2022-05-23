package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/biz"
	"kratosx-fashion/pkg/pagination"
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
		Id: uint64(id),
	}, nil
}
func (s *ResourceService) UpdateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.EmptyReply, error) {
	if err := s.uc.EditMenu(ctx, req); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *ResourceService) DeleteMenu(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	ids := strings.Split(req.Ids, ",")
	if err := s.uc.RemoveMenu(ctx, xcast.ToUintSlice(ids)); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *ResourceService) GetMenuTree(ctx context.Context, req *pb.EmptyRequest) (*pb.MenuReply, error) {
	tree, err := s.uc.MenuTree(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.MenuReply{}
	for _, menu := range tree {
		menuResp := &pb.Menu{}
		_ = copier.Copy(&menuResp, &menu)
		resp.List = append(resp.List, menuResp)
	}
	return resp, nil
}
func (s *ResourceService) GetMenuByRole(ctx context.Context, req *pb.IDRequest) (*pb.MenuReply, error) {
	tree, err := s.uc.RoleMenuList(ctx, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	resp := &pb.MenuReply{}
	for _, menu := range tree {
		menuResp := &pb.Menu{}
		_ = copier.Copy(&menuResp, &menu)
		resp.List = append(resp.List, menuResp)
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
		menuResp := &pb.Menu{}
		_ = copier.Copy(&menuResp, &menu)
		resp.List = append(resp.List, menuResp)
	}
	return resp, nil
}
func (s *ResourceService) GetRouteTree(ctx context.Context, req *pb.EmptyRequest) (*pb.RouterReply, error) {
	tree, err := s.uc.RouterTree(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.RouterReply{}
	for _, route := range tree {
		routeResp := &pb.RouterGroup{}
		_ = copier.Copy(&routeResp, &route)
		resp.Routers = append(resp.Routers, routeResp)
	}
	return resp, nil
}
func (s *ResourceService) GetRouteByRole(ctx context.Context, req *pb.IDRequest) (*pb.RouterReply, error) {
	tree, err := s.uc.RoleRouterList(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp := &pb.RouterReply{}
	for _, route := range tree {
		routeResp := &pb.Router{}
		_ = copier.Copy(&routeResp, &route)
		resp.RoleRouters = append(resp.RoleRouters, routeResp)
	}
	return resp, nil
}
func (s *ResourceService) ListMenu(ctx context.Context, req *pb.ListRequest) (*pb.ListMenuReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	menus, total, err := s.uc.MenuPage(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListMenuReply{}
	resp.Total = total
	for _, menu := range menus {
		menuResp := &pb.Menu{}
		_ = copier.Copy(&menuResp, &menu)
		resp.List = append(resp.List, menuResp)
	}
	return resp, nil
}
func (s *ResourceService) ListRoute(ctx context.Context, req *pb.ListRequest) (*pb.ListRouterReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	routes, total, err := s.uc.RouterPage(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListRouterReply{}
	resp.Total = total
	for _, route := range routes {
		routeResp := &pb.RouterGroup{}
		_ = copier.Copy(&routeResp, &route)
		resp.List = append(resp.List, routeResp)
	}
	return resp, nil
}
func (s *ResourceService) GetActionByRole(ctx context.Context, req *pb.IDRequest) (*pb.ListMenuActionReply, error) {
	actions, err := s.uc.RoleActionList(ctx, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	resp := &pb.ListMenuActionReply{}
	for _, action := range actions {
		resp.List = append(resp.List, action)
	}
	return resp, nil
}
