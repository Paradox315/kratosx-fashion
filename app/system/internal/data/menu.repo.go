package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type MenuRepo struct {
	dao *Data
	log *log.Helper
}

func NewMenuRepo(dao *Data, logger log.Logger) biz.MenuRepo {
	return &MenuRepo{
		dao: dao,
		log: log.NewHelper(logger),
	}
}

func (m *MenuRepo) Select(ctx context.Context, u uint64) (*model.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuRepo) SelectByIDs(ctx context.Context, uint64s []uint64) ([]*model.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuRepo) List(ctx context.Context, request *pb.ListRequest, option ...*pb.QueryOption) ([]*model.Menu, uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuRepo) Insert(ctx context.Context, menu *model.Menu) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuRepo) Update(ctx context.Context, menu *model.Menu) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuRepo) Delete(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuRepo) DeleteByIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}
