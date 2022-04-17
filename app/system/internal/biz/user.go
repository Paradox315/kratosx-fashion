package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/ctxutil"
	"kratosx-fashion/pkg/cypher"
	"kratosx-fashion/pkg/xcast"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	api "kratosx-fashion/api/system/v1"
	pb "kratosx-fashion/api/system/v1"
)

type UserUsecase struct {
	userRepo     UserRepo
	userRoleRepo UserRoleRepo
	logRepo      LoginLogRepo
	roleRepo     RoleRepo
	tx           Transaction
	log          *log.Helper
}

func NewUserUsecase(userRepo UserRepo, userRoleRepo UserRoleRepo, roleRepo RoleRepo, logRepo LoginLogRepo, tx Transaction, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		logRepo:      logRepo,
		tx:           tx,
		log:          log.NewHelper(log.With(logger, "biz", "user")),
	}
}
func (u *UserUsecase) buildLog(ctx context.Context, lpo *model.LoginLog) (log *pb.LoginLog, err error) {
	log = &pb.LoginLog{
		Id:         cast.ToString(lpo.ID),
		Ip:         lpo.Ip,
		Time:       lpo.CreatedAt.Format(timeFormat),
		Agent:      lpo.Agent,
		Os:         lpo.OS,
		Device:     lpo.Device,
		DeviceType: uint32(lpo.DeviceType),
		LoginType:  uint32(lpo.LoginType),
	}
	var loc Location
	if err = encoding.GetCodec("json").Unmarshal([]byte(lpo.Location), &loc); err != nil {
		return
	}
	log.Country = loc.Country
	log.Region = loc.Region
	log.City = loc.City
	log.Position = &pb.LoginLog_Position{
		Lat: loc.Position["lat"],
		Lng: loc.Position["lng"],
	}
	return
}

func (u *UserUsecase) buildUserDo(ctx context.Context, upo *model.User) (user User, err error) {
	_ = copier.Copy(&user, &upo)
	urs, err := u.userRoleRepo.SelectAllByUserID(ctx, uint64(upo.ID))
	if err != nil {
		err = errors.Wrap(err, "userUsecase.buildUserDto.SelectAllByUserID")
		u.log.WithContext(ctx).Error(err)
		return
	}
	var rids []uint
	for _, ur := range urs {
		rids = append(rids, uint(ur.RoleID))
	}
	roles, err := u.roleRepo.SelectByIDs(ctx, rids)
	if err != nil {
		err = errors.Wrap(err, "userUsecase.buildUserDto.SelectByIDs")
		u.log.WithContext(ctx).Error(err)
		return
	}
	for _, role := range roles {
		user.Roles = append(user.Roles, UserRole{
			Id:   strconv.FormatUint(uint64(role.ID), 10),
			Name: role.Name,
		})
	}
	user.Id = cast.ToString(upo.ID)
	user.CreatedAt = upo.CreatedAt.Format(timeFormat)
	user.UpdatedAt = upo.UpdatedAt.Format(timeFormat)
	user.Gender = upo.Gender.String()
	return
}

func (u *UserUsecase) validateUser(ctx context.Context, user *pb.UserRequest) (err error) {
	var cnt int64
	if len(user.Username) > 0 {
		cnt, err = u.userRepo.ExistByUsername(ctx, user.Username)
		if err != nil {
			return
		}
		if len(user.Id) > 0 && cnt > 1 {
			err = api.ErrorUserAlreadyExists("用户名已存在")
			return
		}
		if len(user.Id) == 0 && cnt > 0 {
			err = api.ErrorUserAlreadyExists("用户名已存在")
			return
		}
	}

	if len(user.Email) > 0 {
		cnt, err = u.userRepo.ExistByEmail(ctx, user.Email)
		if err != nil {
			return
		}
		if len(user.Id) > 0 && cnt > 1 {
			err = api.ErrorUserAlreadyExists("邮箱已存在")
			return
		}
		if len(user.Id) == 0 && cnt > 0 {
			err = api.ErrorEmailAlreadyExists("邮箱已存在")
			return
		}
	}

	if len(user.Mobile) > 0 {
		cnt, err = u.userRepo.ExistByMobile(ctx, user.Mobile)
		if err != nil {
			return
		}
		if len(user.Id) > 0 && cnt > 1 {
			err = api.ErrorUserAlreadyExists("手机号已存在")
			return
		}
		if len(user.Id) == 0 && cnt > 0 {
			err = api.ErrorMobileAlreadyExists("手机号已存在")
			return
		}
	}
	return
}

