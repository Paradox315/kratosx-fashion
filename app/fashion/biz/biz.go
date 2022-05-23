package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/google/wire"
	pb "kratosx-fashion/api/fashion/v1"
	"kratosx-fashion/app/fashion/model"
)

var codec = encoding.GetCodec("json")

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewClothesUsecase,
	NewRecommendUsecase,
)

type RecommendRepo interface {
	SelectUserNeighbors(context.Context, string, int, int) ([]string, error)
	SelectClothesNeighbors(context.Context, string, int, int) ([]string, error)
	SelectPopular(context.Context, int, int) ([]model.Item, error)
	SelectLatest(context.Context, int, int) ([]model.Item, error)
	SelectUserRecommend(context.Context, string, int, int) ([]string, error)
}

type ClothesRepo interface {
	Insert(context.Context, *model.Clothes) error
	Select(context.Context, string) (*model.Clothes, error)
	SelectByIDs(context.Context, []string) ([]*model.Clothes, error)
	Update(context.Context, *model.Clothes) error
	Delete(context.Context, string) error
	TryOn(context.Context, *pb.TryOnRequest) (*pb.TryOnReply, error)
}

type FeedbackRepo interface {
	Insert(context.Context, *model.Feedback) error
}
