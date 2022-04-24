package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/math"
	"strconv"
	"strings"
)

var codec = encoding.GetCodec("json")

type ResourceUsecase struct {
	roleResourceRepo RoleResourceRepo
	menuRepo         ResourceMenuRepo
	routeRepo        ResourceRouterRepo
	tx               Transaction
	log              *log.Helper
}

func NewResourceUsecase(menuRepo ResourceMenuRepo, routeRepo ResourceRouterRepo, roleResourceRepo RoleResourceRepo, tx Transaction, logger log.Logger) *ResourceUsecase {
	return &ResourceUsecase{
		menuRepo:         menuRepo,
		routeRepo:        routeRepo,
		roleResourceRepo: roleResourceRepo,
		tx:               tx,
		log:              log.NewHelper(log.With(logger, "biz", "resource")),
	}
}

func longestCommonPrefix(rs []model.Router, group string) (prefix string) {
	if len(rs) == 0 {
		return ""
	}
	group = strings.ToLower(group)
	index := strings.Index(rs[0].Path, group)
	if index != -1 {
		return rs[0].Path[:index+len(group)] + "/*"
	}
	for i := 0; i < len(rs[0].Path); i++ {
		for _, r := range rs {
			if i >= len(r.Path) || r.Path[i] != rs[0].Path[i] {
				return rs[0].Path[:i] + "*"
			}
		}
	}
	return "*"
}

func (r *ResourceUsecase) buildRouters(ctx context.Context, rpos []model.Router) (routers []RouterGroup) {
	routerMap := make(map[string][]model.Router)
	for _, rpo := range rpos {
		routerMap[rpo.Group] = append(routerMap[rpo.Group], rpo)
	}
	for group, children := range routerMap {
		routers = append(routers, RouterGroup{
			Path:     longestCommonPrefix(children, group),
			Name:     group,
			Children: children,
			Method: func(routers []model.Router) string {
				var methods []string
				visited := make(map[string]bool)
				for _, route := range routers {
					if _, ok := visited[route.Method]; ok {
						continue
					}
					methods = append(methods, route.Method)
					visited[route.Method] = true
				}
				return strings.Join(methods, "|")
			}(children),
		})
	}
	return
}
func (r *ResourceUsecase) buildMenuPO(ctx context.Context, menu *pb.MenuRequest) (mpo *model.ResourceMenu, err error) {
	mpo = &model.ResourceMenu{
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
	}
	mpo.ID = cast.ToUint(menu.Id)
	mpo.ParentID = cast.ToUint(menu.ParentId)
	mpo.Locale = menu.Meta.Locale
	mpo.Icon = menu.Meta.Icon
	mpo.Order = menu.Meta.Order
	if menu.Meta.HideInMenu {
		mpo.HideInMenu = 1
	}
	if menu.Meta.IgnoreCache {
		mpo.IgnoreCache = 1
	}
	if menu.Meta.RequireAuth {
		mpo.RequireAuth = 1
	}
	if menu.Meta.NoAffix {
		mpo.NoAffix = 1
	}
	if len(menu.Actions) == 0 {
		return
	}
	var acts []*model.ResourceAction
	for _, act := range menu.Actions {
		acts = append(acts, &model.ResourceAction{
			Name: act.Name,
			Code: act.Code,
		})
	}
	bytes, err := codec.Marshal(acts)
	if err != nil {
		return
	}
	mpo.Actions = string(bytes)
	return
}

