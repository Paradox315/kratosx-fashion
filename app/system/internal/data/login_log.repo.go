package data

import (
	"context"

	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/app/system/internal/data/query"
	"kratosx-fashion/pkg/option"
	"kratosx-fashion/pkg/pagination"

	"github.com/go-kratos/kratos/v2/log"

	pb "kratosx-fashion/api/system/v1"
)

type LoginLogRepo struct {
	dao      *Data
	log      *log.Helper
	baseRepo *query.Query
}

func NewLoginLogRepo(data *Data, logger log.Logger) biz.LoginLogRepo {
	return &LoginLogRepo{
		dao:      data,
		log:      log.NewHelper(logger),
		baseRepo: query.Use(data.DB),
	}
}

func (l *LoginLogRepo) Select(ctx context.Context, id uint) (loginLog *model.LoginLog, err error) {
	lr := l.baseRepo.LoginLog
	return lr.WithContext(ctx).Where(lr.ID.Eq(id)).First()
}

func (l *LoginLogRepo) ListByUserID(ctx context.Context, id uint64, req *pb.ListRequest, opts ...*pb.QueryOption) (logs []*model.LoginLog, total int64, err error) {
	orders, keywords := option.Parse(opts...)
	limit, offset, err := pagination.Parse(req)
	if err != nil {
		l.log.WithContext(ctx).Error("pagination.Parse error", err)
		return
	}
	keywords["user_id"] = id
	tx := l.dao.DB.Where(keywords).Limit(limit).Offset(offset)
	for _, order := range orders {
		tx.Order(order)
	}
	err = tx.Count(&total).Find(logs).Error
	return
}

func (l *LoginLogRepo) Insert(ctx context.Context, loginLog *model.LoginLog) error {
	lr := l.baseRepo.LoginLog
	return lr.WithContext(ctx).Create(loginLog)
}

func (l *LoginLogRepo) Delete(ctx context.Context, id uint) error {
	lr := l.baseRepo.LoginLog
	_, err := lr.WithContext(ctx).Where(lr.ID.Eq(id)).Delete()
	return err
}

func (l *LoginLogRepo) DeleteByUserID(ctx context.Context, uid uint64) error {
	lr := l.baseRepo.LoginLog
	_, err := lr.WithContext(ctx).Where(lr.UserID.Eq(uid)).Delete()
	return err
}

func (l *LoginLogRepo) DeleteByUserIDs(ctx context.Context, uids []uint64) error {
	lr := l.baseRepo.LoginLog
	_, err := lr.WithContext(ctx).Where(lr.UserID.In(uids...)).Delete()
	return err
}
