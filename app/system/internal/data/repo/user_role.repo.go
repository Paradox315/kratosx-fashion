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
	usr, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "user_role.Select")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return usr, nil
}

func (u *userRoleRepo) SelectAll(ctx context.Context) ([]*model.UserRole, error) {
	ur := u.baseRepo.UserRole
	urs, err := ur.WithContext(ctx).Find()
	if err != nil {
		err = errors.Wrap(err, "user_role.SelectAll")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return urs, nil
}

func (u *userRoleRepo) SelectAllByUserID(ctx context.Context, uid uint) ([]*model.UserRole, error) {
	ur := u.baseRepo.UserRole
	list, err := ur.WithContext(ctx).Where(ur.UserID.Eq(uid)).Find()
	if err != nil {
		err = errors.Wrap(err, "user_role.SelectAllByUserID")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return list, nil
}

func (u *userRoleRepo) Insert(ctx context.Context, userRole ...*model.UserRole) error {
	ur := u.baseRepo.UserRole
	if err := ur.WithContext(ctx).Create(userRole...); err != nil {
		err = errors.Wrap(err, "user_role.Insert")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRoleRepo) Update(ctx context.Context, userRole *model.UserRole) error {
	ur := u.baseRepo.UserRole
	if _, err := ur.WithContext(ctx).Where(ur.ID.Eq(userRole.ID)).Updates(userRole); err != nil {
		err = errors.Wrap(err, "user_role.Update")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRoleRepo) UpdateByUserID(ctx context.Context, uid uint, urs []*model.UserRole) error {
	return u.dao.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.UserRole{}).Where("user_id = ?", uid).Delete(&model.UserRole{})
		return tx.Create(&urs).Error
	})
}

func (u *userRoleRepo) DeleteByUserIDs(ctx context.Context, uids []uint) error {
	ur := u.baseRepo.UserRole
	if _, err := ur.WithContext(ctx).Where(ur.UserID.In(uids...)).Delete(); err != nil {
		err = errors.Wrap(err, "user_role.DeleteByUserIDs")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRoleRepo) DeleteByRoleIDs(ctx context.Context, rids []uint) error {
	ur := u.baseRepo.UserRole
	if _, err := ur.WithContext(ctx).Where(ur.RoleID.In(rids...)).Delete(); err != nil {
		err = errors.Wrap(err, "user_role.DeleteByRoleIDs")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}
