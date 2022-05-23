package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
	"kratosx-fashion/app/system/data/model"
	"kratosx-fashion/pkg/xcast"

	api "kratosx-fashion/api/system/v1"
	pb "kratosx-fashion/api/system/v1"
)

type RoleUsecase struct {
	roleRepo         RoleRepo
	roleUserRepo     UserRoleRepo
	roleResourceRepo RoleResourceRepo
	roleRouterRepo   ResourceRouterRepo
	tx               Transaction
	log              *log.Helper
}

func NewRoleUsecase(roleRepo RoleRepo, roleUserRepo UserRoleRepo, roleRouterRepo ResourceRouterRepo, roleResourceRepo RoleResourceRepo, tx Transaction, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		roleRepo:         roleRepo,
		roleUserRepo:     roleUserRepo,
		roleResourceRepo: roleResourceRepo,
		roleRouterRepo:   roleRouterRepo,
		tx:               tx,
		log:              log.NewHelper(log.With(logger, "biz", "role")),
	}
}

func (r *RoleUsecase) buildRoleDto(ctx context.Context, rpo *model.Role) (role *pb.RoleReply) {
	role = &pb.RoleReply{
		Id:          uint64(rpo.ID),
		Name:        rpo.Name,
		Description: rpo.Description,
		CreatedAt:   rpo.CreatedAt.Format(timeFormat),
		UpdatedAt:   rpo.UpdatedAt.Format(timeFormat),
	}
	return
}

func (r *RoleUsecase) buildResources(ctx context.Context, role *pb.RoleRequest) (menus []*model.RoleResource, routes []*model.ResourceRouter) {
	for _, rm := range role.Resources {
		menus = append(menus, &model.RoleResource{
			RoleID:     uint(role.Id),
			ResourceID: rm.ResourceId,
			Type:       model.ResourceType(rm.ResourceType),
		})
	}
	for _, rr := range role.Routers {
		routes = append(routes, &model.ResourceRouter{
			RoleID: cast.ToString(role.Id),
			Path:   rr.Path,
			Method: rr.Method,
		})
	}
	return
}

func (r *RoleUsecase) Save(ctx context.Context, role *pb.RoleRequest) (id uint, err error) {
	rpo := &model.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	if err = r.roleRepo.Insert(ctx, rpo); err != nil {
		r.log.WithContext(ctx).Errorf("role.Insert err: %v", err)
		err = api.ErrorRoleAddFail("角色添加失败")
		return
	}
	role.Id = uint64(rpo.ID)
	resources, routes := r.buildResources(ctx, role)
	if len(resources) > 0 {
		if err = r.roleResourceRepo.Insert(ctx, resources...); err != nil {
			r.log.WithContext(ctx).Errorf("resource.Insert err: %v", err)
			err = api.ErrorRoleAddFail("角色资源添加失败")
			return
		}
	}
	if len(routes) > 0 {
		if err = r.roleRouterRepo.Update(ctx, routes); err != nil {
			r.log.WithContext(ctx).Errorf("router.Update err: %v", err)
			err = api.ErrorRoleAddFail("角色路由添加失败")
			return
		}
	}
	return
}

func (r *RoleUsecase) Edit(ctx context.Context, role *pb.RoleRequest) (err error) {
	rpo := &model.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	rpo.ID = uint(role.Id)
	resources, rRoutes := r.buildResources(ctx, role)
	if err = r.tx.ExecTx(ctx, func(ctx context.Context) error {
		if len(resources) > 0 {
			if err = r.roleResourceRepo.UpdateByRoleID(ctx, rpo.ID, resources); err != nil {
				return err
			}
		}
		if len(rRoutes) > 0 {
			if err = r.roleRouterRepo.Update(ctx, rRoutes); err != nil {
				return err
			}
		}
		return r.roleRepo.Update(ctx, rpo)
	}); err != nil {
		r.log.WithContext(ctx).Errorf("edit role error %s", err.Error())
		err = api.ErrorRoleUpdateFail("编辑角色失败")
		return
	}
	return
}

func (r *RoleUsecase) Remove(ctx context.Context, ids []uint) (err error) {
	if err = r.tx.ExecTx(ctx, func(ctx context.Context) error {
		if err = r.roleRepo.DeleteByIDs(ctx, ids); err != nil {
			return err
		}
		if err = r.roleResourceRepo.DeleteByRoleIDs(ctx, ids); err != nil {
			return err
		}
		if err = r.roleRouterRepo.ClearByRoleIDs(ctx, xcast.ToStringSlice(ids)...); err != nil {
			return err
		}
		return r.roleUserRepo.DeleteByRoleIDs(ctx, ids)
	}); err != nil {
		r.log.WithContext(ctx).Errorf("role.remove %s", err.Error())
		err = api.ErrorRoleDeleteFail("删除角色失败")
	}
	return
}

func (r *RoleUsecase) Get(ctx context.Context, id uint) (role *pb.RoleReply, err error) {
	role = &pb.RoleReply{}
	rpo, err := r.roleRepo.Select(ctx, id)
	if err != nil {
		r.log.WithContext(ctx).Errorf("get role by id %d failed: %v", id, err)
		err = api.ErrorRoleFetchFail("角色获取失败")
		return
	}
	role = r.buildRoleDto(ctx, rpo)
	return
}

func (r *RoleUsecase) Page(ctx context.Context, limit, offset int) (list *pb.ListRoleReply, err error) {
	list = &pb.ListRoleReply{}
	roles, total, err := r.roleRepo.SelectPage(ctx, limit, offset)
	if err != nil {
		r.log.WithContext(ctx).Errorf("roleUseCase.Page.SelectPage: %s", err.Error())
		err = api.ErrorRoleFetchFail("获取角色列表失败")
		return
	}
	for _, role := range roles {
		list.List = append(list.List, r.buildRoleDto(ctx, role))
	}
	list.Total = uint32(total)
	return
}
