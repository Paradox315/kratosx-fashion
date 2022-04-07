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

func (u *userRepo) Select(ctx context.Context, id uint) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.Select")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectByUsername(ctx context.Context, username string) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Username.Eq(username)).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectByUsername")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectPasswordByUID(ctx context.Context, uid uint) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.ID.Eq(uid)).Select(ur.Password).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPasswordByUID")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectPasswordByName(ctx context.Context, username string) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Username.Eq(username)).Select(ur.Password).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPasswordByName")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectPasswordByMobile(ctx context.Context, mobile string) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Mobile.Eq(mobile)).Select(ur.Password).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPasswordByMobile")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectPasswordByEmail(ctx context.Context, email string) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Email.Eq(email)).Select(ur.Password).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPasswordByEmail")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectPage(ctx context.Context, limit, offset int, opt *biz.SQLOption) (users []*model.User, total int64, err error) {
	tx := u.dao.DB.Model(&model.User{})
	if err = tx.Count(&total).Error; err != nil {
		err = errors.Wrap(err, "userRepo.SelectPage")
		u.log.WithContext(ctx).Error(err)
		return
	}
	if opt != nil && len(opt.Where) > 0 {
		tx = tx.Where(opt.Where, opt.Args...)
	}
	if opt != nil && len(opt.Order) > 0 {
		tx = tx.Order(opt.Order)
	}
	err = tx.WithContext(ctx).Where("status = ?", model.UserStatusNormal).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPage")
		u.log.WithContext(ctx).Error(err)
		return
	}
	return
}

func (u *userRepo) Insert(ctx context.Context, user *model.User) error {
	ur := u.baseRepo.User
	if err := ur.WithContext(ctx).Create(user); err != nil {
		err = errors.Wrap(err, "userRepo.Insert")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRepo) Update(ctx context.Context, user *model.User) error {
	ur := u.baseRepo.User
	if _, err := ur.WithContext(ctx).Where(ur.ID.Eq(user.ID)).Updates(user); err != nil {
		err = errors.Wrap(err, "userRepo.Update")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRepo) UpdateStatus(ctx context.Context, id uint, status model.UserStatus) error {
	ur := u.baseRepo.User
	if _, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).Update(ur.Status, status); err != nil {
		err = errors.Wrap(err, "userRepo.UpdateStatus")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRepo) DeleteByIDs(ctx context.Context, ids []uint) error {
	ur := u.baseRepo.User
	if _, err := ur.WithContext(ctx).Where(ur.ID.In(ids...)).Delete(); err != nil {
		err = errors.Wrap(err, "userRepo.DeleteByIDs")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (u *userRepo) ExistByUsername(ctx context.Context, username string) bool {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.Username.Eq(username)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		u.log.WithContext(ctx).Error("userRepo.ExistByUsername error", err)
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
