package repo

import (
	"context"
	iploc "github.com/ip2location/ip2location-go"
	ua "github.com/mileusna/useragent"
	"github.com/pkg/errors"
	"kratosx-fashion/app/system/biz"
	"kratosx-fashion/app/system/data"
	"kratosx-fashion/app/system/data/linq"
	"kratosx-fashion/app/system/data/model"
	"net"

	"github.com/go-kratos/kratos/v2/log"
)

type UserLogRepo struct {
	dao      *data.Data
	log      *log.Helper
	ipdb     *iploc.DB
	baseRepo *linq.Query
}

func isPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func NewLoginLogRepo(data *data.Data, logger log.Logger, ipdb *iploc.DB) biz.UserLogRepo {
	return &UserLogRepo{
		dao:      data,
		log:      log.NewHelper(log.With(logger, "repo", "user_log")),
		ipdb:     ipdb,
		baseRepo: linq.Use(data.DB),
	}
}

func (l *UserLogRepo) SelectLocation(ctx context.Context, ip string) (loc *biz.Location, err error) {
	if len(ip) == 0 || ip == "127.0.0.1" || ip == "localhost" || ip == "::1" || !isPublicIP(net.ParseIP(ip)) {
		loc = &biz.Location{
			Country: "本地",
			Region:  "本地",
			City:    "本地",
			Position: map[string]float32{
				"lat": 0.1,
				"lng": 0.1,
			},
		}
		return
	}
	result, err := l.ipdb.Get_all(ip)
	if err != nil {
		err = errors.Wrap(err, "ipdb.Get_all")
		l.log.Error("ip2location get_all error", err)
		return
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

func (l *UserLogRepo) SelectAgent(ctx context.Context, agentStr string) (agent *biz.Agent, err error) {
	result := ua.Parse(agentStr)
	agent = &biz.Agent{
		Name:   result.Name,
		OS:     result.OS + result.OSVersion,
		Device: result.Device,
	}
	if result.Desktop {
		agent.DeviceType = model.DevicetypePc
	} else if result.Mobile {
		agent.DeviceType = model.DevicetypeMobile
	} else if result.Tablet {
		agent.DeviceType = model.DevicetypePad
	} else if result.Bot {
		agent.DeviceType = model.DevicetypeBot
	}
	return
}

func (l *UserLogRepo) SelectPageByUID(ctx context.Context, id uint, limit, offset int) (logs []*model.UserLog, total int64, err error) {
	lr := l.baseRepo.UserLog
	logs, err = lr.WithContext(ctx).Where(lr.UID.Eq(id)).Limit(limit).Offset(offset).Order(lr.UpdatedAt.Desc()).Select(lr.ALL).Find()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.SelectPageByUserID")
		l.log.WithContext(ctx).Error("pagination.Find error", err)
	}
	total, err = lr.WithContext(ctx).Where(lr.UID.Eq(id)).Count()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.SelectPageByUserID")
		l.log.WithContext(ctx).Error("pagination.Count error", err)
		return
	}
	return
}

func (l *UserLogRepo) Insert(ctx context.Context, loginLog *model.UserLog) error {
	lr := l.baseRepo.UserLog
	err := lr.WithContext(ctx).Create(loginLog)
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.Insert")
		l.log.WithContext(ctx).Error("login_log insert error", err)
		return err
	}
	return nil
}

func (l *UserLogRepo) Delete(ctx context.Context, id uint) error {
	lr := l.baseRepo.UserLog
	_, err := lr.WithContext(ctx).Where(lr.ID.Eq(id)).Delete()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.Delete")
		l.log.WithContext(ctx).Error("login_log delete error", err)
		return err
	}
	return nil
}

func (l *UserLogRepo) DeleteByUserIDs(ctx context.Context, uids []uint) error {
	lr := l.baseRepo.UserLog
	_, err := lr.WithContext(ctx).Where(lr.UID.In(uids...)).Delete()
	if err != nil {
		err = errors.Wrap(err, "baseRepo.LoginLog.DeleteByUserIDs")
		l.log.WithContext(ctx).Error("login_log delete error", err)
		return err
	}
	return nil
}
