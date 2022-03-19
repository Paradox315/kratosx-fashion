package biz

import "github.com/go-kratos/kratos/v2/log"

type MenuUsecase struct {
	menuRepo           MenuRepo
	menuActionRepo     MenuActionRepo
	actionResourceRepo MenuActionResourceRepo
	log                *log.Helper
}

func NewMenuUsecase(menuRepo MenuRepo, menuActionRepo MenuActionRepo, actionResourceRepo MenuActionResourceRepo, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{
		menuRepo:           menuRepo,
		menuActionRepo:     menuActionRepo,
		actionResourceRepo: actionResourceRepo,
		log:                log.NewHelper(logger),
	}
}
