package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
)

type userRepo struct {
	dao *Data
	log *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		dao: data,
		log: log.NewHelper(logger),
	}
}

func (u userRepo) SelectByID(ctx context.Context, u2 uint64) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) SelectByUsername(ctx context.Context, s string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) SelectPasswordByID(ctx context.Context, u2 uint64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) List(ctx context.Context, reply *pb.ListUserReply, option ...*pb.QueryOption) ([]*model.User, uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) Insert(ctx context.Context, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) Update(ctx context.Context, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) Delete(ctx context.Context, u2 uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) DeleteByIDs(ctx context.Context, uint64s []uint64) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) ExistByUserName(ctx context.Context, s string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
