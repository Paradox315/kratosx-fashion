package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/pkg/pagination"
	"kratosx-fashion/pkg/xcast"
	"strings"

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
		log: log.NewHelper(log.With(logger, "service", "role")),
	}
}

func (s *RoleService) CreateRole(ctx context.Context, req *pb.RoleRequest) (*pb.IDReply, error) {
	id, err := s.uc.Save(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.RoleRequest) (*pb.IDReply, error) {
	id, err := s.uc.Edit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	ids := strings.Split(req.Ids, ",")
	err := s.uc.Remove(ctx, xcast.ToUintSlice(ids))
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *RoleService) GetRole(ctx context.Context, req *pb.IDRequest) (*pb.RoleReply, error) {
	return s.uc.Get(ctx, cast.ToUint(req.Id))
}
func (s *RoleService) ListRole(ctx context.Context, req *pb.ListRequest) (*pb.ListRoleReply, error) {
	limit, offset := pagination.Parse(req.PageNum, req.PageSize)
	return s.uc.Page(ctx, limit, offset)
}
