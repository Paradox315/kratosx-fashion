package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
)

type userRoleRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewUserRoleRepo(data *data.Data, logger log.Logger) biz.UserRoleRepo {
	return &userRoleRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "user_role")),
		baseRepo: linq.Use(data.DB),
	}
}

func (u *userRoleRepo) Select(ctx context.Context, id uint) (*model.UserRole, error) {
	ur := u.baseRepo.UserRole
	return ur.WithContext(ctx).Where(ur.ID.Eq(id)).First()
}

func (u *userRoleRepo) SelectAll(ctx context.Context) ([]*model.UserRole, error) {
	ur := u.baseRepo.UserRole
	return ur.WithContext(ctx).Find()
}

func (u *userRoleRepo) SelectAllByUserID(ctx context.Context, uid uint64) ([]*model.UserRole, error) {
	ur := u.baseRepo.UserRole
	return ur.WithContext(ctx).Where(ur.UserID.Eq(uid)).Find()
}

func (u *userRoleRepo) Insert(ctx context.Context, userRole ...*model.UserRole) error {
	ur := u.baseRepo.UserRole
	return ur.WithContext(ctx).Create(userRole...)
}

func (u *userRoleRepo) Update(ctx context.Context, userRole *model.UserRole) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Updates(userRole)
	return err
}

func (u *userRoleRepo) UpdateByUserID(ctx context.Context, uid uint64, urs []*model.UserRole) error {
	return u.dao.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ?", uid).Delete(&model.UserRole{})
		return tx.Create(&urs).Error
	})
}
func (u *userRoleRepo) Delete(ctx context.Context, id uint) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).Delete()
	return err
}

func (u *userRoleRepo) DeleteByUserIDs(ctx context.Context, uids []uint64) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.UserID.In(uids...)).Delete()
	return err
}

func (u *userRoleRepo) DeleteByRoleIDs(ctx context.Context, rids []uint64) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.RoleID.In(rids...)).Delete()
	return err
}
