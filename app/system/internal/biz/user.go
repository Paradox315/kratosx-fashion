package biz

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"kratosx-fashion/app/system/internal/data/model"
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
	roleRepo     RoleRepo
	tx           Transaction
	log          *log.Helper
}

func NewUserUsecase(userRepo UserRepo, userRoleRepo UserRoleRepo, roleRepo RoleRepo, tx Transaction, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		tx:           tx,
		log:          log.NewHelper(log.With(logger, "biz", "user")),
	}
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
		user.UserRoles = append(user.UserRoles, UserRole{
			ID:          strconv.FormatUint(uint64(role.ID), 10),
			Name:        role.Name,
			Description: role.Description,
		})
	}
	user.Id = cast.ToString(upo.ID)
	user.CreatedAt = upo.CreatedAt.Format(timeFormat)
	user.UpdatedAt = upo.UpdatedAt.Format(timeFormat)
	user.Gender = upo.Gender.String()
	return
}
func (u *UserUsecase) validateUser(ctx context.Context, user *pb.UserRequest) (err error) {
	if len(user.Username) != 0 && u.userRepo.ExistByUsername(ctx, user.Username) {
		err = api.ErrorUserAlreadyExists("用户名已存在")
		return
	}
	if len(user.Email) != 0 && u.userRepo.ExistByEmail(ctx, user.Email) {
		err = api.ErrorEmailAlreadyExists("邮箱已存在")
		return
	}
	if len(user.Mobile) != 0 && u.userRepo.ExistByMobile(ctx, user.Mobile) {
		err = api.ErrorMobileAlreadyExists("手机号已存在")
		return
	}
	return
}

func (u *UserUsecase) Save(ctx context.Context, user *pb.UserRequest) (id string, err error) {
	upo := &model.User{}
	if err = u.validateUser(ctx, user); err != nil {
		return
	}
	_ = copier.Copy(&upo, &user)
	upo.Password = cypher.BcryptMake(user.Password)
	if err = u.userRepo.Insert(ctx, upo); err != nil {
		return
	}
	uid := upo.ID
	id = cast.ToString(uid)
	var urs []*model.UserRole
	for _, ur := range user.UserRoles {
		rid := cast.ToUint64(ur.RoleId)
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
	for _, ur := range user.UserRoles {
		rid, _ := strconv.ParseUint(ur.RoleId, 10, 64)
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

func (u *UserUsecase) EditPassword(ctx context.Context, oldpwd, newpwd string, uid uint) error {
	user, err := u.userRepo.SelectPasswordByUID(ctx, uid)
	if err != nil {
		err = errors.Wrap(err, "userUsecase.EditPassword.Select")
		u.log.WithContext(ctx).Error(err)
		return err
	}
	if !cypher.BcryptCheck(oldpwd, user.Password) {
		return api.ErrorPasswordInvalid("密码错误")
	}
	user.Password = cypher.BcryptMake(newpwd)
	return u.userRepo.Update(ctx, user)
}
