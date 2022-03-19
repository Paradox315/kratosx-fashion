package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type RoleRepo struct {
	dao *Data
	log *log.Helper
}

func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &RoleRepo{
		dao: data,
		log: log.NewHelper(logger),
	}
}

func (r *RoleRepo) Select(ctx context.Context, u uint64) (*model.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) SelectByIDs(ctx context.Context, uint64s []uint64) ([]*model.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) List(ctx context.Context, request *pb.ListRequest, option ...*pb.QueryOption) ([]*model.Role, uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) Insert(ctx context.Context, role *model.Role) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) Update(ctx context.Context, role *model.Role) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) UpdateStatus(ctx context.Context, u uint64, u2 uint8) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) Delete(ctx context.Context, u uint64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RoleRepo) DeleteByIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}
