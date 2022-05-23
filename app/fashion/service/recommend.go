package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	api "kratosx-fashion/api/fashion/v1"
	"kratosx-fashion/app/fashion/biz"
	"kratosx-fashion/pkg/pagination"

	pb "kratosx-fashion/api/fashion/v1"
)

type RecommendService struct {
	log *log.Helper
	uc  *biz.RecommendUsecase
}

func NewRecommendService(uc *biz.RecommendUsecase, logger log.Logger) *RecommendService {
	return &RecommendService{
		log: log.NewHelper(log.With(logger, "service", "recommend")),
		uc:  uc,
	}
}

func (s *RecommendService) GetPopular(ctx context.Context, req *pb.ListRequest) (*pb.RecommendReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	resp, err := s.uc.GetPopular(ctx, limit, offset)
	if err != nil {
		err = api.ErrorRecommendFailed("获取热门推荐失败")
		return nil, err
	}
	return resp, nil
}
func (s *RecommendService) GetLatest(ctx context.Context, req *pb.ListRequest) (*pb.RecommendReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	resp, err := s.uc.GetLatest(ctx, limit, offset)
	if err != nil {
		err = api.ErrorRecommendFailed("获取最新推荐失败")
		return nil, err
	}
	return resp, nil
}
func (s *RecommendService) GetClothesNeighbors(ctx context.Context, req *pb.ListIDRequest) (*pb.RecommendReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	resp, err := s.uc.GetClothesNeighbors(ctx, req.Id, limit, offset)
	if err != nil {
		err = api.ErrorRecommendFailed("获取服饰相似推荐失败")
		return nil, err
	}
	return resp, nil
}
func (s *RecommendService) GetUserNeighbors(ctx context.Context, req *pb.ListRequest) (*pb.RecommendReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	resp, err := s.uc.GetUserNeighbors(ctx, limit, offset)
	if err != nil {
		err = api.ErrorRecommendFailed("获取用户相似推荐失败")
		return nil, err
	}
	return resp, nil
}
func (s *RecommendService) GetUserRecommend(ctx context.Context, req *pb.ListRequest) (*pb.RecommendReply, error) {
	limit, offset := pagination.Parse(req.Current, req.PageSize)
	resp, err := s.uc.GetUserRecommend(ctx, limit, offset)
	if err != nil {
		err = api.ErrorRecommendFailed("获取用户推荐失败")
		return nil, err
	}
	return resp, nil
}
