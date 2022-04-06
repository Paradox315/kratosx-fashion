package repo

import (
	"context"
	iploc "github.com/ip2location/ip2location-go"
	ua "github.com/mileusna/useragent"
	"github.com/pkg/errors"
	"kratosx-fashion/app/system/internal/biz"
	"kratosx-fashion/app/system/internal/data"
	"kratosx-fashion/app/system/internal/data/linq"
	"kratosx-fashion/app/system/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type LoginLogRepo struct {
	dao      *data.Data
	log      *log.Helper
	ipdb     *iploc.DB
	baseRepo *linq.Query
}

func NewLoginLogRepo(data *data.Data, logger log.Logger, ipdb *iploc.DB) biz.LoginLogRepo {
	return &LoginLogRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "login_log")),
		ipdb:     ipdb,
		baseRepo: linq.Use(data.DB),
	}
}

func (l *LoginLogRepo) SelectLocation(ctx context.Context, ip string) (loc *biz.Location, err error) {
	if len(ip) == 0 || ip == "127.0.0.1" || ip == "localhost" || ip == "::1" {
		loc = &biz.Location{
			Country:  "本地",
			Region:   "本地",
			City:     "本地",
			Position: nil,
		}
	}
	result, err := l.ipdb.Get_all(ip)
	if err != nil {
		err = errors.Wrap(err, "ipdb.Get_all")
		l.log.Error("ip2location get_all error", err)
		return nil, err
	}
	loc = &biz.Location{
		Country: result.Country_short,
		Region:  result.Region,
		City:    result.City,
		Position: map[string]float32{
			"lat": result.Latitude,
			"lng": result.Longitude,
		},
	}
	return
}

func (l *LoginLogRepo) SelectAgent(ctx context.Context, agentStr string) (agent *biz.Agent, err error) {
	result := ua.Parse(agentStr)
	agent = &biz.Agent{
		Name:   result.Name,
		OS:     result.OS,
		Device: result.Device,
	}
	if result.Desktop {
		agent.DeviceType = model.DeviceType_PC
	} else if result.Mobile {
		agent.DeviceType = model.DeviceType_Mobile
	} else if result.Tablet {
		agent.DeviceType = model.DeviceType_Pad
	} else if result.Bot {
		agent.DeviceType = model.DeviceType_Bot
	}
	return
}

func (l *LoginLogRepo) Select(ctx context.Context, id uint) (loginLog *model.LoginLog, err error) {
	lr := l.baseRepo.LoginLog
	loginLog, err = lr.WithContext(ctx).Where(lr.ID.Eq(id)).First()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.Select")
		l.log.WithContext(ctx).Error("login_log select error", err)
		return
	}
	return
}

func (l *LoginLogRepo) SelectPageByUserID(ctx context.Context, id uint64, limit, offset int) (logs []*model.LoginLog, total int64, err error) {
	if err != nil {
		l.log.WithContext(ctx).Error("pagination.Parse error", err)
		return
	}
	lr := l.baseRepo.LoginLog
	tx := lr.WithContext(ctx).Where(lr.UserID.Eq(id)).Limit(limit).Offset(offset)
	total, err = tx.Count()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.SelectPageByUserID")
		l.log.WithContext(ctx).Error("pagination.Count error", err)
		return
	}
	logs, err = tx.Find()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.SelectPageByUserID")
		l.log.WithContext(ctx).Error("pagination.Find error", err)
	}
	return
}

func (l *LoginLogRepo) Insert(ctx context.Context, loginLog *model.LoginLog) error {
	lr := l.baseRepo.LoginLog
	err := lr.WithContext(ctx).Create(loginLog)
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.Insert")
		l.log.WithContext(ctx).Error("login_log insert error", err)
		return err
	}
	return nil
}

func (l *LoginLogRepo) Delete(ctx context.Context, id uint) error {
	lr := l.baseRepo.LoginLog
	_, err := lr.WithContext(ctx).Where(lr.ID.Eq(id)).Delete()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.Delete")
		l.log.WithContext(ctx).Error("login_log delete error", err)
		return err
	}
	return nil
}

func (l *LoginLogRepo) DeleteByUserID(ctx context.Context, uid uint64) error {
	lr := l.baseRepo.LoginLog
	_, err := lr.WithContext(ctx).Where(lr.UserID.Eq(uid)).Delete()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.DeleteByUserID")
		l.log.WithContext(ctx).Error("login_log delete error", err)
		return err
	}
	return nil
}

func (l *LoginLogRepo) DeleteByUserIDs(ctx context.Context, uids []uint64) error {
	lr := l.baseRepo.LoginLog
	_, err := lr.WithContext(ctx).Where(lr.UserID.In(uids...)).Delete()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.DeleteByUserIDs")
		l.log.WithContext(ctx).Error("login_log delete error", err)
		return err
	}
	return nil
}
