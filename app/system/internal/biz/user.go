package biz

import (
	"context"
	"github.com/spf13/cast"
	"golang.org/x/crypto/openpgp/errors"
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
	log          *log.Helper
}

func NewUserUsecase(userRepo UserRepo, userRoleRepo UserRoleRepo, roleRepo RoleRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		log:          log.NewHelper(log.With(logger, "biz", "user")),
	}
}
func (u *UserUsecase) buildUserReply(ctx context.Context, upo *model.User) (user *pb.UserReply, err error) {
	if err = copier.Copy(user, upo); err != nil {
		return
	}
	urs, err := u.userRoleRepo.SelectAllByUserID(ctx, uint64(upo.ID))
	if err != nil {
		return
	}
	var rids []uint
	for _, ur := range urs {
		rids = append(rids, uint(ur.RoleID))
	}
	roles, err := u.roleRepo.SelectByIDs(ctx, rids)
	if err != nil {
		return
	}
	var userRoles []*pb.UserRole
	for _, role := range roles {
		userRoles = append(userRoles, &pb.UserRole{
			Id:          strconv.Itoa(int(role.ID)),
			Name:        role.Name,
			Description: role.Description,
		})
	}
	user.UserRoles = userRoles
	user.CreatedAt = upo.CreatedAt.Format(timeFormat)
	user.UpdatedAt = upo.UpdatedAt.Format(timeFormat)
	user.Gender = upo.Gender.String()
	return
}
func (u *UserUsecase) validateUser(ctx context.Context, user *pb.UserRequest) (err error) {
	if u.userRepo.ExistByUserName(ctx, user.Username) {
		err = api.ErrorUserAlreadyExists("用户名已存在")
		return
	}
	if u.userRepo.ExistByEmail(ctx, user.Email) {
		err = api.ErrorEmailAlreadyExists("邮箱已存在")
		return
	}
	if u.userRepo.ExistByMobile(ctx, user.Mobile) {
		err = api.ErrorMobileAlreadyExists("手机号已存在")
		return
	}
	return
}

func (u *UserUsecase) Save(ctx context.Context, user *pb.UserRequest) (id string, err error) {
	var upo *model.User
	if err = u.validateUser(ctx, user); err != nil {
		return
	}
	if err = copier.Copy(&upo, user); err != nil {
		return
	}
	upo.Password = cypher.BcryptMake(user.Password)
	err = u.userRepo.Insert(ctx, upo)
	if err != nil {
		return
	}
	id = upo.GetUid()
	uid := upo.ID
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
			return
		}
	}

	return upo.GetUid(), nil
}

func (u *UserUsecase) Edit(ctx context.Context, user *pb.UserRequest) (id string, err error) {
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
	if len(urs) > 0 {
		err = u.userRoleRepo.UpdateByUserID(ctx, uid, urs)
		if err != nil {
			return
		}
	}
	var upo *model.User
	if err = u.validateUser(ctx, user); err != nil {
		return
	}
	if err = copier.Copy(upo, user); err != nil {
		return
	}
	err = u.userRepo.Update(ctx, upo)
	id = user.Id
	return
}

func (u UserUsecase) Remove(ctx context.Context, uids []uint) (err error) {
	if len(uids) == 0 {
		return errors.InvalidArgumentError("uids is null")
	}
	err = u.userRoleRepo.DeleteByUserIDs(ctx, xcast.ToUint64Slice[uint](uids))
	if err != nil {
		return
	}
	return u.userRepo.DeleteByIDs(ctx, uids)
}

func (u *UserUsecase) Get(ctx context.Context, uid uint) (user *pb.UserReply, err error) {
	upo, err := u.userRepo.Select(ctx, uid)
	if err != nil {
		return
	}
	return u.buildUserReply(ctx, upo)
}

func (u *UserUsecase) Search(ctx context.Context, limit, offset int, opt SQLOption) (list *pb.ListUserReply, err error) {
	users, total, err := u.userRepo.List(ctx, limit, offset, opt)
	if err != nil {
		return
	}
	for _, user := range users {
		var userReply *pb.UserReply
		userReply, err = u.buildUserReply(ctx, user)
		if err != nil {
			return
		}
		list.Users = append(list.Users, userReply)
	}
	list.Total = uint32(total)
	return
}

func (u *UserUsecase) EditStatus(ctx context.Context, uid uint, status model.UserStatus) error {
	return u.userRepo.UpdateStatus(ctx, uid, status)
}

func (u *UserUsecase) EditPassword(ctx context.Context, oldpwd, newpwd string, uid uint) error {
	user, err := u.userRepo.Select(ctx, uid)
	if err != nil {
		return err
	}
	if !cypher.BcryptCheck(oldpwd, user.Password) {
		return api.ErrorPasswordInvalid("密码错误")
	}
	user.Password = cypher.BcryptMake(newpwd)
	return u.userRepo.Update(ctx, user)
}
