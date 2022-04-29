package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"
	"marwan.io/singleflight"
	"time"
)

const (
	menuAll = "menu:all"
)

var codec = encoding.GetCodec("json")

type ResourceMenuRepo struct {
	dao      *data.Data
	log      *log.Helper
	baseRepo *linq.Query
	sf       *singleflight.Group[[]*model.ResourceMenu]
}

func NewResourceMenuRepo(dao *data.Data, logger log.Logger) biz.ResourceMenuRepo {
	return &ResourceMenuRepo{
		dao:      dao,
		log:      log.NewHelper(log.With(logger, "repo", "resource_menu")),
		baseRepo: linq.Use(dao.DB),
		sf:       &singleflight.Group[[]*model.ResourceMenu]{},
	}
}

func (r *ResourceMenuRepo) Select(ctx context.Context, id uint) (menu *model.ResourceMenu, err error) {
	rr := r.baseRepo.ResourceMenu
	menu, err = rr.WithContext(ctx).Where(rr.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "resource_menu.repo.Select")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return
}

func (r *ResourceMenuRepo) SelectByIDs(ctx context.Context, ids []uint) ([]*model.ResourceMenu, error) {
	rr := r.baseRepo.ResourceMenu
	menus, err := rr.WithContext(ctx).Where(rr.ID.In(ids...)).Find()
	if err != nil {
		err = errors.Wrap(err, "resource_menu.repo.SelectByIDs")
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return menus, nil
}

func (r *ResourceMenuRepo) SelectAll(ctx context.Context) ([]*model.ResourceMenu, error) {
	result, err, _ := r.sf.Do("resource_menu:all", func() (menus []*model.ResourceMenu, err error) {
		bytes, err := r.dao.RDB.WithContext(ctx).Get(ctx, menuAll).Bytes()
		if err == nil {
			_ = codec.Unmarshal(bytes, &menus)
			return
		}
		rr := r.baseRepo.ResourceMenu
		menus, err = rr.WithContext(ctx).Find()
		if err != nil {
			err = errors.Wrap(err, "resource_menu.repo.SelectAll")
			r.log.WithContext(ctx).Error(err)
			return nil, err
		}
		bytes, _ = codec.Marshal(menus)
		err = r.dao.RDB.Set(ctx, menuAll, bytes, time.Hour*1).Err()
		return
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ResourceMenuRepo) Insert(ctx context.Context, menu ...*model.ResourceMenu) error {
	rr := r.baseRepo.ResourceMenu
	if err := rr.WithContext(ctx).Create(menu...); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.Insert")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	if err := r.dao.RDB.Del(ctx, menuAll).Err(); err != nil {
		err = errors.Wrap(err, "redis.Del")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *ResourceMenuRepo) Update(ctx context.Context, menu *model.ResourceMenu) error {
	rr := r.baseRepo.ResourceMenu
	if _, err := rr.WithContext(ctx).Where(rr.ID.Eq(menu.ID)).Updates(menu); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.Update")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	if err := r.dao.RDB.Del(ctx, menuAll).Err(); err != nil {
		err = errors.Wrap(err, "redis.Del")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return r.dao.RDB.WithContext(ctx).Del(ctx, menuAll).Err()
}

func (r *ResourceMenuRepo) DeleteByIDs(ctx context.Context, ids []uint) error {
	rr := r.baseRepo.ResourceMenu
	if _, err := rr.WithContext(ctx).Where(rr.ID.In(ids...)).Delete(); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.DeleteByIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	if _, err := rr.WithContext(ctx).Where(rr.ParentID.In(ids...)).Delete(); err != nil {
		err = errors.Wrap(err, "resource_menu.repo.DeleteByIDs")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	if err := r.dao.RDB.Del(ctx, menuAll).Err(); err != nil {
		err = errors.Wrap(err, "redis.Del")
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return r.dao.RDB.WithContext(ctx).Del(ctx, menuAll).Err()
}
