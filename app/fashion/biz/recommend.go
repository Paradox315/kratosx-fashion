package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/fashion/v1"
	"kratosx-fashion/app/fashion/do"
	"kratosx-fashion/app/fashion/model"
	"kratosx-fashion/pkg/ctxutil"
	"strings"
)

type RecommendUsecase struct {
	log           *log.Helper
	recommendRepo RecommendRepo
	clothesRepo   ClothesRepo
}

func NewRecommendUsecase(recommendRepo RecommendRepo, clothesRepo ClothesRepo, logger log.Logger) *RecommendUsecase {
	return &RecommendUsecase{
		recommendRepo: recommendRepo,
		clothesRepo:   clothesRepo,
		log:           log.NewHelper(log.With(logger, "biz", "recommend")),
	}
}
func (u *RecommendUsecase) buildClothesDto(ctx context.Context, cpo *model.Clothes) (clothes *pb.ClothesReply, err error) {
	var comment do.ClothesComment
	if err = codec.Unmarshal([]byte(cpo.Comment), &comment); err != nil {
		return
	}
	typ, brand, style := cpo.Labels[0], cpo.Labels[1], cpo.Labels[2]
	clothes = &pb.ClothesReply{
		Id:          cpo.ItemId,
		Type:        typ,
		Description: comment.Description,
		Image:       comment.Image,
		Brand:       brand,
		Style:       style,
		Region:      comment.Region,
		Time:        cpo.Timestamp.Format("2006-01-02 15:04:05"),
		Price:       cast.ToFloat32(comment.Price),
		Colors:      strings.Split(comment.Colors, ","),
	}
	return
}
func (u *RecommendUsecase) GetPopular(ctx context.Context, limit, offset int) (*pb.RecommendReply, error) {
	resp := &pb.RecommendReply{}
	items, err := u.recommendRepo.SelectPopular(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, item.Id)
	}
	cpos, err := u.clothesRepo.SelectByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, cpo := range cpos {
		var clothes *pb.ClothesReply
		clothes, err = u.buildClothesDto(ctx, cpo)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, clothes)
	}
	return resp, nil
}
func (u *RecommendUsecase) GetLatest(ctx context.Context, limit, offset int) (*pb.RecommendReply, error) {
	resp := &pb.RecommendReply{}
	items, err := u.recommendRepo.SelectLatest(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, item.Id)
	}
	cpos, err := u.clothesRepo.SelectByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, cpo := range cpos {
		var clothes *pb.ClothesReply
		clothes, err = u.buildClothesDto(ctx, cpo)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, clothes)
	}
	return resp, nil
}
func (u *RecommendUsecase) GetClothesNeighbors(ctx context.Context, id string, limit, offset int) (*pb.RecommendReply, error) {
	resp := &pb.RecommendReply{}
	items, err := u.recommendRepo.SelectClothesNeighbors(ctx, id, limit, offset)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, item)
	}
	cpos, err := u.clothesRepo.SelectByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, cpo := range cpos {
		var clothes *pb.ClothesReply
		clothes, err = u.buildClothesDto(ctx, cpo)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, clothes)
	}
	return resp, nil
}

func (u *RecommendUsecase) GetUserNeighbors(ctx context.Context, limit, offset int) (*pb.RecommendReply, error) {
	resp := &pb.RecommendReply{}
	uid := ctxutil.GetUid(ctx)
	items, err := u.recommendRepo.SelectUserNeighbors(ctx, cast.ToString(uid), limit, offset)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, item)
	}
	cpos, err := u.clothesRepo.SelectByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, cpo := range cpos {
		var clothes *pb.ClothesReply
		clothes, err = u.buildClothesDto(ctx, cpo)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, clothes)
	}
	return resp, nil
}
func (u *RecommendUsecase) GetUserRecommend(ctx context.Context, limit, offset int) (*pb.RecommendReply, error) {
	resp := &pb.RecommendReply{}
	uid := ctxutil.GetUid(ctx)
	items, err := u.recommendRepo.SelectUserRecommend(ctx, cast.ToString(uid), limit, offset)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, item)
	}
	cpos, err := u.clothesRepo.SelectByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, cpo := range cpos {
		var clothes *pb.ClothesReply
		clothes, err = u.buildClothesDto(ctx, cpo)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, clothes)
	}
	return resp, nil
}
