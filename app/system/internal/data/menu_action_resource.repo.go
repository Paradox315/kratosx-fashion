package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type MenuActionResourceRepo struct {
	dao *Data
	log *log.Helper
}

func NewMenuActionResourceRepo(dao *Data, logger log.Logger) biz.MenuActionResourceRepo {
	return &MenuActionResourceRepo{
		dao: dao,
		log: log.NewHelper(logger),
	}
}

func (m *MenuActionResourceRepo) Select(ctx context.Context, u uint64) (*model.MenuActionResource, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) SelectByActionID(ctx context.Context, u uint64) ([]*model.MenuActionResource, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) Insert(ctx context.Context, resource *model.MenuActionResource) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) Update(ctx context.Context, resource *model.MenuActionResource) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) Delete(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) DeleteByIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) DeleteByActionID(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionResourceRepo) DeleteByActionIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}
