package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"

	pb "kratosx-fashion/api/system/v1"
)

type PubService struct {
	pb.UnimplementedPubServer

	uc  *biz.PublicUsecase
	log *log.Helper
}

func NewPubService(uc *biz.PublicUsecase, logger log.Logger) *PubService {
	return &PubService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *PubService) Generate(ctx context.Context, req *pb.EmptyRequest) (*pb.CaptchaReply, error) {
	return &pb.CaptchaReply{}, nil
}
func (s *PubService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *PubService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
func (s *PubService) Logout(ctx context.Context, req *pb.EmptyRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
func (s *PubService) RetrievePwd(ctx context.Context, req *pb.RetrieveRequest) (*pb.EmptyReply, error) {
	return &pb.EmptyReply{}, nil
}
