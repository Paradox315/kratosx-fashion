package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	api "kratosx-fashion/api/system/v1"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
	"kratosx-fashion/pkg/ctxutil"
	"kratosx-fashion/pkg/cypher"
	"kratosx-fashion/pkg/xsync"
	"time"
)

type UserUsecase struct {
	jwtRepo      JwtRepo
	userRepo     UserRepo
	logRepo      UserLogRepo
	userRoleRepo UserRoleRepo
	roleRepo     RoleRepo
	tx           Transaction
	log          *log.Helper
	rdb          *redis.Client
}

func NewUserUsecase(jwtRepo JwtRepo, userRepo UserRepo, logRepo UserLogRepo, userRoleRepo UserRoleRepo, roleRepo RoleRepo, rdb *redis.Client, tx Transaction, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		jwtRepo:      jwtRepo,
		userRepo:     userRepo,
		logRepo:      logRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		tx:           tx,
		log:          log.NewHelper(log.With(logger, "biz", "user")),
		rdb:          rdb,
	}
}

func (u *UserUsecase) buildUserDto(ctx context.Context, upo *model.User) (user *pb.UserReply, err error) {
	urs, err := u.userRoleRepo.SelectAllByUserID(ctx, upo.ID)
	if err != nil {
		u.log.WithContext(ctx).Errorf("select user role error %s", err.Error())
		err = api.ErrorUserFetchFail("获取用户角色失败")
		return
	}
	var rids []uint
	for _, ur := range urs {
		rids = append(rids, ur.RoleID)
	}
	roles, err := u.roleRepo.SelectByIDs(ctx, rids)
	if err != nil {
		return
	}
	user = &pb.UserReply{
		Id:        uint64(upo.ID),
		Username:  upo.Username,
		Email:     upo.Email,
		Mobile:    upo.Mobile,
		Age:       uint32(upo.Age),
		Avatar:    upo.Avatar,
		Nickname:  upo.Nickname,
		Gender:    upo.Gender.String(),
		Creator:   upo.Creator,
		Status:    upo.Status == model.UserStatusNormal,
		CreatedAt: upo.CreatedAt.Format(timeFormat),
		UpdatedAt: upo.UpdatedAt.Format(timeFormat),
	}
	extras := upo.GetExtras()
	if extras != nil {
		user.Address = extras.Address
		user.Country = extras.Country
		user.City = extras.City
		user.Birthday = extras.Birthday
		user.Description = extras.Description
	}
	for _, role := range roles {
		user.Roles = append(user.Roles, &pb.UserRole{
			Id:   uint64(role.ID),
			Name: role.Name,
		})
	}
	return
}

func (u *UserUsecase) buildUserPo(ctx context.Context, user *pb.UserRequest) (upo *model.User) {
	upo = &model.User{
		Username: user.Username,
		Password: cypher.BcryptMake("123456"),
		Email:    user.Email,
		Mobile:   user.Mobile,
		Age:      uint8(user.Age),
		Avatar:   user.Avatar,
		Nickname: user.Nickname,
		Gender:   model.GenderStatus(user.Gender),
		Creator:  ctxutil.GetUsername(ctx),
		Status:   lo.If(user.Status, model.UserStatusNormal).Else(model.UserStatusForbid),
	}
	upo.ID = uint(user.Id)
	upo.SetExtras(&model.UserExtra{
		Address:     user.Address,
		Country:     user.Country,
		City:        user.City,
		Birthday:    user.Birthday,
		Description: user.Description,
	})
	return
}

func (u *UserUsecase) Save(ctx context.Context, user *pb.UserRequest) (id uint, err error) {
	upo := u.buildUserPo(ctx, user)
	if err = u.userRepo.Insert(ctx, upo); err != nil {
		u.log.WithContext(ctx).Errorf("save user role failed: %s", err.Error())
		err = api.ErrorUserAddFail("添加用户失败")
		return
	}
	var urs []*model.UserRole
	for _, ridStr := range user.Roles {
		rid := cast.ToUint(ridStr)
		if upo.ID == 0 || rid == 0 {
			continue
		}
		urs = append(urs, &model.UserRole{
			UserID: upo.ID,
			RoleID: rid,
		})
	}
	if len(urs) > 0 {
		err = u.userRoleRepo.Insert(ctx, urs...)
		if err != nil {
			u.log.WithContext(ctx).Errorf("save user role failed: %s", err.Error())
			err = api.ErrorUserAddFail("添加用户角色关系失败")
			return
		}
	}
	return upo.ID, nil
}

func (u *UserUsecase) Edit(ctx context.Context, user *pb.UserRequest) (err error) {
	var urs []*model.UserRole
	for _, rid := range user.Roles {
		if user.Id == 0 || rid == 0 {
			continue
		}
		urs = append(urs, &model.UserRole{
			UserID: uint(user.Id),
			RoleID: uint(rid),
		})
	}
	upo := u.buildUserPo(ctx, user)
	upo.Password = ""
	err = u.tx.ExecTx(ctx, func(ctx context.Context) error {
		if len(urs) > 0 {
			err = u.userRoleRepo.UpdateByUserID(ctx, uint(user.Id), urs)
			if err != nil {
				return err
			}
		}
		return u.userRepo.Update(ctx, upo)
	})
	if err != nil {
		u.log.WithContext(ctx).Errorf("edit user failed: %s", err.Error())
		err = api.ErrorUserUpdateFail("修改用户失败")
		return
	}
	return
}