func (r *ResourceUsecase) buildTree(ctx context.Context, mpos []*model.ResourceMenu) (menus []Menu, err error) {
	for _, mpo := range mpos {
		var menu Menu
		menu, err = r.buildMenuDO(ctx, mpo)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	menuMap := lo.GroupBy(menus, func(menu Menu) string {
		return cast.ToString(menu.ParentId)
	})
	menus = menuMap["0"]
	for i := 0; i < len(menus); i++ {
		r.buildMenuChild(ctx, &menus[i], menuMap)
	}
	return
}

func (r *ResourceUsecase) buildMenuChild(ctx context.Context, menu *Menu, menuMap map[string][]Menu) {
	menu.Children = menuMap[menu.Id]
	for _, child := range menu.Children {
		r.buildMenuChild(ctx, &child, menuMap)
	}
	return
}

func (r *ResourceUsecase) buildMenuDO(ctx context.Context, mpo *model.ResourceMenu) (menu Menu, err error) {
	menu = Menu{
		Id:        cast.ToString(mpo.ID),
		ParentId:  cast.ToString(mpo.ParentID),
		Path:      mpo.Path,
		Name:      mpo.Name,
		Component: mpo.Component,
		CreatedAt: mpo.CreatedAt.Format(timeFormat),
		UpdatedAt: mpo.UpdatedAt.Format(timeFormat),
		Actions:   nil,
	}
	rrs, err := r.roleResourceRepo.SelectByResourceID(ctx, strconv.FormatUint(uint64(mpo.ID), 10))
	if err != nil {
		err = errors.Wrap(err, "failed to select role resource")
		r.log.WithContext(ctx).Error(err)
		return
	}
	var roles []string
	for _, rr := range rrs {
		roles = append(roles, cast.ToString(rr.RoleID))
	}
	menu.Meta = &MenuMeta{
		Roles:       roles,
		RequireAuth: mpo.RequireAuth == 1,
		Icon:        mpo.Icon,
		Locale:      mpo.Locale,
		Order:       mpo.Order,
		HideInMenu:  mpo.HideInMenu == 1,
		NoAffix:     mpo.NoAffix == 1,
		IgnoreCache: mpo.IgnoreCache == 1,
	}
	if len(mpo.Actions) == 0 {
		return
	}
	var acts []MenuAction
	if err = codec.Unmarshal([]byte(mpo.Actions), &acts); err != nil {
		return
	}
	menu.Actions = acts
	return
}

func (r *ResourceUsecase) SaveMenu(ctx context.Context, menu *pb.MenuRequest) (id string, err error) {
	mpo, err := r.buildMenuPO(ctx, menu)
	if err != nil {
		return
	}
	err = r.menuRepo.Insert(ctx, mpo)
	id = cast.ToString(mpo.ID)
	return
}

func (r *ResourceUsecase) EditMenu(ctx context.Context, menu *pb.MenuRequest) (id string, err error) {
	mpo, err := r.buildMenuPO(ctx, menu)
	if err != nil {
		return
	}
	err = r.menuRepo.Update(ctx, mpo)
	id = cast.ToString(mpo.ID)
	return
}

func (r *ResourceUsecase) RemoveMenu(ctx context.Context, ids []uint) (err error) {
	return r.menuRepo.DeleteByIDs(ctx, ids)
}

func (r *ResourceUsecase) MenuTree(ctx context.Context) (menus []Menu, err error) {
	mpos, err := r.menuRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	return r.buildTree(ctx, mpos)
}

func (r *ResourceUsecase) RoleMenuTree(ctx context.Context, rid uint) (menus []Menu, err error) {
	rMenus, err := r.roleResourceRepo.SelectByRoleID(ctx, uint64(rid), model.ResourceTypeMenu)
	if err != nil {
		return
	}
	var mids []uint
	for _, menu := range rMenus {
		mids = append(mids, cast.ToUint(menu.ResourceID))
	}
	mpos, err := r.menuRepo.SelectByIDs(ctx, mids)
	if err != nil {
		return
	}
	for _, mpo := range mpos {
		var menu Menu
		menu, err = r.buildMenuDO(ctx, mpo)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	return
}

func (r *ResourceUsecase) MenuPage(ctx context.Context, limit, offset int) (list []Menu, total uint32, err error) {
	mpos, err := r.menuRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	list, err = r.buildTree(ctx, mpos)
	if err != nil {
		return
	}
	if offset > len(list)-1 {
		return nil, 0, err
	}
	total = uint32(len(list))
	end := math.Min(limit+offset, len(list))
	list = list[offset:end]
	return
}

func (r *ResourceUsecase) RouterTree(ctx context.Context) (groups []RouterGroup, err error) {
	allRoutes, err := r.routeRepo.SelectAll(ctx)
	if err != nil {
		err = errors.Wrap(err, "ResourceUsecase.RouterTree")
		r.log.WithContext(ctx).Error(err)
		return
	}
	groups = r.buildRouters(ctx, allRoutes)
	return
}

func (r *ResourceUsecase) RoleRouterTree(ctx context.Context, rids ...string) (routers []model.Router, err error) {
	roleRouters, err := r.routeRepo.SelectByRoleIDs(ctx, rids)
	if err != nil {
		return
	}
	for _, rr := range roleRouters {
		routers = append(routers, model.Router{
			Method: rr.Method,
			Path:   rr.Path,
		})
	}
	return
}

func (r *ResourceUsecase) RouterPage(ctx context.Context, limit, offset int) (list []RouterGroup, total uint32, err error) {
	routers, err := r.routeRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	list = r.buildRouters(ctx, routers)
	if offset > len(list)-1 {
		return nil, 0, err
	}
	total = uint32(len(list))
	end := math.Min(limit+offset, len(list))
	list = list[offset:end]
	return
}
