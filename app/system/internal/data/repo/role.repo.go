package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
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
	role, err := rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "role.repo.Select")
		return nil, err
	}
	return role, nil
}

func (r *RoleRepo) SelectByIDs(ctx context.Context, rids []uint) ([]*model.Role, error) {
	rr := r.baseRepo.Role
	roles, err := rr.WithContext(ctx).Where(rr.ID.In(rids...)).Find()
	if err != nil {
		err = errors.Wrap(err, "role.repo.SelectByIDs")
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepo) List(ctx context.Context, limit, offset int) (roles []*model.Role, total int64, err error) {
	rr := r.baseRepo.Role
	roles, total, err = rr.WithContext(ctx).FindByPage(offset, limit)
	if err != nil {
		err = errors.Wrap(err, "role.repo.List")
		return
	}
	return
}

func (r *RoleRepo) Insert(ctx context.Context, role ...*model.Role) error {
	rr := r.baseRepo.Role
	err := rr.WithContext(ctx).Create(role...)
	if err != nil {
		err = errors.Wrap(err, "role.repo.Insert")
		return err
	}
	return nil
}

func (r *RoleRepo) Update(ctx context.Context, role *model.Role) error {
	rr := r.baseRepo.Role
	if _, err := rr.WithContext(ctx).Where(rr.ID.Eq(role.ID)).Updates(role); err != nil {
		err = errors.Wrap(err, "role.repo.Update")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleRepo) DeleteByIDs(ctx context.Context, rids []uint) error {
	rr := r.baseRepo.Role
	if _, err := rr.WithContext(ctx).Where(rr.ID.In(rids...)).Delete(); err != nil {
		err = errors.Wrap(err, "role.repo.DeleteByIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}
