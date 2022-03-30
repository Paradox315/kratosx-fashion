package data

import (
	"context"
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-kratos/kratos/v2/log"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type ResourceRouterRepo struct {
	dao        *Data
	casbinRepo *casbin.SyncedEnforcer
	log        *log.Helper
}

func NewResourceRouterRepo(dao *Data, logger log.Logger) biz.ResourceRouterRepo {
	r := &ResourceRouterRepo{
		dao: dao,
		log: log.NewHelper(log.With(logger, "repo", "resource_router")),
	}
	adapter, err := gormAdapter.NewAdapterByDB(dao.DB)
	if err != nil {
		panic(err)
	}
	r.casbinRepo, err = casbin.NewSyncedEnforcer("./resource/model.conf", adapter)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *ResourceRouterRepo) SelectAll(ctx context.Context) ([]*model.ResourceRouter, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceRouterRepo) SelectByRoleIDs(ctx context.Context, strings []string) ([]*model.ResourceRouter, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceRouterRepo) Insert(ctx context.Context, router *model.ResourceRouter) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceRouterRepo) Update(ctx context.Context, router *model.ResourceRouter) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceRouterRepo) DeleteByRoleIDs(ctx context.Context, strings []string) error {
	//TODO implement me
	panic("implement me")
}

func (r *ResourceRouterRepo) Exist(ctx context.Context, roleID string, path string, method string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
