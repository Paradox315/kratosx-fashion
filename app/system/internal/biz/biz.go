package biz

import (
	"context"
	"kratosx-fashion/app/system/internal/data/linq"

	"kratosx-fashion/app/system/internal/data/model"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewPublicUsecase,
	NewUserUsecase,
	NewRoleUsecase,
	NewResourceUsecase,
)

type UserRepo interface {
	Select(context.Context, uint) (*model.User, error)
	SelectByUsername(context.Context, string) (*model.User, error)
	SelectPasswordByName(context.Context, string) (*model.User, error)
	SelectPasswordByMobile(context.Context, string) (*model.User, error)
	SelectPasswordByEmail(context.Context, string) (*model.User, error)
	List(context.Context, int, int, SQLOption) ([]*model.User, int64, error)
	Insert(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	UpdateStatus(context.Context, uint, model.UserStatus) error
	DeleteByIDs(context.Context, []uint) error
	ExistByUserName(context.Context, string) bool
	ExistByEmail(context.Context, string) bool
	ExistByMobile(context.Context, string) bool
	BaseRepo(ctx context.Context) *linq.Query
}

type UserRoleRepo interface {
	Select(context.Context, uint) (*model.UserRole, error)
	SelectAll(context.Context) ([]*model.UserRole, error)
	SelectAllByUserID(context.Context, uint64) ([]*model.UserRole, error)
	Insert(context.Context, ...*model.UserRole) error
	Update(context.Context, *model.UserRole) error
	UpdateByUserID(context.Context, uint64, []*model.UserRole) error
	Delete(context.Context, uint) error
	DeleteByUserIDs(context.Context, []uint64) error
	DeleteByRoleIDs(context.Context, []uint64) error
	Exist(context.Context, uint64, uint64) bool
}

type LoginLogRepo interface {
	Select(context.Context, uint) (*model.LoginLog, error)
	SelectLocation(context.Context, string) (*Location, error)
	SelectAgent(context.Context, string) (*Agent, error)
	ListByUserID(context.Context, uint64, int, int) ([]*model.LoginLog, int64, error)
	Insert(context.Context, *model.LoginLog) error
	Delete(context.Context, uint) error
	DeleteByUserIDs(context.Context, []uint64) error
}

type RoleRepo interface {
	Select(context.Context, uint) (*model.Role, error)
	SelectByIDs(context.Context, []uint) ([]*model.Role, error)
	List(context.Context, int, int) ([]*model.Role, int64, error)
	Insert(context.Context, ...*model.Role) error
	Update(context.Context, *model.Role) error
	DeleteByIDs(context.Context, []uint) error
	BaseRepo(ctx context.Context) *linq.Query
}

type RoleResourceRepo interface {
	Select(context.Context, uint) (*model.RoleResource, error)
	SelectByRoleID(context.Context, uint64, ...model.ResourceType) ([]*model.RoleResource, error)
	Insert(context.Context, ...*model.RoleResource) error
	Update(context.Context, *model.RoleResource) error
	UpdateByRoleID(context.Context, uint64, []*model.RoleResource) error
	Delete(context.Context, uint) error
	DeleteByRoleIDs(context.Context, []uint64) error
	DeleteByResourceIDs(context.Context, []uint64, model.ResourceType) error
	BaseRepo(ctx context.Context) *linq.Query
}

type ResourceMenuRepo interface {
	Select(context.Context, uint64) (*model.ResourceMenu, error)
	SelectByIDs(context.Context, []uint) ([]*model.ResourceMenu, error)
	Insert(context.Context, *model.ResourceMenu) error
	Update(context.Context, *model.ResourceMenu) error
	DeleteByIDs(context.Context, []uint) error
}

type ResourceActionRepo interface {
	Select(context.Context, uint64) (*model.ResourceAction, error)
	SelectByMenuID(context.Context, uint64) ([]*model.ResourceAction, error)
	Insert(context.Context, *model.ResourceAction) error
	Update(context.Context, *model.ResourceAction) error
	Delete(context.Context, uint64) error
	DeleteByMenuIDs(context.Context, []uint64) error
}

type ResourceRouterRepo interface {
	SelectAll(context.Context) ([]*model.ResourceRouter, error)
	SelectByRoleIDs(context.Context, []string) ([]*model.ResourceRouter, error)
	Insert(context.Context, *model.ResourceRouter) error
	Update(context.Context, *model.ResourceRouter) error
	DeleteByRoleIDs(context.Context, []string) error
	Exist(context.Context, string, string, string) (bool, error)
}

type CaptchaRepo interface {
	Create(context.Context) (Captcha, error)
	Verify(context.Context, Captcha) bool
}
