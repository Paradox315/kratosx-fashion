package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/option"
	"kratosx-fashion/pkg/pagination"
	"kratosx-fashion/pkg/xcast"
	"strconv"
	"strings"
)

type UserService struct {
	pb.UnimplementedUserServer

	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "service", "user")),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	id, err := s.uc.Save(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	id, err := s.uc.Edit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *pb.PasswordRequest) (*pb.IDReply, error) {
	err := s.uc.EditPassword(ctx, req.OldPassword, req.NewPassword, cast.ToUint(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: req.Id,
	}, nil
}
func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.StatusRequest) (*pb.IDReply, error) {
	err := s.uc.EditStatus(ctx, cast.ToUint(req.Id), model.UserStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: req.Id,
	}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	if err := s.uc.Remove(ctx, xcast.ToUintSlice[string](strings.Split(req.Ids, ","))); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.IDRequest) (*pb.UserReply, error) {
	uid, _ := strconv.ParseUint(req.Id, 10, 64)
	return s.uc.Get(ctx, uint(uid))
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListSearchRequest) (*pb.ListUserReply, error) {
	limit, offset := pagination.Parse(req.PageNum, req.PageSize)
	where, order, args := option.Parse(req.Query...)
	return s.uc.Search(ctx, limit, offset, biz.SQLOption{
		Where: where,
		Order: order,
		Args:  args,
	})
}
