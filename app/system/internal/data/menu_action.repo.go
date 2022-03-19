package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type MenuActionRepo struct {
	dao *Data
	log *log.Helper
}

func NewMenuActionRepo(dao *Data, logger log.Logger) biz.MenuActionRepo {
	return &MenuActionRepo{
		dao: dao,
		log: log.NewHelper(logger),
	}
}

func (m *MenuActionRepo) Select(ctx context.Context, u uint64) (*model.MenuAction, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) SelectByMenuID(ctx context.Context, u uint64) ([]*model.MenuAction, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) ExistByMenuIDAndActionID(ctx context.Context, u uint64, u2 uint64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) Insert(ctx context.Context, action *model.MenuAction) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) Update(ctx context.Context, action *model.MenuAction) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) Delete(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) DeleteByMenuID(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuActionRepo) DeleteByMenuIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}
