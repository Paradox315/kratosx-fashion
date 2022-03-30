package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
)

type ResourceActionRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewResourceActionRepo(dao *data.Data, logger log.Logger) biz.ResourceActionRepo {
	return &ResourceActionRepo{
		dao:      dao,
		log:      log.NewHelper(log.With(logger, "repo", "resource_action")),
		baseRepo: linq.Use(dao.DB),
	}
}

func (r *ResourceActionRepo) Select(ctx context.Context, id uint) (*model.ResourceAction, error) {
	ar := r.baseRepo.ResourceAction
	return ar.WithContext(ctx).Where(ar.ID.Eq(id)).First()
}

func (r *ResourceActionRepo) SelectByMenuID(ctx context.Context, menuID uint64) ([]*model.ResourceAction, error) {
	ar := r.baseRepo.ResourceAction
	return ar.WithContext(ctx).Where(ar.MenuID.Eq(menuID)).Find()
}

func (r *ResourceActionRepo) Insert(ctx context.Context, action ...*model.ResourceAction) error {
	ar := r.baseRepo.ResourceAction
	return ar.WithContext(ctx).Create(action...)
}

func (r *ResourceActionRepo) Update(ctx context.Context, action *model.ResourceAction) error {
	ar := r.baseRepo.ResourceAction
	_, err := ar.WithContext(ctx).Updates(action)
	return err
}

func (r *ResourceActionRepo) Delete(ctx context.Context, id uint) error {
	ar := r.baseRepo.ResourceAction
	_, err := ar.WithContext(ctx).Where(ar.ID.Eq(id)).Delete()
	return err
}

func (r *ResourceActionRepo) DeleteByMenuIDs(ctx context.Context, mids []uint64) error {
	ar := r.baseRepo.ResourceAction
	_, err := ar.WithContext(ctx).Where(ar.MenuID.In(mids...)).Delete()
	return err
}
