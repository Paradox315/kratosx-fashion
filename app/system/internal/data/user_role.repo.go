package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/app/system/internal/data/query"
)

type userRoleRepo struct {
	dao      *Data
	log      *log.Helper
	baseRepo *query.Query
}

func NewUserRoleRepo(data *Data, logger log.Logger) biz.UserRoleRepo {
	return &userRoleRepo{
		dao:      data,
		log:      log.NewHelper(logger),
		baseRepo: query.Use(data.DB),
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

func (u *userRoleRepo) SelectRoleIDByUserID(ctx context.Context, uid uint64) (rids []uint64, err error) {
	ur := u.baseRepo.UserRole
	err = ur.WithContext(ctx).Where(ur.UserID.Eq(uid)).Scan(&rids)
	return
}

func (u *userRoleRepo) Insert(ctx context.Context, userRole *model.UserRole) error {
	ur := u.baseRepo.UserRole
	return ur.WithContext(ctx).Create(userRole)
}

func (u *userRoleRepo) Update(ctx context.Context, userRole *model.UserRole) error {
	ur := u.baseRepo.UserRole
	return ur.WithContext(ctx).Save(userRole)
}

func (u *userRoleRepo) Delete(ctx context.Context, id uint) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).Delete()
	return err
}

func (u *userRoleRepo) DeleteByUserID(ctx context.Context, uid uint64) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.UserID.Eq(uid)).Delete()
	return err
}

func (u *userRoleRepo) DeleteByUserIDs(ctx context.Context, uids []uint64) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.UserID.In(uids...)).Delete()
	return err
}

func (u *userRoleRepo) DeleteByRoleID(ctx context.Context, rid uint64) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.RoleID.Eq(rid)).Delete()
	return err
}

func (u *userRoleRepo) DeleteByRoleIDs(ctx context.Context, rids []uint64) error {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.RoleID.In(rids...)).Delete()
	return err
}

func (u *userRoleRepo) ExistByUserIDAndRoleID(ctx context.Context, uid uint64, rid uint64) (bool, error) {
	ur := u.baseRepo.UserRole
	_, err := ur.WithContext(ctx).Where(ur.UserID.Eq(uid)).Where(ur.RoleID.Eq(rid)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
}
