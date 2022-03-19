package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/model"
)

type RoleMenuRepo struct {
	dao *data.Data
	log *log.Helper
}

func NewRoleMenuRepo(data *data.Data, logger log.Logger) biz.RoleMenuRepo {
	return &RoleMenuRepo{
		dao: data,
		log: log.NewHelper(logger),
	}
}
func (r *RoleMenuRepo) Select(ctx context.Context, u uint64) (*model.RoleMenu, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) SelectByRoleID(ctx context.Context, u uint64) ([]*model.RoleMenu, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) ExistByRoleIDAndMenuID(ctx context.Context, u uint64, u2 uint64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) Insert(ctx context.Context, menu *model.RoleMenu) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) Update(ctx context.Context, menu *model.RoleMenu) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) Delete(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) DeleteByRoleID(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) DeleteByRoleIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) DeleteByMenuID(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleMenuRepo) DeleteByMenuIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}
