package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"kratosx-fashion/pkg/ctxutil"
	"net/http"
	"strings"
	"time"

	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/model"

	"github.com/casbin/casbin/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gofiber/fiber/v2"

	kerrors "github.com/go-kratos/kratos/v2/errors"
)

const routerAll = "router:all"

type ResourceRouterRepo struct {
	dao   *data.Data
	cRepo *casbin.SyncedEnforcer
	log   *log.Helper
}

func NewResourceRouterRepo(dao *data.Data, logger log.Logger, casbinRepo *casbin.SyncedEnforcer) biz.ResourceRouterRepo {
	return &ResourceRouterRepo{
		dao:   dao,
		cRepo: casbinRepo,
		log:   log.NewHelper(log.With(logger, "repo", "resource_router")),
	}
}

func parseGroup(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.Split(name, "-")[0]
}

func (r *ResourceRouterRepo) SelectAll(ctx context.Context) (rs []model.Router, err error) {
	bytes, err := r.dao.RDB.WithContext(ctx).Get(ctx, routerAll).Bytes()
	if err == nil {
		_ = codec.Unmarshal(bytes, &rs)
		return
	}
	if err != redis.Nil {
		err = errors.Wrap(err, "redis.Get")
		r.log.WithContext(ctx).Error(err)
		return
	}
	var routers [][]*fiber.Route
	if fiberCtx, ok := ctxutil.GetFiberCtx(ctx); !ok {
		return nil, errors.New("get fiber context failed")
	} else {
		routers = fiberCtx.App().Stack()
	}
	for _, router := range routers {
		for _, ro := range router {
			if ro.Name == "" {
				continue
			}
			switch ro.Method {
			case http.MethodHead, http.MethodOptions, http.MethodTrace, http.MethodConnect, http.MethodPatch:
				continue
			}
			rs = append(rs, model.Router{
				Method: "(" + ro.Method + ")",
				Path:   ro.Path,
				Name:   ro.Name,
				Params: ro.Params,
				Group:  parseGroup(ro.Name),
			})
		}
	}
	bytes, _ = codec.Marshal(&rs)
	err = r.dao.RDB.Set(ctx, routerAll, bytes, time.Hour*1).Err()
	return
}

func (r *ResourceRouterRepo) SelectByRoleIDs(ctx context.Context, rids []string) (rrs []model.ResourceRouter, err error) {
	list := r.cRepo.GetFilteredPolicy(0, rids...)
	for _, rr := range list {
		rrs = append(rrs, model.ResourceRouter{
			RoleID: rr[0],
			Path:   rr[1],
			Method: rr[2],
		})
	}
	return
}

func (r *ResourceRouterRepo) Update(ctx context.Context, router []model.ResourceRouter) (err error) {
	rid := router[0].RoleID
	if err = r.ClearByRoleIDs(ctx, []string{rid}); err != nil {
		err = kerrors.InternalServer("CASBIN", "清除角色资源路由失败")
		return
	}
	var rules [][]string
	for _, v := range router {
		rules = append(rules, []string{v.RoleID, v.Path, v.Method})
	}
	if ok, err := r.cRepo.AddPolicies(rules); ok {
		return nil
	} else {
		return err
	}
}

func (r *ResourceRouterRepo) ClearByRoleIDs(ctx context.Context, rids []string) error {
	if ok, err := r.cRepo.RemoveFilteredPolicy(0, rids...); ok {
		return nil
	} else {
		return err
	}
}