func (u *UserUsecase) Save(ctx context.Context, user *pb.UserRequest) (id string, err error) {
	upo := &model.User{}
	if err = u.validateUser(ctx, user); err != nil {
		return
	}
	upo.Creator = ctxutil.GetUsername(ctx)
	_ = copier.Copy(&upo, &user)
	upo.Password = cypher.BcryptMake(user.Password)
	if err = u.userRepo.Insert(ctx, upo); err != nil {
		return
	}
	uid := upo.ID
	id = cast.ToString(uid)
	var urs []*model.UserRole
	for _, ridstr := range user.Roles {
		rid := cast.ToUint64(ridstr)
		if uid == 0 || rid == 0 {
			continue
		}
		urs = append(urs, &model.UserRole{
			UserID: uint64(uid),
			RoleID: rid,
		})
	}
	if len(urs) > 0 {
		err = u.userRoleRepo.Insert(ctx, urs...)
		if err != nil {
			err = errors.Wrap(err, "userUsecase.Save.Insert")
			u.log.WithContext(ctx).Error(err)
			return
		}
	}

	return
}

func (u *UserUsecase) Edit(ctx context.Context, user *pb.UserRequest) (id string, err error) {
	id = user.Id
	uid := cast.ToUint64(user.Id)
	var urs []*model.UserRole
	for _, ridstr := range user.Roles {
		rid := cast.ToUint64(ridstr)
		if uid == 0 || rid == 0 {
			continue
		}
		urs = append(urs, &model.UserRole{
			UserID: uid,
			RoleID: rid,
		})
	}
	upo := &model.User{}
	if err = u.validateUser(ctx, user); err != nil {
		return
	}
	_ = copier.Copy(&upo, &user)
	upo.ID = uint(uid)
	upo.Password = ""
	err = u.tx.ExecTx(ctx, func(ctx context.Context) error {
		if len(urs) > 0 {
			err = u.userRoleRepo.UpdateByUserID(ctx, uid, urs)
			if err != nil {
				err = errors.Wrap(err, "useUsecase.Edit.UpdateByUserID")
				u.log.WithContext(ctx).Error(err)
				return err
			}
		}
		return u.userRepo.Update(ctx, upo)
	})
	return
}

func (u *UserUsecase) Remove(ctx context.Context, uids []uint) (err error) {
	if len(uids) == 0 {
		return
	}
	return u.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = u.userRoleRepo.DeleteByUserIDs(ctx, xcast.ToUint64Slice(uids))
		if err != nil {
			return err
		}
		return u.userRepo.DeleteByIDs(ctx, uids)
	})
}

func (u *UserUsecase) Get(ctx context.Context, uid uint) (user User, err error) {
	upo, err := u.userRepo.Select(ctx, uid)
	if err != nil {
		err = errors.Wrap(err, "userUsecase.Get.Select")
		u.log.WithContext(ctx).Error(err)
		return
	}
	return u.buildUserDo(ctx, upo)
}

func (u *UserUsecase) Search(ctx context.Context, limit, offset int, opt *SQLOption) (list []User, total int64, err error) {
	users, total, err := u.userRepo.SelectPage(ctx, limit, offset, opt)
	if err != nil {
		return
	}
	for _, upo := range users {
		var user User
		user, err = u.buildUserDo(ctx, upo)
		if err != nil {
			return
		}
		list = append(list, user)
	}
	return
}

func (u *UserUsecase) EditStatus(ctx context.Context, uid uint, status model.UserStatus) error {
	return u.userRepo.UpdateStatus(ctx, uid, status)
}

func (u *UserUsecase) EditPassword(ctx context.Context, oldpwd, newpwd, confirmPwd string, uid uint) error {
	user, err := u.userRepo.SelectPasswordByUID(ctx, uid)
	if err != nil {
		err = errors.Wrap(err, "userUsecase.EditPassword.Select")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	if !cypher.BcryptCheck(oldpwd, user.Password) {
		return api.ErrorPasswordInvalid("密码错误")
	}
	if newpwd != confirmPwd {
		return api.ErrorPasswordNotMatch("密码不匹配")
	}
	user.Password = cypher.BcryptMake(newpwd)
	return u.userRepo.Update(ctx, user)
}

func (u *UserUsecase) LogPage(ctx context.Context, uid uint64, limit, offset int) (list *pb.ListLoginLogReply, err error) {
	logs, total, err := u.logRepo.SelectPageByUserID(ctx, uid, limit, offset)
	if err != nil {
		return
	}
	list = &pb.ListLoginLogReply{}
	for _, l := range logs {
		var logDto *pb.LoginLog
		logDto, err = u.buildLog(ctx, l)
		if err != nil {
			return
		}
		list.List = append(list.List, logDto)
	}
	list.Total = uint32(total)
	return
}
