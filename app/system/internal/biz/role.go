package biz

import "github.com/go-kratos/kratos/v2/log"

type RoleUsecase struct {
	roleRepo RoleRepo
	log      *log.Helper
}

func NewRoleUsecase(roleRepo RoleRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		roleRepo: roleRepo,
		log:      log.NewHelper(logger),
	}
}