func (u *UserUsecase) Remove(ctx context.Context, uids []uint) (err error) {
	if len(uids) == 0 {
		return
	}
	err = u.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = u.userRoleRepo.DeleteByUserIDs(ctx, uids)
		if err != nil {
			return err
		}
		return u.userRepo.DeleteByIDs(ctx, uids)
	})
	if err != nil {
		u.log.WithContext(ctx).Errorf("remove user failed: %s", err.Error())
		err = api.ErrorUserDeleteFail("删除用户失败")
		return
	}
	return
}

func (u *UserUsecase) Get(ctx context.Context, uid uint) (user *pb.UserReply, err error) {
	upo, err := u.userRepo.Select(ctx, uid)
	if err != nil {
		u.log.WithContext(ctx).Errorf("get user failed: %s", err.Error())
		err = api.ErrorUserFetchFail("用户不存在")
		return
	}
	return u.buildUserDto(ctx, upo)
}

func (u *UserUsecase) Search(ctx context.Context, limit, offset int, opt *SQLOption) (list *pb.ListUserReply, err error) {
	users, total, err := u.userRepo.SelectPage(ctx, limit, offset, opt)
	if err != nil {
		u.log.WithContext(ctx).Errorf("search user failed %s", err.Error())
		err = api.ErrorUserFetchFail("获取用户列表失败")
		return
	}
	list = &pb.ListUserReply{}
	for _, upo := range users {
		var user *pb.UserReply
		user, err = u.buildUserDto(ctx, upo)
		if err != nil {
			return
		}
		list.List = append(list.List, user)
	}
	list.Total = uint32(total)
	return
}

func (u *UserUsecase) EditPassword(ctx context.Context, pwdReq *pb.PasswordRequest) error {
	user, err := u.userRepo.SelectPasswordByUID(ctx, uint(pwdReq.Id))
	if err != nil {
		return err
	}
	if !cypher.BcryptCheck(pwdReq.OldPassword, user.Password) {
		return api.ErrorPasswordInvalid("密码错误")
	}
	if pwdReq.ConfirmPassword != pwdReq.NewPassword {
		return api.ErrorPasswordNotMatch("密码不匹配")
	}
	user.Password = cypher.BcryptMake(pwdReq.NewPassword)
	lock := xsync.Lock("edit_password_lock", int64(time.Second*2), u.rdb)
	if lock.Get() {
		defer lock.Release()
		if err = u.userRepo.Update(ctx, user); err != nil {
			u.log.WithContext(ctx).Errorf("edit password failed: %s", err.Error())
			err = api.ErrorUpdatePasswordFail("修改密码失败")
			return err
		}
		var tokens []string
		tokens, err = u.userRepo.SelectTokens(ctx, uint(pwdReq.Id))
		if err != nil {
			u.log.WithContext(ctx).Errorf("edit password failed: %s", err.Error())
			err = api.ErrorUpdatePasswordFail("修改密码失败")
			return err
		}
		for _, token := range tokens {
			if err = u.jwtRepo.JoinInBlackList(ctx, token); err != nil {
				u.log.WithContext(ctx).Errorf("edit password failed: %s", err.Error())
				err = api.ErrorUpdatePasswordFail("修改密码失败")
				return err
			}
		}
	}
	return nil
}

func (u *UserUsecase) ResetPassword(ctx context.Context, id uint) error {
	user, err := u.userRepo.SelectPasswordByUID(ctx, id)
	if err != nil {
		return err
	}
	user.Password = cypher.BcryptMake("123456")
	lock := xsync.Lock("reset_password_lock", int64(time.Second*2), u.rdb)
	if lock.Get() {
		defer lock.Release()
		if err = u.userRepo.Update(ctx, user); err != nil {
			u.log.WithContext(ctx).Errorf("reset password failed: %s", err.Error())
			err = api.ErrorUserUpdateFail("重置密码失败")
			return err
		}
		var tokens []string
		tokens, err = u.userRepo.SelectTokens(ctx, id)
		if err != nil {
			u.log.WithContext(ctx).Errorf("reset password failed: %s", err.Error())
			err = api.ErrorResetPasswordFail("重置密码失败")
			return err
		}
		for _, token := range tokens {
			if err = u.jwtRepo.JoinInBlackList(ctx, token); err != nil {
				u.log.WithContext(ctx).Errorf("reset password failed: %s", err.Error())
				err = api.ErrorResetPasswordFail("重置密码失败")
				return err
			}
		}
	}

	return nil
}

func (u *UserUsecase) LogPage(ctx context.Context, uid uint, limit, offset int) (list *pb.ListUserLogReply, err error) {
	logs, total, err := u.logRepo.SelectPageByUID(ctx, uid, limit, offset)
	if err != nil {
		u.log.WithContext(ctx).Errorf("search user log failed %s", err.Error())
		err = api.ErrorLogFetchFail("获取用户日志列表失败")
		return
	}
	list = &pb.ListUserLogReply{}
	for _, l := range logs {
		list.List = append(list.List, &pb.UserLog{
			Id:      uint64(l.ID),
			Ip:      l.Ip,
			Method:  l.Method,
			Path:    l.Path,
			Status:  l.Status,
			Country: l.Country,
			Region:  l.Region,
			City:    l.City,
			Position: func() *pb.Position {
				pos := l.GetPosition()
				if pos == nil {
					return nil
				}
				return &pb.Position{
					Lat: pos["lat"],
					Lng: pos["lng"],
				}
			}(),
			Time:       l.CreatedAt.Format(timeFormat),
			UserAgent:  l.UserAgent,
			Client:     l.Client,
			Os:         l.OS,
			Device:     l.Device,
			DeviceType: l.DeviceType.String(),
			Type:       l.Type,
		})
	}
	list.Total = uint32(total)
	return
}
