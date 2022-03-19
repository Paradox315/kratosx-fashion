package biz

import "github.com/go-kratos/kratos/v2/log"

type UserUsecase struct {
	userRepo     UserRepo
	userRoleRepo UserRoleRepo
	log          *log.Helper
}

func NewUserUsecase(userRepo UserRepo, userRoleRepo UserRoleRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		log:          log.NewHelper(logger),
	}
}
