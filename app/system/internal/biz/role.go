package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/xcast"
	"strconv"
)

type RoleUsecase struct {
	roleRepo         RoleRepo
	roleUserRepo     UserRoleRepo
	roleResourceRepo RoleResourceRepo
	tx               Transaction
	log              *log.Helper
}

func NewRoleUsecase(roleRepo RoleRepo, roleUserRepo UserRoleRepo, roleResourceRepo RoleResourceRepo, tx Transaction, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		roleRepo:         roleRepo,
		roleUserRepo:     roleUserRepo,
		roleResourceRepo: roleResourceRepo,
		tx:               tx,
		log:              log.NewHelper(log.With(logger, "biz", "role")),
	}
}

func (r *RoleUsecase) buildRoleReply(ctx context.Context, rpo *model.Role) (role *pb.RoleReply, err error) {
	role = &pb.RoleReply{}
	if err = copier.Copy(&role, &rpo); err != nil {
		return
	}
	rrs, err := r.roleResourceRepo.SelectByRoleID(ctx, uint64(rpo.ID))
	for _, rr := range rrs {
		role.RoleResources = append(role.RoleResources, &pb.RoleResource{
			RoleId:       strconv.FormatUint(rr.RoleID, 10),
			ResourceId:   strconv.FormatUint(rr.ResourceID, 10),
			ResourceType: uint32(rr.Type),
		})
	}
	role.Id = cast.ToString(rpo.ID)
	role.CreatedAt = rpo.CreatedAt.Format(timeFormat)
	role.UpdatedAt = rpo.UpdatedAt.Format(timeFormat)
	return
}

func (r *RoleUsecase) Save(ctx context.Context, role *pb.RoleRequest) (id string, err error) {
	rpo := &model.Role{}
	if err = copier.Copy(&rpo, &role); err != nil {
		return
	}
	var rrs []*model.RoleResource
	for _, rr := range role.RoleResources {
		rrs = append(rrs, &model.RoleResource{
			RoleID:     cast.ToUint64(rr.RoleId),
			ResourceID: cast.ToUint64(rr.ResourceId),
			Type:       model.ResourceType(rr.ResourceType),
		})
	}
	err = r.roleRepo.Insert(ctx, rpo)
	if err != nil {
		return
	}
	id = strconv.Itoa(int(rpo.ID))
	if len(rrs) > 0 {
		err = r.roleResourceRepo.Insert(ctx, rrs...)
	}
	return
}

func (r *RoleUsecase) Edit(ctx context.Context, role *pb.RoleRequest) (id string, err error) {
	rpo := &model.Role{}
	if err = copier.Copy(&rpo, &role); err != nil {
		return
	}
	rpo.ID = cast.ToUint(role.Id)
	var rrs []*model.RoleResource
	for _, rr := range role.RoleResources {
		rrs = append(rrs, &model.RoleResource{
			RoleID:     cast.ToUint64(rr.RoleId),
			ResourceID: cast.ToUint64(rr.ResourceId),
			Type:       model.ResourceType(rr.ResourceType),
		})
	}
	if len(rrs) > 0 {
		err = r.roleResourceRepo.UpdateByRoleID(ctx, uint64(rpo.ID), rrs)
		if err != nil {
			return
		}
	}
	err = r.roleRepo.Update(ctx, rpo)
	id = role.Id
	return
}

func (r *RoleUsecase) Remove(ctx context.Context, ids []uint) (err error) {
	return r.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = r.roleRepo.DeleteByIDs(ctx, ids)
		if err != nil {
			return err
		}
		err = r.roleResourceRepo.DeleteByRoleIDs(ctx, xcast.ToUint64Slice(ids))
		if err != nil {
			return err
		}
		return r.roleUserRepo.DeleteByRoleIDs(ctx, xcast.ToUint64Slice(ids))
	})
}

func (r *RoleUsecase) Get(ctx context.Context, id uint) (role *pb.RoleReply, err error) {
	role = &pb.RoleReply{}
	rpo, err := r.roleRepo.Select(ctx, id)
	if err != nil {
		return
	}
	return r.buildRoleReply(ctx, rpo)
}

func (r *RoleUsecase) List(ctx context.Context, limit, offset int) (list *pb.ListRoleReply, err error) {
	list = &pb.ListRoleReply{}
	roles, total, err := r.roleRepo.List(ctx, limit, offset)
	if err != nil {
		return
	}
	for _, role := range roles {
		var rr *pb.RoleReply
		rr, err = r.buildRoleReply(ctx, role)
		if err != nil {
			return
		}
		list.Roles = append(list.Roles, rr)
	}
	list.Total = uint32(total)
	return
}
