package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	api "kratosx-fashion/api/fashion/v1"
	"kratosx-fashion/app/fashion/biz"

	pb "kratosx-fashion/api/fashion/v1"
)

type TryOnService struct {
	log *log.Helper
	uc  *biz.ClothesUsecase
}

func NewTryOnService(uc *biz.ClothesUsecase, logger log.Logger) *TryOnService {
	return &TryOnService{
		log: log.NewHelper(log.With(logger, "service", "tryon")),
		uc:  uc,
	}
}

func (s *TryOnService) TryOnClothes(ctx context.Context, req *pb.TryOnRequest) (*pb.TryOnReply, error) {
	resp, err := s.uc.TryOn(ctx, req)
	if err != nil {
		err = api.ErrorTryonFailed("试穿失败")
		return nil, err
	}
	return resp, nil
}
