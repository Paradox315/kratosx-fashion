package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type ResourceActionRepo struct {
	dao *Data
	log *log.Helper
}

func NewResourceActionRepo(dao *Data, logger log.Logger) biz.ResourceActionRepo {
	return &ResourceActionRepo{
		dao: dao,
		log: log.NewHelper(logger),
	}
}

func (r *ResourceActionRepo) Select(ctx context.Context, u uint64) (*model.ResourceAction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceActionRepo) SelectByMenuID(ctx context.Context, u uint64) ([]*model.ResourceAction, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceActionRepo) Insert(ctx context.Context, action *model.ResourceAction) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceActionRepo) Update(ctx context.Context, action *model.ResourceAction) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceActionRepo) Delete(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceActionRepo) DeleteByMenuIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}
