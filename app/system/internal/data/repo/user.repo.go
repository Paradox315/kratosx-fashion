package repo

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
)

type userRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
}

func NewUserRepo(data *data.Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "user")),
		baseRepo: linq.Use(data.DB),
	}
}

func (u userRepo) Select(ctx context.Context, id uint) (*model.User, error) {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Where(ur.ID.Eq(id)).First()
}

func (u userRepo) SelectByUsername(ctx context.Context, username string) (*model.User, error) {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Where(ur.Username.Eq(username)).First()
}

func (u *userRepo) SelectPasswordByName(ctx context.Context, username string) (*model.User, error) {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Where(ur.Username.Eq(username)).Select(ur.Password, ur.ID, ur.Status).First()
}

func (u *userRepo) SelectPasswordByMobile(ctx context.Context, mobile string) (*model.User, error) {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Where(ur.Mobile.Eq(mobile)).Select(ur.Password, ur.ID, ur.Status).First()
}

func (u *userRepo) SelectPasswordByEmail(ctx context.Context, email string) (*model.User, error) {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Where(ur.Email.Eq(email)).Select(ur.Password, ur.ID, ur.Status).First()
}

func (u userRepo) List(ctx context.Context, limit, offset int, sqlopt biz.SQLOption) (users []model.User, total int64, err error) {
	err = u.dao.DB.Where(sqlopt.Where, sqlopt.Args).Order(sqlopt.Order).Offset(offset).Limit(limit).Count(&total).Find(&users).Error
	return
}

func (u userRepo) Insert(ctx context.Context, user *model.User) error {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Create(user)
}

func (u userRepo) Update(ctx context.Context, user *model.User) error {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Updates(user)
	return err
}

func (u *userRepo) UpdateStatus(ctx context.Context, id uint, status model.UserStatus) error {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).Update(ur.Status, status)
	return err
}

func (u userRepo) DeleteByIDs(ctx context.Context, ids []uint) error {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.ID.In(ids...)).Delete()
	return err
}

func (u *userRepo) ExistByUserName(ctx context.Context, username string) bool {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.Username.Eq(username)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		u.log.WithContext(ctx).Error("userRepo.ExistByUserName error", err)
		return false
	}
	return true
}

func (u *userRepo) ExistByEmail(ctx context.Context, email string) bool {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.Email.Eq(email)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		u.log.WithContext(ctx).Error("userRepo.ExistByEmail error", err)
		return false
	}

	return true
}

func (u *userRepo) ExistByMobile(ctx context.Context, mobile string) bool {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.Mobile.Eq(mobile)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		u.log.WithContext(ctx).Error("userRepo.ExistByMobile error", err)
		return false
	}
	return true
}

func (u *userRepo) BaseRepo(ctx context.Context) *linq.Query {
	return u.baseRepo
}
