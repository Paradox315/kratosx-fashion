package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/xcast"
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

func (r *RoleUsecase) buildRoleReply(ctx context.Context, rpo *model.Role) (role *pb.RoleReply, err error) {
	role = &pb.RoleReply{
		Id:          cast.ToString(rpo.ID),
		Name:        rpo.Name,
		Description: rpo.Description,
		CreatedAt:   rpo.CreatedAt.Format(timeFormat),
		UpdatedAt:   rpo.UpdatedAt.Format(timeFormat),
	}
	return
}

func (r *RoleUsecase) buildResources(ctx context.Context, role *pb.RoleRequest) (menus []*model.RoleResource, routes []model.ResourceRouter) {
	for _, rm := range role.RoleResources {
		menus = append(menus, &model.RoleResource{
			RoleID:     cast.ToUint64(role.Id),
			ResourceID: rm.ResourceId,
			Type:       model.ResourceType(rm.ResourceType),
		})
	}
	for _, rr := range role.RoleRouters {
		routes = append(routes, model.ResourceRouter{
			RoleID: role.Id,
			Path:   rr.Path,
			Method: rr.Method,
		})
	}
	return
}

func (r *RoleUsecase) Save(ctx context.Context, role *pb.RoleRequest) (id string, err error) {
	rpo := &model.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	if err = r.roleRepo.Insert(ctx, rpo); err != nil {
		return
	}
	id = cast.ToString(rpo.ID)
	role.Id = id
	resources, rRoutes := r.buildResources(ctx, role)
	if len(resources) > 0 {
		if err = r.roleResourceRepo.Insert(ctx, resources...); err != nil {
			return
		}
	}
	if len(rRoutes) > 0 {
		if err = r.roleRouterRepo.Update(ctx, rRoutes); err != nil {
			return
		}
	}
	return
}

func (r *RoleUsecase) Edit(ctx context.Context, role *pb.RoleRequest) (id string, err error) {
	id = role.Id
	rpo := &model.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	rpo.ID = cast.ToUint(role.Id)
	resources, rRoutes := r.buildResources(ctx, role)

	err = r.tx.ExecTx(ctx, func(ctx context.Context) error {
		if len(resources) > 0 {
			if err = r.roleResourceRepo.UpdateByRoleID(ctx, uint64(rpo.ID), resources); err != nil {
				return err
			}
		}
		if len(rRoutes) > 0 {
			if err = r.roleRouterRepo.Update(ctx, rRoutes); err != nil {
				return err
			}
		}
		return r.roleRepo.Update(ctx, rpo)
	})
	return
}

func (r *RoleUsecase) Remove(ctx context.Context, ids []uint) (err error) {
	return r.tx.ExecTx(ctx, func(ctx context.Context) error {
		if err = r.roleRepo.DeleteByIDs(ctx, ids); err != nil {
			return err
		}
		if err = r.roleResourceRepo.DeleteByRoleIDs(ctx, xcast.ToUint64Slice(ids)); err != nil {
			return err
		}
		if err = r.roleRouterRepo.ClearByRoleIDs(ctx, xcast.ToStringSlice(ids)); err != nil {
			return err
		}
		return r.roleUserRepo.DeleteByRoleIDs(ctx, xcast.ToUint64Slice(ids))
	})
}

func (r *RoleUsecase) Get(ctx context.Context, id uint) (role *pb.RoleReply, err error) {
	role = &pb.RoleReply{}
	rpo, err := r.roleRepo.Select(ctx, id)
	if err != nil {
		err = errors.Wrap(err, "roleUseCase.Get.Select")
		r.log.WithContext(ctx).Error(err)
		return
	}
	return r.buildRoleReply(ctx, rpo)
}

func (r *RoleUsecase) Page(ctx context.Context, limit, offset int) (list *pb.ListRoleReply, err error) {
	list = &pb.ListRoleReply{}
	roles, total, err := r.roleRepo.SelectPage(ctx, limit, offset)
	if err != nil {
		return
	}
	for _, role := range roles {
		var rr *pb.RoleReply
		rr, err = r.buildRoleReply(ctx, role)
		if err != nil {
			return
		}
		list.List = append(list.List, rr)
	}
	list.Total = uint32(total)
	return
}
