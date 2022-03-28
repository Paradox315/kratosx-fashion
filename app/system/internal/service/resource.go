package service

import (
	"context"

	pb "kratosx-fashion/api/system/v1"
)

type ResourceService struct {
	pb.UnimplementedResourceServer
}

func NewResourceService() *ResourceService {
	return &ResourceService{}
}

func (s *ResourceService) CreateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *ResourceService) UpdateMenu(ctx context.Context, req *pb.MenuRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *ResourceService) DeleteMenu(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *ResourceService) GetMenuTree(ctx context.Context, req *pb.IDRequest) (*pb.MenuReply, error) {
	return &pb.MenuReply{}, nil
}
func (s *ResourceService) GetMenuTreeByRole(ctx context.Context, req *pb.IDRequest) (*pb.MenuReply, error) {
	return &pb.MenuReply{}, nil
}
func (s *ResourceService) GetRouteTree(ctx context.Context, req *pb.EmptyRequest) (*pb.RouterReply, error) {
	return &pb.RouterReply{}, nil
}
