package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type ResourceMenuRepo struct {
	dao *Data
	log *log.Helper
}

func NewResourceMenuRepo(dao *Data, logger log.Logger) biz.ResourceMenuRepo {
	return &ResourceMenuRepo{
		dao: dao,
		log: log.NewHelper(logger),
	}
}

func (r *ResourceMenuRepo) Select(ctx context.Context, u uint64) (*model.ResourceMenu, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceMenuRepo) SelectByIDs(ctx context.Context, uints []uint) ([]*model.ResourceMenu, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceMenuRepo) Insert(ctx context.Context, menu *model.ResourceMenu) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceMenuRepo) Update(ctx context.Context, menu *model.ResourceMenu) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceMenuRepo) DeleteByIDs(ctx context.Context, uints []uint) error {
	//TODO implement me
	panic("implement me")
}
