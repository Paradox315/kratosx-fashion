package biz

import (
	"context"
	"github.com/google/wire"
	pb "kratosx-fashion/api/system/v1"
	"kratosx-fashion/app/system/internal/data/model"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewPublicUsecase,
	NewUserUsecase,
	NewRoleUsecase,
	NewMenuUsecase,
)

type UserRepo interface {
	SelectByID(context.Context, uint64) (*model.User, error)
	SelectByUsername(context.Context, string) (*model.User, error)
	SelectPasswordByID(context.Context, uint64) (string, error)
	List(context.Context, *pb.ListUserReply, ...*pb.QueryOption) ([]*model.User, uint64, error)
	Insert(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	Delete(context.Context, uint64) error
	DeleteByIDs(context.Context, []uint64) error
	ExistByUserName(context.Context, string) (bool, error)
}

type UserRoleRepo interface {
	Select(context.Context, uint64) (*model.UserRole, error)
	SelectAll(context.Context) ([]*model.UserRole, error)
	SelectAllByUserID(context.Context, uint64) ([]*model.UserRole, error)
	SelectRoleIDByUserID(context.Context, uint64) ([]uint64, error)
	Insert(context.Context, *model.UserRole) error
	Update(context.Context, *model.UserRole) error
	UpdateStatus(context.Context, uint64, uint8) error
	Delete(context.Context, uint64) error
	DeleteByUserID(context.Context, uint64) error
	DeleteByUserIDs(context.Context, []uint64) error
	DeleteByRoleID(context.Context, uint64) error
	DeleteByRoleIDs(context.Context, []uint64) error
	ExistByUserIDAndRoleID(context.Context, uint64, uint64) (bool, error)
}

type LoginLogRepo interface {
	Select(context.Context, uint64) (*model.LoginLog, error)
	ListByUserID(context.Context, uint64, *pb.ListRequest, ...*pb.QueryOption) ([]*model.LoginLog, error)
	Insert(context.Context, *model.LoginLog) error
	Delete(context.Context, uint64) error
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
