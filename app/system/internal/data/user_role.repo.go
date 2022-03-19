package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/model"
)

type userRoleRepo struct {
	dao *data.Data
	log *log.Helper
}

func NewUserRoleRepo(data *data.Data, logger log.Logger) biz.UserRoleRepo {
	return &userRoleRepo{
		dao: data,
		log: log.NewHelper(logger),
	}
}

func (u *userRoleRepo) Select(ctx context.Context, u2 uint64) (*model.UserRole, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) SelectAll(ctx context.Context) ([]*model.UserRole, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) SelectAllByUserID(ctx context.Context, u2 uint64) ([]*model.UserRole, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) SelectRoleIDByUserID(ctx context.Context, u2 uint64) ([]uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) Insert(ctx context.Context, role *model.UserRole) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) Update(ctx context.Context, role *model.UserRole) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) UpdateStatus(ctx context.Context, u3 uint64, u2 uint8) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) Delete(ctx context.Context, u2 uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) DeleteByUserID(ctx context.Context, u2 uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) DeleteByUserIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) DeleteByRoleID(ctx context.Context, u2 uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) DeleteByRoleIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRoleRepo) ExistByUserIDAndRoleID(ctx context.Context, u3 uint64, u2 uint64) (bool, error) {
	//TODO implement me
	panic("implement me")
}
