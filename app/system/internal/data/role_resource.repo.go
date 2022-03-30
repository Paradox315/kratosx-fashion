package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
)

type RoleResourceRepo struct {
	dao      *Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewRoleResourceRepo(data *Data, logger log.Logger) biz.RoleResourceRepo {
	return &RoleResourceRepo{
		dao:      data,
		log:      log.NewHelper(logger),
		baseRepo: linq.Use(data.DB),
	}
}

func (r *RoleResourceRepo) Select(ctx context.Context, id uint) (*model.RoleResource, error) {
	rr := r.baseRepo.RoleResource
	return rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
}

func (r *RoleResourceRepo) SelectByRoleID(ctx context.Context, rid uint64, resourceType ...model.ResourceType) ([]*model.RoleResource, error) {
	rr := r.baseRepo.RoleResource
	if len(resourceType) == 0 {
		return rr.WithContext(ctx).Where(rr.RoleID.Eq(rid)).Find()
	}
	typ := resourceType[0]
	return rr.WithContext(ctx).Where(rr.RoleID.Eq(rid), rr.Type.Eq(uint8(typ))).Find()

}

func (r *RoleResourceRepo) Insert(ctx context.Context, resource ...*model.RoleResource) error {
	rr := r.baseRepo.RoleResource
	return rr.WithContext(ctx).Create(resource...)
}

func (r *RoleResourceRepo) Update(ctx context.Context, resource *model.RoleResource) error {
	rr := r.baseRepo.RoleResource
	_, err := rr.WithContext(ctx).Updates(resource)
	return err
}

func (r *RoleResourceRepo) Delete(ctx context.Context, id uint) error {
	rr := r.baseRepo.RoleResource
	_, err := rr.WithContext(ctx).Where(rr.ID.Eq(id)).Delete()
	return err
}

func (r *RoleResourceRepo) DeleteByRoleIDs(ctx context.Context, rids []uint64) error {
	rr := r.baseRepo.RoleResource
	_, err := rr.WithContext(ctx).Where(rr.RoleID.In(rids...)).Delete()
	return err
}

func (r *RoleResourceRepo) DeleteByResourceIDs(ctx context.Context, resIDs []uint64, resourceType model.ResourceType) error {
	rr := r.baseRepo.RoleResource
	_, err := rr.WithContext(ctx).Where(rr.ResourceID.In(resIDs...), rr.Type.Eq(uint8(resourceType))).Delete()
	return err
}

func (r *RoleResourceRepo) BaseRepo(ctx context.Context) *linq.Query {
	return r.baseRepo
}

func (r *RoleResourceRepo) UpdateByRoleID(ctx context.Context, rid uint64, rrs []*model.RoleResource) error {
	return r.dao.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("role_id = ?", rid).Delete(&model.RoleResource{})
		return tx.Create(&rrs).Error
	})
}
