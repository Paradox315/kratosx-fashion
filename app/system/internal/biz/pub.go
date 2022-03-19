package biz

import "github.com/go-kratos/kratos/v2/log"

type PublicUsecase struct {
	userRepo UserRepo
	logRepo  LoginLogRepo
	log      *log.Helper
}

func NewPublicUsecase(userRepo UserRepo, logRepo LoginLogRepo, logger log.Logger) *PublicUsecase {
	return &PublicUsecase{
		userRepo: userRepo,
		logRepo:  logRepo,
		log:      log.NewHelper(logger),
	}
}
