package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
)

type RoleResourceRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewRoleResourceRepo(data *data.Data, logger log.Logger) biz.RoleResourceRepo {
	return &RoleResourceRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "role_resource")),
		baseRepo: linq.Use(data.DB),
	}
}

func (r *RoleResourceRepo) Select(ctx context.Context, id uint) (*model.RoleResource, error) {
	rr := r.baseRepo.RoleResource
	rrs, err := rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "role_resource.repo.Select")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return rrs, nil
}

func (r *RoleResourceRepo) SelectByRoleID(ctx context.Context, rid uint64, resourceType ...model.ResourceType) ([]*model.RoleResource, error) {
	rr := r.baseRepo.RoleResource
	if len(resourceType) == 0 {
		return rr.WithContext(ctx).Where(rr.RoleID.Eq(rid)).Find()
	}
	typ := resourceType[0]
	list, err := rr.WithContext(ctx).Where(rr.RoleID.Eq(rid), rr.Type.Eq(uint8(typ))).Find()
	if err != nil {
		err = errors.Wrap(err, "role_resource.repo.SelectByRoleID")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return list, nil
}

func (r *RoleResourceRepo) Insert(ctx context.Context, resource ...*model.RoleResource) error {
	rr := r.baseRepo.RoleResource
	if err := rr.WithContext(ctx).Create(resource...); err != nil {
		err = errors.Wrap(err, "role_resource.repo.Insert")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) Update(ctx context.Context, resource *model.RoleResource) error {
	rr := r.baseRepo.RoleResource
	if _, err := rr.WithContext(ctx).Where(rr.ID.Eq(resource.ID)).Updates(resource); err != nil {
		err = errors.Wrap(err, "role_resource.repo.Update")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) Delete(ctx context.Context, id uint) error {
	rr := r.baseRepo.RoleResource
	if _, err := rr.WithContext(ctx).Where(rr.ID.Eq(id)).Delete(); err != nil {
		err = errors.Wrap(err, "role_resource.repo.Delete")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) DeleteByRoleIDs(ctx context.Context, rids []uint64) error {
	rr := r.baseRepo.RoleResource
	if _, err := rr.WithContext(ctx).Where(rr.RoleID.In(rids...)).Delete(); err != nil {
		err = errors.Wrap(err, "role_resource.repo.DeleteByRoleIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) DeleteByResourceIDs(ctx context.Context, resIDs []uint64, resourceType model.ResourceType) error {
	rr := r.baseRepo.RoleResource
	if _, err := rr.WithContext(ctx).Where(rr.ResourceID.In(resIDs...), rr.Type.Eq(uint8(resourceType))).Delete(); err != nil {
		err = errors.Wrap(err, "role_resource.repo.DeleteByResourceIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) UpdateByRoleID(ctx context.Context, rid uint64, rrs []*model.RoleResource) error {
	return r.dao.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.RoleResource{}).Where("role_id = ?", rid).Delete(&model.RoleResource{})
		return tx.Create(&rrs).Error
	})
}
