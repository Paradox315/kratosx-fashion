package service

import (
	"context"

	pb "kratosx-fashion/api/system/v1"
)

type RoleService struct {
	pb.UnimplementedRoleServer
}

func NewRoleService() *RoleService {
	return &RoleService{}
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
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *RoleService) GetRole(ctx context.Context, req *pb.IDRequest) (*pb.RoleReply, error) {
	return &pb.RoleReply{}, nil
}
func (s *RoleService) ListRole(ctx context.Context, req *pb.ListRequest) (*pb.ListRoleReply, error) {
	return &pb.ListRoleReply{}, nil
}
