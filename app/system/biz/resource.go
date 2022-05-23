package biz

import (
	"context"
	"strconv"
	"strings"

	"kratosx-fashion/app/system/data/model"
	"kratosx-fashion/pkg/math"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	api "kratosx-fashion/api/system/v1"
	pb "kratosx-fashion/api/system/v1"
)

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

func longestCommonPrefix(rs []*model.Router, group string) (prefix string) {
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

func (r *ResourceUsecase) buildRouters(ctx context.Context, rpos []*model.Router) (routers []*RouterGroup) {
	routerMap := lo.GroupBy(rpos, func(rpo *model.Router) string {
		return rpo.Group
	})
	for group, children := range routerMap {
		routers = append(routers, &RouterGroup{
			Path:     longestCommonPrefix(children, group),
			Name:     group,
			Children: children,
			Method: func(routers []*model.Router) string {
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
func (r *ResourceUsecase) buildMenuPO(ctx context.Context, menu *pb.MenuRequest) (mpo *model.ResourceMenu) {
	mpo = &model.ResourceMenu{
		Name:        menu.Name,
		Path:        menu.Path,
		Component:   menu.Component,
		Description: menu.Description,
		ParentID:    uint(menu.ParentId),
		Locale:      menu.Meta.Locale,
		RequireAuth: model.AuthType(lo.Ternary(menu.Meta.RequireAuth, model.RequireAuth, model.NotRequireAuth)),
		Icon:        menu.Meta.Icon,
		Order:       menu.Meta.Order,
		HideInMenu:  model.HideType(lo.Ternary(menu.Meta.HideInMenu, model.HideInMenu, model.ShowInMenu)),
		NoAffix:     model.AffixType(lo.Ternary(menu.Meta.NoAffix, model.NoAffix, model.Affix)),
		IgnoreCache: model.CacheType(lo.Ternary(menu.Meta.IgnoreCache, model.IgnoreCache, model.Cache)),
	}
	mpo.ID = uint(menu.Id)
	var acts []*model.ResourceAction
	for _, action := range menu.Actions {
		acts = append(acts, &model.ResourceAction{
			Name: action.Name,
			Code: action.Code,
		})
	}
	mpo.SetActions(acts)
	return
}

func (r *ResourceUsecase) buildTree(ctx context.Context, mpos []*model.ResourceMenu) (menus []*Menu, err error) {
	for _, mpo := range mpos {
		var menu *Menu
		menu, err = r.buildMenuDO(ctx, mpo)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	menuMap := lo.GroupBy(menus, func(menu *Menu) uint {
		return menu.ParentId
	})
	menus = menuMap[0]
	for i := 0; i < len(menus); i++ {
		r.buildMenuChild(ctx, menus[i], menuMap)
	}
	return
}

func (r *ResourceUsecase) buildMenuChild(ctx context.Context, menu *Menu, menuMap map[uint][]*Menu) {
	menu.Children = menuMap[menu.Id]
	for _, child := range menu.Children {
		r.buildMenuChild(ctx, child, menuMap)
	}
	return
}

func (r *ResourceUsecase) buildMenuDO(ctx context.Context, mpo *model.ResourceMenu) (menu *Menu, err error) {
	menu = &Menu{
		Id:          mpo.ID,
		ParentId:    mpo.ParentID,
		Path:        mpo.Path,
		Name:        mpo.Name,
		Description: mpo.Description,
		Component:   mpo.Component,
		CreatedAt:   mpo.CreatedAt.Format(timeFormat),
		UpdatedAt:   mpo.UpdatedAt.Format(timeFormat),
	}
	rrs, err := r.roleResourceRepo.SelectByResourceID(ctx, strconv.FormatUint(uint64(mpo.ID), 10), model.ResourceTypeMenu)
	if err != nil {
		r.log.WithContext(ctx).Errorf("select role resource by resource id error: %v", err)
		err = api.ErrorMenuFetchFail("获取菜单权限失败")
		return
	}
	var roles []uint64
	for _, rr := range rrs {
		roles = append(roles, uint64(rr.RoleID))
	}
	menu.Meta = &MenuMeta{
		Roles:       roles,
		RequireAuth: mpo.RequireAuth == model.RequireAuth,
		Icon:        mpo.Icon,
		Locale:      mpo.Locale,
		Order:       mpo.Order,
		HideInMenu:  mpo.HideInMenu == model.HideInMenu,
		NoAffix:     mpo.NoAffix == model.NoAffix,
		IgnoreCache: mpo.IgnoreCache == model.IgnoreCache,
	}
	for _, action := range mpo.GetActions() {
		menu.Actions = append(menu.Actions, &MenuAction{
			Name: action.Name,
			Code: action.Code,
		})
	}
	return
}

func (r *ResourceUsecase) SaveMenu(ctx context.Context, menu *pb.MenuRequest) (id uint, err error) {
	mpo := r.buildMenuPO(ctx, menu)
	if err = r.menuRepo.Insert(ctx, mpo); err != nil {
		r.log.WithContext(ctx).Errorf("save menu failed: %v", err)
		err = api.ErrorMenuAddFail("添加菜单失败")
		return
	}
	return mpo.ID, nil
}

func (r *ResourceUsecase) EditMenu(ctx context.Context, menu *pb.MenuRequest) (err error) {
	mpo := r.buildMenuPO(ctx, menu)
	if err = r.menuRepo.Update(ctx, mpo); err != nil {
		r.log.WithContext(ctx).Errorf("edit menu failed: %v", err)
		err = api.ErrorMenuUpdateFail("编辑菜单失败")
		return
	}
	return
}

func (r *ResourceUsecase) RemoveMenu(ctx context.Context, ids []uint) (err error) {
	if err = r.menuRepo.DeleteByIDs(ctx, ids); err != nil {
		r.log.WithContext(ctx).Errorf("remove menu failed: %v", err)
		err = api.ErrorMenuDeleteFail("删除菜单失败")
		return
	}
	return
}

func (r *ResourceUsecase) MenuTree(ctx context.Context) (menus []*Menu, err error) {
	mpos, err := r.menuRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	tree, err := r.buildTree(ctx, mpos)
	if err != nil {
		r.log.WithContext(ctx).Errorf("build menu tree failed: %v", err)
		err = api.ErrorMenuFetchFail("菜单树构建失败")
		return
	}
	return tree, nil
}
func (r *ResourceUsecase) RoleMenuTree(ctx context.Context, rid uint) (tree []*Menu, err error) {
	rrs, err := r.roleResourceRepo.SelectByRoleID(ctx, rid, model.ResourceTypeMenu)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get role resource failed: %v", err)
		err = api.ErrorRoleFetchFail("获取角色资源失败")
		return
	}
	var mids []uint
	for _, rr := range rrs {
		mids = append(mids, cast.ToUint(rr.ResourceID))
	}
	mpos, err := r.menuRepo.SelectByIDs(ctx, mids)
	if err != nil {
		return
	}
	tree, err = r.buildTree(ctx, mpos)
	if err != nil {
		r.log.WithContext(ctx).Errorf("build menu tree failed: %v", err)
		err = api.ErrorMenuFetchFail("菜单树构建失败")
		return
	}
	return
}
func (r *ResourceUsecase) RoleMenuList(ctx context.Context, rid uint) (menus []*Menu, err error) {
	rrs, err := r.roleResourceRepo.SelectByRoleID(ctx, rid, model.ResourceTypeMenu)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get role resource failed: %v", err)
		err = api.ErrorRoleFetchFail("获取角色资源失败")
		return
	}
	var mids []uint
	for _, menu := range rrs {
		mids = append(mids, cast.ToUint(menu.ResourceID))
	}
	mpos, err := r.menuRepo.SelectByIDs(ctx, mids)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get menu by ids failed: %v", err)
		err = api.ErrorMenuFetchFail("获取菜单失败")
		return
	}
	for _, mpo := range mpos {
		var menu *Menu
		menu, err = r.buildMenuDO(ctx, mpo)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	return
}

func (r *ResourceUsecase) MenuPage(ctx context.Context, limit, offset int) (list []*Menu, total uint32, err error) {
	mpos, err := r.menuRepo.SelectAll(ctx)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get menu page failed: %v", err)
		err = api.ErrorMenuFetchFail("获取菜单失败")
		return
	}
	list, err = r.buildTree(ctx, mpos)
	if err != nil {
		r.log.WithContext(ctx).Errorf("build menu tree failed: %v", err)
		err = api.ErrorMenuFetchFail("菜单树构建失败")
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

func (r *ResourceUsecase) RouterTree(ctx context.Context) (groups []*RouterGroup, err error) {
	allRoutes, err := r.routeRepo.SelectAll(ctx)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get all routes failed: %v", err)
		err = api.ErrorRouterFetchFail("获取路由失败")
		return
	}
	groups = r.buildRouters(ctx, allRoutes)
	return
}

func (r *ResourceUsecase) RoleRouterList(ctx context.Context, rid string) (routers []*model.Router, err error) {
	roleRouters, err := r.routeRepo.SelectByRoleID(ctx, rid)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get role routers failed: %v", err)
		err = api.ErrorRouterFetchFail("获取角色路由失败")
		return
	}
	for _, rr := range roleRouters {
		routers = append(routers, &model.Router{
			Method: rr.Method,
			Path:   rr.Path,
		})
	}
	return
}
func (r *ResourceUsecase) RoleActionList(ctx context.Context, rid uint) (actions []string, err error) {
	roleActions, err := r.roleResourceRepo.SelectByRoleID(ctx, rid, model.ResourceTypeAction)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get role actions failed: %v", err)
		err = api.ErrorActionFetchFail("获取角色动作失败")
		return
	}
	for _, ra := range roleActions {
		actions = append(actions, ra.ResourceID)
	}
	return
}
func (r *ResourceUsecase) RouterPage(ctx context.Context, limit, offset int) (list []*RouterGroup, total uint32, err error) {
	routers, err := r.routeRepo.SelectAll(ctx)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get all routers failed: %v", err)
		err = api.ErrorRouterFetchFail("获取路由失败")
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
