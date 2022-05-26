package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/fashion/v1"
	"kratosx-fashion/app/fashion/do"
	"kratosx-fashion/app/fashion/model"
	"kratosx-fashion/pkg/ctxutil"
	"math"
	"strings"
	"time"
)

type ClothesUsecase struct {
	log          *log.Helper
	clothesRepo  ClothesRepo
	feedbackRepo FeedbackRepo
}

func NewClothesUsecase(clothesRepo ClothesRepo, feedbackRepo FeedbackRepo, logger log.Logger) *ClothesUsecase {
	return &ClothesUsecase{
		log:          log.NewHelper(log.With(logger, "biz", "clothes")),
		clothesRepo:  clothesRepo,
		feedbackRepo: feedbackRepo,
	}
}
func (u *ClothesUsecase) buildClothesDto(ctx context.Context, cpo *model.Clothes) (clothes *pb.ClothesReply, err error) {
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

func (u *ClothesUsecase) buildClothesPo(ctx context.Context, clothes *pb.ClothesRequest) (cpo *model.Clothes, err error) {
	comment := do.ClothesComment{
		Description: clothes.Description,
		Price:       cast.ToString(clothes.Price),
		Region:      clothes.Region,
		Image:       clothes.Image,
		Colors:      strings.Join(clothes.Colors, ","),
	}
	bytes, err := codec.Marshal(&comment)
	if err != nil {
		return
	}
	regions := strings.Split(clothes.Region, "-")
	cpo = &model.Clothes{
		ItemId:     clothes.Id,
		IsHidden:   false,
		Categories: []string{"clothes"},
		Timestamp:  time.Now(),
		Labels:     []string{clothes.Type, clothes.Brand, clothes.Style, cast.ToString(int(math.Floor(float64(clothes.Price + 0/5)))), regions[0], regions[1]},
		Comment:    string(bytes),
	}
	return
}

func (u *ClothesUsecase) Get(ctx context.Context, id string) (clothes *pb.ClothesReply, err error) {
	cpo, err := u.clothesRepo.Select(ctx, id)
	if err != nil {
		return
	}
	clothes, err = u.buildClothesDto(ctx, cpo)
	if err != nil {
		return nil, err
	}
	if len(id) == 0 {
		return
	}
	if err = u.feedbackRepo.Insert(ctx, &model.Feedback{
		FeedbackType: "view",
		UserId:       cast.ToString(ctxutil.GetUid(ctx)),
		ItemId:       id,
		Timestamp:    time.Now().Format(time.RFC3339),
	}); err != nil {
		u.log.WithContext(ctx).Error(err)
	}
	return
}

func (u *ClothesUsecase) Save(ctx context.Context, clothes *pb.ClothesRequest) (err error) {
	cpo, err := u.buildClothesPo(ctx, clothes)
	if err != nil {
		return err
	}
	return u.clothesRepo.Insert(ctx, cpo)
}

func (u *ClothesUsecase) Edit(ctx context.Context, clothes *pb.ClothesRequest) (err error) {
	cpo, err := u.buildClothesPo(ctx, clothes)
	if err != nil {
		return
	}
	return u.clothesRepo.Update(ctx, cpo)
}

func (u *ClothesUsecase) Delete(ctx context.Context, id string) (err error) {
	return u.clothesRepo.Delete(ctx, id)
}

func (u *ClothesUsecase) TryOn(ctx context.Context, req *pb.TryOnRequest) (resp *pb.TryOnReply, err error) {
	resp, err = u.clothesRepo.TryOn(ctx, req)
	if err != nil {
		return
	}
	if len(req.ClothesId) == 0 {
		return
	}
	if err = u.feedbackRepo.Insert(ctx, &model.Feedback{
		FeedbackType: "tryon",
		UserId:       cast.ToString(ctxutil.GetUid(ctx)),
		ItemId:       req.ClothesId,
		Timestamp:    time.Now().Format(time.RFC3339),
	}); err != nil {
		u.log.WithContext(ctx).Error(err)
	}
	return
}
