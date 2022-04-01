package repo

import (
	"context"
	"net/http"
	"strings"

	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/conf"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/model"

	"github.com/casbin/casbin/v2"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/xhttp/apistate"
	"github.com/gofiber/fiber/v2"

	kerrors "github.com/go-kratos/kratos/v2/errors"
)

type ResourceRouterRepo struct {
	dao   *data.Data
	cRepo *casbin.SyncedEnforcer
	cli   *fiber.Agent
	log   *log.Helper
}

func NewResourceRouterRepo(dao *data.Data, c *conf.Server, logger log.Logger, casbinRepo *casbin.SyncedEnforcer) biz.ResourceRouterRepo {
	return &ResourceRouterRepo{
		dao:   dao,
		cRepo: casbinRepo,
		log:   log.NewHelper(log.With(logger, "repo", "resource_router")),
		cli:   fiber.Get("http://" + c.Http.Addr),
	}
}

func (r *ResourceRouterRepo) parseGroup(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.Split(name, "-")[0]
}

func (r *ResourceRouterRepo) SelectAll(ctx context.Context) (rs []biz.Router, err error) {
	_, body, errs := r.cli.Bytes()
	if errs != nil {
		r.log.Error("获取资源路由失败", errs)
		return nil, kerrors.InternalServer("HTTP_CLIENT", "资源路由获取失败")
	}
	var resp apistate.Resp[[][]*fiber.Route]
	if err = encoding.GetCodec("json").Unmarshal(body, &resp); err != nil {
		r.log.Error("序列化失败", err)
		return nil, kerrors.InternalServer("CODEC", "序列化失败")
	}
	if resp.Code != http.StatusOK {
		r.log.Error("获取资源路由失败", resp.Message)
		return nil, kerrors.InternalServer("HTTP_CLIENT", "资源路由获取失败")
	}
	routers := resp.Metadata
	for _, router := range routers {
		for _, ro := range router {
			switch ro.Method {
			case http.MethodHead, http.MethodOptions, http.MethodTrace, http.MethodConnect, http.MethodPatch:
				continue
			}
			if ro.Name == "" {
				continue
			}
			rs = append(rs, biz.Router{
				Method: ro.Method,
				Path:   ro.Path,
				Name:   ro.Name,
				Params: ro.Params,
				Group:  r.parseGroup(ro.Name),
			})
		}
	}
	return
}

func (r *ResourceRouterRepo) SelectByRoleIDs(ctx context.Context, rids []string) (rrs []*model.ResourceRouter, err error) {
	list := r.cRepo.GetFilteredPolicy(0, rids...)
	for _, rr := range list {
		rrs = append(rrs, &model.ResourceRouter{
			Method: rr[1],
			Path:   rr[2],
		})
	}
	return
}

func (r *ResourceRouterRepo) Update(ctx context.Context, router []*model.ResourceRouter) (err error) {
	rid := router[0].RoleID
	if err = r.ClearByRoleIDs(ctx, []string{rid}); err != nil {
		err = kerrors.InternalServer("CASBIN", "清除角色资源路由失败")
		return
	}
	var rules [][]string
	for _, v := range router {
		rules = append(rules, []string{v.RoleID, v.Path, v.Method})
	}
	if ok, err := r.cRepo.AddPolicy(rules); ok {
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
