package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"

	pb "kratosx-fashion/api/system/v1"
)

type RoleService struct {
	pb.UnimplementedRoleServer

	uc  *biz.RoleUsecase
	log *log.Helper
}

func NewRoleService(uc *biz.RoleUsecase, logger log.Logger) *RoleService {
	return &RoleService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *RoleService) CreateRole(ctx context.Context, req *pb.RoleRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.RoleRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *RoleService) UpdateRoleStatus(ctx context.Context, req *pb.IDRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.IDRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *RoleService) GetRole(ctx context.Context, req *pb.IDRequest) (*pb.RoleReply, error) {
	return &pb.RoleReply{}, nil
}
func (s *RoleService) ListRole(ctx context.Context, req *pb.ListRequest) (*pb.ListRoleReply, error) {
	return &pb.ListRoleReply{}, nil
}
