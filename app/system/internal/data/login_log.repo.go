package data

import (
	"context"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type LoginLogRepo struct {
	dao *Data
	log *log.Helper
}

func NewLoginLogRepo(data *Data, logger log.Logger) biz.LoginLogRepo {
	return &LoginLogRepo{dao: data, log: log.NewHelper(logger)}
}

func (l *LoginLogRepo) Select(ctx context.Context, id uint64) (*model.LoginLog, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LoginLogRepo) ListByUserID(ctx context.Context, u uint64, request *pb.ListRequest, option ...*pb.QueryOption) ([]*model.LoginLog, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LoginLogRepo) Insert(ctx context.Context, loginLog *model.LoginLog) error {
	//TODO implement me
	panic("implement me")
}

func (l *LoginLogRepo) Delete(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (l *LoginLogRepo) DeleteByUserID(ctx context.Context, uid uint64) error {
	//TODO implement me
	panic("implement me")
}

func (l *LoginLogRepo) DeleteByUserIDs(ctx context.Context, uids []uint64) error {
	//TODO implement me
	panic("implement me")
}
