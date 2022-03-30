package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "kratosx-fashion/api/system/v1"
)

type ResourceUsecase struct {
	menuRepo  ResourceMenuRepo
	actRepo   ResourceActionRepo
	routeRepo ResourceRouterRepo
	tx        Transaction
	log       *log.Helper
}

func NewResourceUsecase(menuRepo ResourceMenuRepo, actRepo ResourceActionRepo, routeRepo ResourceRouterRepo, tx Transaction, logger log.Logger) *ResourceUsecase {
	return &ResourceUsecase{
		menuRepo:  menuRepo,
		actRepo:   actRepo,
		routeRepo: routeRepo,
		tx:        tx,
		log:       log.NewHelper(log.With(logger, "biz", "resource")),
	}
}

func (r *ResourceUsecase) SaveMenu(ctx context.Context, menu *pb.MenuRequest) (id string, err error) {
	panic("implement me")
}

func (r *ResourceUsecase) EditMenu(ctx context.Context, menu *pb.MenuRequest) (id string, err error) {
	panic("implement me")
}

func (r *ResourceUsecase) RemoveMenu(ctx context.Context, ids []uint) (err error) {
	panic("implement me")
}

func (r *ResourceUsecase) UserMenuTree(ctx context.Context, uid uint) (tree *pb.MenuReply, err error) {
	panic("implement me")
}

func (r *ResourceUsecase) RoleMenuTree(ctx context.Context, rid uint) (tree *pb.MenuReply, err error) {
	panic("implement me")
}

func (r *ResourceUsecase) UserRouterTree(ctx context.Context, uid uint) (tree *pb.RouterReply, err error) {
	panic("implement me")
}

func (r *ResourceUsecase) RoleRouterTree(ctx context.Context, rid uint) (tree *pb.RouterReply, err error) {
	panic("implement me")
}

func (r ResourceUsecase) EditRouters(ctx context.Context, req *pb.RoutePolicyRequest) (err error) {
	panic("implement me")
}
