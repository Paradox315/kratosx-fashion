package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
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
	return rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
}

func (r *ResourceMenuRepo) SelectByIDs(ctx context.Context, ids []uint) ([]*model.ResourceMenu, error) {
	rr := r.baseRepo.ResourceMenu
	return rr.WithContext(ctx).Where(rr.ID.In(ids...)).Find()
}

func (r *ResourceMenuRepo) Insert(ctx context.Context, menu ...*model.ResourceMenu) error {
	rr := r.baseRepo.ResourceMenu
	return rr.WithContext(ctx).Create(menu...)
}

func (r *ResourceMenuRepo) Update(ctx context.Context, menu *model.ResourceMenu) error {
	rr := r.baseRepo.ResourceMenu
	_, err := rr.WithContext(ctx).Updates(menu)
	return err
}

func (r *ResourceMenuRepo) DeleteByIDs(ctx context.Context, ids []uint) error {
	rr := r.baseRepo.ResourceMenu
	_, err := rr.WithContext(ctx).Where(rr.ID.In(ids...)).Delete()
	return err
}
