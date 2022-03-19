package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jassue/go-storage/storage"
)

type Greeter struct {
	Hello string
}

type GreeterRepo interface {
	CreateGreeter(context.Context, *Greeter) error
	UpdateGreeter(context.Context, *Greeter) error
}

type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
	disk storage.Storage
}

func NewGreeterUsecase(repo GreeterRepo, logger log.Logger, disk storage.Storage) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger), disk: disk}
}

func (uc *GreeterUsecase) Create(ctx context.Context, g *Greeter) error {
	return uc.repo.CreateGreeter(ctx, g)
}

func (uc *GreeterUsecase) Update(ctx context.Context, g *Greeter) error {
	return uc.repo.UpdateGreeter(ctx, g)
}
