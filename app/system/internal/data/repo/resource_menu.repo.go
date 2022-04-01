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

type ResourceMenuRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewResourceMenuRepo(dao *data.Data, logger log.Logger) biz.ResourceMenuRepo {
	return &ResourceMenuRepo{
		dao:      dao,
		log:      log.NewHelper(log.With(logger, "repo", "resource_menu")),
		baseRepo: linq.Use(dao.DB),
	}
}

func (r *ResourceMenuRepo) Select(ctx context.Context, id uint) (*model.ResourceMenu, error) {
	rr := r.baseRepo.ResourceMenu
	menu, err := rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "resource_menu.repo.Select")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return menu, nil
}

func (r *ResourceMenuRepo) SelectByIDs(ctx context.Context, ids []uint) ([]*model.ResourceMenu, error) {
	rr := r.baseRepo.ResourceMenu
	menus, err := rr.WithContext(ctx).Where(rr.ID.In(ids...)).Find()
	if err != nil {
		err = errors.Wrap(err, "resource_menu.repo.SelectByIDs")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return menus, nil
}

func (r *ResourceMenuRepo) Insert(ctx context.Context, menu ...*model.ResourceMenu) error {
	rr := r.baseRepo.ResourceMenu
	if err := rr.WithContext(ctx).Create(menu...); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.Insert")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *ResourceMenuRepo) Update(ctx context.Context, menu *model.ResourceMenu) error {
	rr := r.baseRepo.ResourceMenu
	if _, err := rr.WithContext(ctx).Where(rr.ID.Eq(menu.ID)).Updates(menu); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.Update")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *ResourceMenuRepo) DeleteByIDs(ctx context.Context, ids []uint) error {
	rr := r.baseRepo.ResourceMenu
	if _, err := rr.WithContext(ctx).Where(rr.ID.In(ids...)).Delete(); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.DeleteByIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}
