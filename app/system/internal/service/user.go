package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"

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
		log: log.NewHelper(logger),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, req *pb.PasswordRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.IDRequest) (*pb.IDReply, error) {
	return &pb.IDReply{}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.IDRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.IDRequest) (*pb.UserReply, error) {
	return &pb.UserReply{}, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListRequest) (*pb.ListUserReply, error) {
	return &pb.ListUserReply{}, nil
}
