package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/biz"
	"kratosx-fashion/app/system/data"
	"kratosx-fashion/app/system/data/linq"
	"kratosx-fashion/app/system/data/model"
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

func (r *RoleResourceRepo) SelectByRoleID(ctx context.Context, rid uint, resourceType ...model.ResourceType) (list []*model.RoleResource, err error) {
	rr := r.baseRepo.RoleResource
	if len(resourceType) == 0 {
		list, err = rr.WithContext(ctx).Where(rr.RoleID.Eq(rid)).Select(rr.ID, rr.RoleID, rr.ResourceID, rr.Type).Find()
	} else {
		typ := resourceType[0]
		list, err = rr.WithContext(ctx).Where(rr.RoleID.Eq(rid), rr.Type.Eq(uint8(typ))).Select(rr.ID, rr.RoleID, rr.ResourceID, rr.Type).Find()
	}
	if err != nil {
		err = errors.Wrap(err, "role_resource.repo.SelectByRoleID")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	var result []*model.RoleResource
	visited := make(map[string]struct{})
	for _, e := range list {
		if _, ok := visited[e.ResourceID]; ok {
			continue
		}
		result = append(result, e)
		visited[e.ResourceID] = struct{}{}
	}
	return result, nil
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

func (r *RoleResourceRepo) DeleteByRoleIDs(ctx context.Context, rids []uint) error {
	rr := r.baseRepo.RoleResource
	if _, err := rr.WithContext(ctx).Where(rr.RoleID.In(rids...)).Delete(); err != nil {
		err = errors.Wrap(err, "role_resource.repo.DeleteByRoleIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) DeleteByResourceIDs(ctx context.Context, resIDs []string, resourceType model.ResourceType) error {
	rr := r.baseRepo.RoleResource
	if _, err := rr.WithContext(ctx).Where(rr.ResourceID.In(resIDs...), rr.Type.Eq(uint8(resourceType))).Delete(); err != nil {
		err = errors.Wrap(err, "role_resource.repo.DeleteByResourceIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *RoleResourceRepo) UpdateByRoleID(ctx context.Context, rid uint, rrs []*model.RoleResource) error {
	return r.dao.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.RoleResource{}).Where("role_id = ?", rid).Delete(&model.RoleResource{})
		return tx.Create(&rrs).Error
	})
}

func (r *RoleResourceRepo) SelectByResourceID(ctx context.Context, rid string, resourceType ...model.ResourceType) ([]*model.RoleResource, error) {
	rr := r.baseRepo.RoleResource
	if len(resourceType) == 0 {
		return rr.WithContext(ctx).Where(rr.ResourceID.Eq(rid)).Find()
	}
	typ := resourceType[0]
	list, err := rr.WithContext(ctx).Where(rr.ResourceID.Eq(rid), rr.Type.Eq(uint8(typ))).Select(rr.ID, rr.RoleID, rr.ResourceID, rr.Type).Find()
	if err != nil {
		err = errors.Wrap(err, "role_resource.repo.SelectByResourceID")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return list, nil
}
