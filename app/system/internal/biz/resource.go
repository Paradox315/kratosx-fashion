package biz

import (
	"context"
	"github.com/casbin/casbin/v2/util"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/math"
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

func match(r model.Router, p model.ResourceRouter) bool {
	return (util.KeyMatch(r.Path, p.Path) || util.KeyMatch2(r.Path, p.Path)) && util.RegexMatch(r.Method, p.Method)
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
			Path:   longestCommonPrefix(children, group),
			Name:   group,
			Router: children,
			Methods: func(routers []model.Router) string {
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
	mpo = &model.ResourceMenu{}
	if err = copier.Copy(&mpo, &menu); err != nil {
		return
	}
	if len(menu.Id) > 0 && menu.Id != "0" {
		mpo.ID = cast.ToUint(menu.Id)
	}
	if len(menu.ParentId) > 0 && menu.ParentId != "0" {
		mpo.ParentID = cast.ToUint(menu.ParentId)
	}
	mpo.Locale = menu.Meta.Locale
	if menu.Meta.RequireAuth {
		mpo.RequireAuth = 1
	}
	mpo.Icon = menu.Meta.Icon
	mpo.Order = menu.Meta.Order
	if menu.ParentId == "" {
		mpo.ParentID = 0
	}
	if menu.Hidden {
		mpo.Hidden = 1
	}
	if menu.Keepalive {
		mpo.Keepalive = 1
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
		menu, err = r.buildMenuDTO(ctx, mpo)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	menuMap := r.groupMenu(ctx, menus)
	menus = menuMap["0"]
	for i := 0; i < len(menus); i++ {
		r.buildMenuChild(ctx, &menus[i], menuMap)
	}
	return
}

func (r *ResourceUsecase) groupMenu(ctx context.Context, menus []Menu) (menuMap map[string][]Menu) {
	menuMap = make(map[string][]Menu)
	for _, menu := range menus {
		menuMap[menu.ParentId] = append(menuMap[menu.ParentId], menu)
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

func (r *ResourceUsecase) buildMenuDTO(ctx context.Context, mpo *model.ResourceMenu) (menu Menu, err error) {
	_ = copier.Copy(&menu, &mpo)
	menu.ParentId = cast.ToString(mpo.ParentID)
	rrs, err := r.roleResourceRepo.SelectByResourceID(ctx, uint64(mpo.ID))
	if err != nil {
		return
	}
	var roles []string
	for _, rr := range rrs {
		roles = append(roles, cast.ToString(rr.RoleID))
	}
	menu.Meta = &MenuMeta{
		Locale:      mpo.Locale,
		RequireAuth: mpo.RequireAuth == model.RequireAuthStatusOpen,
		Roles:       roles,
		Icon:        mpo.Icon,
		Order:       mpo.Order,
	}
	menu.Id = cast.ToString(mpo.ID)
	menu.Keepalive = mpo.Keepalive == model.KeepAliveStatusOpen
	menu.Hidden = mpo.Hidden == model.HiddenStatusShow
	menu.CreatedAt = mpo.CreatedAt.Format(timeFormat)
	menu.UpdatedAt = mpo.UpdatedAt.Format(timeFormat)
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
	rMenus, err := r.roleResourceRepo.SelectByRoleID(ctx, uint64(rid))
	if err != nil {
		return
	}
	var mids []uint
	for _, menu := range rMenus {
		mids = append(mids, uint(menu.ResourceID))
	}
	mpos, err := r.menuRepo.SelectByIDs(ctx, mids)
	if err != nil {
		return
	}
	return r.buildTree(ctx, mpos)
}

func (r *ResourceUsecase) MenuPage(ctx context.Context, limit, offset int) (menus []Menu, total int64, err error) {
	mpos, total, err := r.menuRepo.SelectPage(ctx, limit, offset)
	if err != nil {
		return
	}
	menus, err = r.buildTree(ctx, mpos)
	if err != nil {
		return
	}
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

func (r *ResourceUsecase) RoleRouterTree(ctx context.Context, rids ...string) (groups []RouterGroup, err error) {
	var filtered []model.Router
	allRoutes, err := r.routeRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	myRoutes, err := r.routeRepo.SelectByRoleIDs(ctx, rids)
	if err != nil {
		return
	}
	for _, route := range allRoutes {
		for _, myRoute := range myRoutes {
			if !match(route, myRoute) {
				continue
			}
			filtered = append(filtered, route)
		}
	}
	groups = r.buildRouters(ctx, filtered)
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
