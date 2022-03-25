package data

import (
	"context"
	"errors"

	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/app/system/internal/data/query"
	"kratosx-fashion/pkg/option"
	"kratosx-fashion/pkg/pagination"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	pb "kratosx-fashion/api/system/v1"
)

type userRepo struct {
	dao      *Data
	log      *log.Helper
	baseRepo *query.Query
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		dao:      data,
		log:      log.NewHelper(logger),
		baseRepo: query.Use(data.DB),
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

func (u userRepo) SelectPasswordByName(ctx context.Context, username string) (uid uint, pwd string, err error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Username.Eq(username)).Select(ur.Password, ur.ID).First()
	if err != nil {
		return
	}
	uid = user.ID
	pwd = user.Password
	return
}

func (u userRepo) List(ctx context.Context, req *pb.ListRequest, opts ...*pb.QueryOption) (users []*model.User, total int64, err error) {
	orders, keywords := option.Parse(opts...)
	limit, offset, err := pagination.Parse(req)
	if err != nil {
		u.log.WithContext(ctx).Error("pagination.Parse error", err)
		return
	}
	tx := u.dao.DB.Where(keywords).Limit(limit).Offset(offset)
	for _, order := range orders {
		tx.Order(order)
	}
	err = tx.Count(&total).Find(users).Error
	return
}

func (u userRepo) Insert(ctx context.Context, user *model.User) error {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Create(user)
}

func (u userRepo) Update(ctx context.Context, user *model.User) error {
	ur := u.baseRepo.User
	return ur.WithContext(ctx).Save(user)
}

func (u *userRepo) UpdateStatus(ctx context.Context, id uint, status uint8) error {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Update(ur.Status, status)
	return err
}

func (u userRepo) Delete(ctx context.Context, id uint) error {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).Delete()
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
	}
	return true
}

func (u *userRepo) ExistByEmail(ctx context.Context, email string) bool {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.Email.Eq(email)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (u *userRepo) ExistByMobile(ctx context.Context, mobile string) bool {
	ur := u.baseRepo.User
	_, err := ur.WithContext(ctx).Where(ur.Mobile.Eq(mobile)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
