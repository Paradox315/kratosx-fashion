package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"

	pb "kratosx-fashion/api/system/v1"
)

type MenuService struct {
	pb.UnimplementedMenuServer

	uc  *biz.MenuUsecase
	log *log.Helper
}

func NewMenuService(uc *biz.MenuUsecase, logger log.Logger) *MenuService {
	return &MenuService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *MenuService) CreateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *MenuService) UpdateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *MenuService) UpdateMenuStatus(ctx context.Context, req *pb.IDRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *MenuService) DeleteMenu(ctx context.Context, req *pb.IDRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *MenuService) GetMenu(ctx context.Context, req *pb.IDRequest) (*pb.MenuReply, error) {
	return &pb.MenuReply{}, nil
}
func (s *MenuService) ListMenu(ctx context.Context, req *pb.ListRequest) (*pb.ListMenuReply, error) {
	return &pb.ListMenuReply{}, nil
}
func (s *MenuService) GetMenuTree(ctx context.Context, req *pb.TreeRequest) (*pb.MenuTreeReply, error) {
	return &pb.MenuTreeReply{}, nil
}
