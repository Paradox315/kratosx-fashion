package biz

import "github.com/go-kratos/kratos/v2/log"

type RoleUsecase struct {
	roleRepo     RoleRepo
	roleMenuRepo RoleMenuRepo
	log          *log.Helper
}

func NewRoleUsecase(roleRepo RoleRepo, roleMenuRepo RoleMenuRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		roleRepo:     roleRepo,
		roleMenuRepo: roleMenuRepo,
		log:          log.NewHelper(logger),
	}
}
