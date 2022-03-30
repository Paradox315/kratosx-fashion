package biz

import "github.com/go-kratos/kratos/v2/log"

type ResourceUsecase struct {
	log *log.Helper
}

func NewResourceUsecase(logger log.Logger) *ResourceUsecase {
	return &ResourceUsecase{
		log: log.NewHelper(logger),
	}
}
