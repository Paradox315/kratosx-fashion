package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
)

type RoleRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewRoleRepo(data *data.Data, logger log.Logger) biz.RoleRepo {
	return &RoleRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "role")),
		baseRepo: linq.Use(data.DB),
	}
}

func (r *RoleRepo) Select(ctx context.Context, id uint) (*model.Role, error) {
	rr := r.baseRepo.Role
	return rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
}

func (r *RoleRepo) SelectByIDs(ctx context.Context, rids []uint) ([]*model.Role, error) {
	rr := r.baseRepo.Role
	return rr.WithContext(ctx).Where(rr.ID.In(rids...)).Find()
}

func (r *RoleRepo) List(ctx context.Context, limit, offset int) (roles []*model.Role, total int64, err error) {
	err = r.dao.DB.Limit(limit).Offset(offset).Count(&total).Find(&roles).Error
	return
}

func (r *RoleRepo) Insert(ctx context.Context, role ...*model.Role) error {
	rr := r.baseRepo.Role
	return rr.WithContext(ctx).Create(role...)
}

func (r *RoleRepo) Update(ctx context.Context, role *model.Role) error {
	rr := r.baseRepo.Role
	_, err := rr.WithContext(ctx).Updates(role)
	return err
}

func (r *RoleRepo) DeleteByIDs(ctx context.Context, rids []uint) error {
	rr := r.baseRepo.Role
	_, err := rr.WithContext(ctx).Where(rr.ID.In(rids...)).Delete()
	return err
}
