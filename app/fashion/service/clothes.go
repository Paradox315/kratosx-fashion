package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/fashion/biz"

	pb "kratosx-fashion/api/fashion/v1"
)

type ClothesService struct {
	log *log.Helper
	uc  *biz.ClothesUsecase
}

func NewClothesService(uc *biz.ClothesUsecase, logger log.Logger) *ClothesService {
	return &ClothesService{
		log: log.NewHelper(log.With(logger, "service", "clothes")),
		uc:  uc,
	}
}

func (s *ClothesService) GetClothes(ctx context.Context, req *pb.IDRequest) (*pb.ClothesReply, error) {
	return s.uc.Get(ctx, req.Id)
}
func (s *ClothesService) CreateClothes(ctx context.Context, req *pb.ClothesRequest) (*pb.EmptyReply, error) {

	if err := s.uc.Save(ctx, req); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *ClothesService) UpdateClothes(ctx context.Context, req *pb.ClothesRequest) (*pb.EmptyReply, error) {
	if err := s.uc.Edit(ctx, req); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *ClothesService) DeleteClothes(ctx context.Context, req *pb.IDRequest) (*pb.EmptyReply, error) {
	if err := s.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
