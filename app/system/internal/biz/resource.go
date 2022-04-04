package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
	"strconv"
	"strings"
)

var codec = encoding.GetCodec("json")

type ResourceUsecase struct {
	menuRepo  ResourceMenuRepo
	routeRepo ResourceRouterRepo
	tx        Transaction
	log       *log.Helper
}

func NewResourceUsecase(menuRepo ResourceMenuRepo, routeRepo ResourceRouterRepo, tx Transaction, logger log.Logger) *ResourceUsecase {
	return &ResourceUsecase{
		menuRepo:  menuRepo,
		routeRepo: routeRepo,
		tx:        tx,
		log:       log.NewHelper(log.With(logger, "biz", "resource")),
	}
}
func longestCommonPrefix(rs []Router, group string) (prefix string) {
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
func (r *ResourceUsecase) buildRouters(ctx context.Context, rpos []Router) (routers []RouterGroup) {
	routerMap := make(map[string][]Router)
	for _, rpo := range rpos {
		routerMap[rpo.Group] = append(routerMap[rpo.Group], rpo)
	}
	for group, children := range routerMap {
		routers = append(routers, RouterGroup{
			Path:   longestCommonPrefix(children, group),
			Name:   group,
			Router: children,
		})
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

func (r *ResourceUsecase) buildMenuChild(ctx context.Context, menu Menu, menuMap map[string][]Menu) {
	menu.Children = menuMap[menu.Id]
	for _, child := range menu.Children {
		r.buildMenuChild(ctx, child, menuMap)
	}
	return
}

func (r *ResourceUsecase) buildMenuDTO(ctx context.Context, mpo *model.ResourceMenu, rids []string) (menu Menu, err error) {
	_ = copier.Copy(&menu, &mpo)
	menu.Meta = &MenuMeta{
		Locale:      mpo.Locale,
		RequireAuth: mpo.RequireAuth == model.RequireAuthStatusOpen,
		Roles:       rids,
		Icon:        mpo.Icon,
		Order:       mpo.Order,
	}
	menu.Id = cast.ToString(mpo.ID)
	menu.Keepalive = mpo.KeepAlive == model.KeepAliveStatusOpen
	menu.Hidden = mpo.Hidden == model.HiddenStatusShow
	menu.CreatedAt = mpo.CreatedAt.Format(timeFormat)
	menu.UpdatedAt = mpo.UpdatedAt.Format(timeFormat)
	var acts []MenuAction
	if err = codec.Unmarshal([]byte(mpo.Actions), &acts); err != nil {
		return
	}
	menu.Actions = acts
	return
}

func (r *ResourceUsecase) SaveMenu(ctx context.Context, menu *pb.MenuRequest) (id string, err error) {
	mpo := &model.ResourceMenu{}
	if err = copier.Copy(&mpo, &menu); err != nil {
		return
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
	return cast.ToString(mpo.ID), r.menuRepo.Insert(ctx, mpo)
}

func (r *ResourceUsecase) EditMenu(ctx context.Context, menu *pb.MenuRequest) (id string, err error) {
	mpo := &model.ResourceMenu{}
	if err = copier.Copy(&mpo, &menu); err != nil {
		return
	}
	mpo.ID = cast.ToUint(menu.Id)
	mpo.Locale = menu.Meta.Locale
	if menu.Meta.RequireAuth {
		mpo.RequireAuth = 1
	}
	mpo.Icon = menu.Meta.Icon
	mpo.Order = menu.Meta.Order
	if menu.ParentId == "" {
		mpo.ParentID = 0
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
	return cast.ToString(mpo.ID), r.menuRepo.Update(ctx, mpo)
}

func (r *ResourceUsecase) RemoveMenu(ctx context.Context, ids []uint) (err error) {
	return r.menuRepo.DeleteByIDs(ctx, ids)
}
func (r *ResourceUsecase) MenuTree(ctx context.Context) (menus []Menu, err error) {
	mpos, err := r.menuRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	for _, mpo := range mpos {
		var menu Menu
		menu, err = r.buildMenuDTO(ctx, mpo, nil)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	menuMap := r.groupMenu(ctx, menus)
	for _, menu := range menuMap["0"] {
		r.buildMenuChild(ctx, menu, menuMap)
	}
	return
}

func (r *ResourceUsecase) RoleMenuTree(ctx context.Context, rids ...uint) (menus []Menu, err error) {
	mpos, err := r.menuRepo.SelectByIDs(ctx, rids)
	if err != nil {
		return
	}
	var roles []string
	for _, rid := range rids {
		roles = append(roles, strconv.FormatUint(uint64(rid), 10))
	}
	for _, mpo := range mpos {
		var menu Menu
		menu, err = r.buildMenuDTO(ctx, mpo, roles)
		if err != nil {
			return
		}
		menus = append(menus, menu)
	}
	menuMap := r.groupMenu(ctx, menus)
	for _, child := range menuMap["0"] {
		r.buildMenuChild(ctx, child, menuMap)
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
	allRoutes, err := r.routeRepo.SelectAll(ctx)
	if err != nil {
		return
	}
	groups = r.buildRouters(ctx, allRoutes)
	//ownedRoutes, err := r.routeRepo.SelectByRoleIDs(ctx, rids)
	if err != nil {
		return
	}
	return
}

func (r *ResourceUsecase) EditRouters(ctx context.Context, rid string, routers []RouterGroup) (err error) {
	var rrs []*model.ResourceRouter
	//for _, group := range routers {
	//	if group.Owned {
	//		rrs = append(rrs, &model.ResourceRouter{
	//			RoleID: rid,
	//			Path:   group.Path,
	//			Method: strings.Join(group.Methods, "|"),
	//		})
	//		continue
	//	}
	//	for _, router := range group.Router {
	//		if router.Owned {
	//			rrs = append(rrs, &model.ResourceRouter{
	//				RoleID: rid,
	//				Path:   router.Path,
	//				Method: router.Method,
	//			})
	//		}
	//	}
	//}
	return r.routeRepo.Update(ctx, rrs)
}
