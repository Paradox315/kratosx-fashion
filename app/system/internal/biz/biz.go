package biz

import (
	"context"

	"kratosx-fashion/app/system/internal/data/model"

	"github.com/google/wire"

	pb "kratosx-fashion/api/system/v1"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewPublicUsecase,
	NewUserUsecase,
	NewRoleUsecase,
	NewMenuUsecase,
)

type UserRepo interface {
	Select(context.Context, uint) (*model.User, error)
	SelectByUsername(context.Context, string) (*model.User, error)
	SelectPasswordByName(context.Context, string) (uint, string, error)
	List(context.Context, *pb.ListRequest, ...*pb.QueryOption) ([]*model.User, int64, error)
	Insert(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	UpdateStatus(context.Context, uint, uint8) error
	Delete(context.Context, uint) error
	DeleteByIDs(context.Context, []uint) error
	ExistByUserName(context.Context, string) bool
	ExistByEmail(context.Context, string) bool
	ExistByMobile(context.Context, string) bool
}

type UserRoleRepo interface {
	Select(context.Context, uint) (*model.UserRole, error)
	SelectAll(context.Context) ([]*model.UserRole, error)
	SelectAllByUserID(context.Context, uint64) ([]*model.UserRole, error)
	SelectRoleIDByUserID(context.Context, uint64) ([]uint64, error)
	Insert(context.Context, *model.UserRole) error
	Update(context.Context, *model.UserRole) error
	Delete(context.Context, uint) error
	DeleteByUserID(context.Context, uint64) error
	DeleteByUserIDs(context.Context, []uint64) error
	DeleteByRoleID(context.Context, uint64) error
	DeleteByRoleIDs(context.Context, []uint64) error
	ExistByUserIDAndRoleID(context.Context, uint64, uint64) (bool, error)
}

type LoginLogRepo interface {
	Select(context.Context, uint) (*model.LoginLog, error)
	ListByUserID(context.Context, uint64, *pb.ListRequest, ...*pb.QueryOption) ([]*model.LoginLog, int64, error)
	Insert(context.Context, *model.LoginLog) error
	Delete(context.Context, uint) error
	DeleteByUserID(context.Context, uint64) error
	DeleteByUserIDs(context.Context, []uint64) error
}

type RoleRepo interface {
	Select(context.Context, uint64) (*model.Role, error)
	SelectByIDs(context.Context, []uint64) ([]*model.Role, error)
	List(context.Context, *pb.ListRequest, ...*pb.QueryOption) ([]*model.Role, uint64, error)
	Insert(context.Context, *model.Role) error
	Update(context.Context, *model.Role) error
	UpdateStatus(context.Context, uint64, uint8) error
	Delete(context.Context, uint64) error
	DeleteByIDs(context.Context, []uint64) error
}

type RoleMenuRepo interface {
	Select(context.Context, uint64) (*model.RoleMenu, error)
	SelectByRoleID(context.Context, uint64) ([]*model.RoleMenu, error)
	ExistByRoleIDAndMenuID(context.Context, uint64, uint64) (bool, error)
	Insert(context.Context, *model.RoleMenu) error
	Update(context.Context, *model.RoleMenu) error
	Delete(context.Context, uint64) error
	DeleteByRoleID(context.Context, uint64) error
	DeleteByRoleIDs(context.Context, []uint64) error
	DeleteByMenuID(context.Context, uint64) error
	DeleteByMenuIDs(context.Context, []uint64) error
}

type MenuRepo interface {
	Select(context.Context, uint64) (*model.Menu, error)
	SelectByIDs(context.Context, []uint64) ([]*model.Menu, error)
	List(context.Context, *pb.ListRequest, ...*pb.QueryOption) ([]*model.Menu, uint64, error)
	Insert(context.Context, *model.Menu) error
	Update(context.Context, *model.Menu) error
	Delete(context.Context, uint64) error
	DeleteByIDs(context.Context, []uint64) error
}

type MenuActionRepo interface {
	Select(context.Context, uint64) (*model.MenuAction, error)
	SelectByMenuID(context.Context, uint64) ([]*model.MenuAction, error)
	ExistByMenuIDAndActionID(context.Context, uint64, uint64) (bool, error)
	Insert(context.Context, *model.MenuAction) error
	Update(context.Context, *model.MenuAction) error
	Delete(context.Context, uint64) error
	DeleteByMenuID(context.Context, uint64) error
	DeleteByMenuIDs(context.Context, []uint64) error
}

type MenuActionResourceRepo interface {
	Select(context.Context, uint64) (*model.MenuActionResource, error)
	SelectByActionID(context.Context, uint64) ([]*model.MenuActionResource, error)
	Insert(context.Context, *model.MenuActionResource) error
	Update(context.Context, *model.MenuActionResource) error
	Delete(context.Context, uint64) error
	DeleteByIDs(context.Context, []uint64) error
	DeleteByActionID(context.Context, uint64) error
	DeleteByActionIDs(context.Context, []uint64) error
}

type CaptchaRepo interface {
	Create(context.Context) (string, string, error)
	Verify(context.Context, string, string) bool
}
