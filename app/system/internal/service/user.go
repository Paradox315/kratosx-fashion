package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/pkg/option"
	"kratosx-fashion/pkg/pagination"
	"strconv"
	"strings"

	pb "kratosx-fashion/api/system/v1"
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
	id, err := s.uc.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	id, err := s.uc.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.IDReply{
		Id: id,
	}, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *pb.PasswordRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.IDRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.IDsRequest) (*pb.EmptyReply, error) {
	var uids []uint
	for _, uid := range strings.Split(req.Ids, ",") {
		id, _ := strconv.ParseUint(uid, 10, 64)
		uids = append(uids, uint(id))
	}
	if err := s.uc.DeleteUsers(ctx, uids); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.IDRequest) (*pb.UserReply, error) {
	uid, _ := strconv.ParseUint(req.Id, 10, 64)
	return s.uc.SelectUser(ctx, uint(uid))
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListSearchRequest) (*pb.ListUserReply, error) {
	limit, offset := pagination.Parse(req.PageNum, req.PageSize)
	where, order, args := option.Parse(req.Query...)
	return s.uc.ListUser(ctx, limit, offset, biz.SQLOption{
		Where: where,
		Order: order,
		Args:  args,
	})
}
