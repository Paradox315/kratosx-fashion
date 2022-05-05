package repo

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
	"marwan.io/singleflight"
	"time"
)

const (
	userKey = "user:%d"
)

type userRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
	sf       *singleflight.Group[*model.User]
}

func NewUserRepo(data *data.Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "user")),
		baseRepo: linq.Use(data.DB),
		sf:       &singleflight.Group[*model.User]{},
	}
}

func (u *userRepo) Select(ctx context.Context, id uint) (*model.User, error) {
	result, err, _ := u.sf.Do(fmt.Sprintf(userKey, id), func() (*model.User, error) {
		bytes, err := u.dao.RDB.Get(ctx, fmt.Sprintf(userKey, id)).Bytes()
		if err == nil {
			var user *model.User
			_ = codec.Unmarshal(bytes, &user)
			return user, nil
		}
		ur := u.baseRepo.User
		user, err := ur.WithContext(ctx).Where(ur.ID.Eq(id)).First()
		if err != nil {
			err = errors.Wrap(err, "userRepo.Select")
			u.log.WithContext(ctx).Error(err)
			return nil, err
		}
		bytes, _ = codec.Marshal(user)
		if err = u.dao.RDB.Set(ctx, fmt.Sprintf(userKey, id), bytes, time.Hour*1).Err(); err != nil {
			err = errors.Wrap(err, "userRepo.Select.redis.Set")
			u.log.WithContext(ctx).Error(err)
			return nil, err
		}
		return user, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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
	user, err := ur.WithContext(ctx).Where(ur.ID.Eq(uid)).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPasswordByUID")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectByMobile(ctx context.Context, mobile string) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Mobile.Eq(mobile)).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectByMobile")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepo) SelectByEmail(ctx context.Context, email string) (*model.User, error) {
	ur := u.baseRepo.User
	user, err := ur.WithContext(ctx).Where(ur.Email.Eq(email)).First()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectByEmail")
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
	err = tx.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectPage")
		u.log.WithContext(ctx).Error(err)
		return
	}
	return
}

func (u *userRepo) SelectTokens(ctx context.Context, id uint) ([]string, error) {
	result, err := u.dao.RDB.MGet(ctx, fmt.Sprintf(userRefreshTokenKey, id), fmt.Sprintf(userAccessTokenKey, id)).Result()
	if err != nil {
		err = errors.Wrap(err, "userRepo.SelectTokens")
		u.log.WithContext(ctx).Error(err)
		return nil, err
	}
	var tokens []string
	for _, v := range result {
		if v != nil {
			tokens = append(tokens, v.(string))
		}
	}
	return tokens, nil
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
	return u.dao.RDB.Del(ctx, fmt.Sprintf(userKey, user.ID)).Err()
}

func (u *userRepo) DeleteByIDs(ctx context.Context, ids []uint) error {
	ur := u.baseRepo.User
	if _, err := ur.WithContext(ctx).Where(ur.ID.In(ids...)).Delete(); err != nil {
		err = errors.Wrap(err, "userRepo.DeleteByIDs")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	return u.dao.RDB.Del(ctx, lo.Map(ids, func(id uint, _ int) string {
		return fmt.Sprintf(userKey, id)
	})...).Err()
}

func (u *userRepo) Verify(ctx context.Context, id uint) bool {
	user, err := u.Select(ctx, id)
	if err != nil {
		return false
	}
	return user.Status == model.UserStatusNormal
}
