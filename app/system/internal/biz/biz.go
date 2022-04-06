package biz

import (
	"context"
	"github.com/google/wire"
	"kratosx-fashion/app/system/internal/data/model"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewPublicUsecase,
	NewUserUsecase,
	NewRoleUsecase,
	NewResourceUsecase,
)

type JwtUser interface {
	GetUid() string
	GetUsername() string
	GetRoleIDs() []string
	GetNickname() string
}

// Transaction 新增事务接口方法
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type UserRepo interface {
	Select(context.Context, uint) (*model.User, error)
	SelectByUsername(context.Context, string) (*model.User, error)
	SelectPasswordByName(context.Context, string) (*model.User, error)
	SelectPasswordByMobile(context.Context, string) (*model.User, error)
	SelectPasswordByEmail(context.Context, string) (*model.User, error)
	SelectPage(context.Context, int, int, *SQLOption) ([]*model.User, int64, error)
	Insert(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	UpdateStatus(context.Context, uint, model.UserStatus) error
	DeleteByIDs(context.Context, []uint) error
	ExistByUserName(context.Context, string) bool
	ExistByEmail(context.Context, string) bool
	ExistByMobile(context.Context, string) bool
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
}

type LoginLogRepo interface {
	Select(context.Context, uint) (*model.LoginLog, error)
	SelectLocation(context.Context, string) (*Location, error)
	SelectAgent(context.Context, string) (*Agent, error)
	SelectPageByUserID(context.Context, uint64, int, int) ([]*model.LoginLog, int64, error)
	Insert(context.Context, *model.LoginLog) error
	Delete(context.Context, uint) error
	DeleteByUserIDs(context.Context, []uint64) error
}

type RoleRepo interface {
	Select(context.Context, uint) (*model.Role, error)
	SelectByIDs(context.Context, []uint) ([]*model.Role, error)
	SelectPage(context.Context, int, int) ([]*model.Role, int64, error)
	Insert(context.Context, ...*model.Role) error
	Update(context.Context, *model.Role) error
	DeleteByIDs(context.Context, []uint) error
}

type RoleResourceRepo interface {
	Select(context.Context, uint) (*model.RoleResource, error)
	SelectByRoleID(context.Context, uint64, ...model.ResourceType) ([]*model.RoleResource, error)
	SelectByResourceID(context.Context, uint64, ...model.ResourceType) ([]*model.RoleResource, error)
	Insert(context.Context, ...*model.RoleResource) error
	Update(context.Context, *model.RoleResource) error
	UpdateByRoleID(context.Context, uint64, []*model.RoleResource) error
	Delete(context.Context, uint) error
	DeleteByRoleIDs(context.Context, []uint64) error
	DeleteByResourceIDs(context.Context, []uint64, model.ResourceType) error
}

type ResourceMenuRepo interface {
	Select(context.Context, uint) (*model.ResourceMenu, error)
	SelectAll(context.Context) ([]*model.ResourceMenu, error)
	SelectPage(context.Context, int, int) ([]*model.ResourceMenu, int64, error)
	SelectByIDs(context.Context, []uint) ([]*model.ResourceMenu, error)
	SelectPageByIDs(context.Context, []uint, int, int) ([]*model.ResourceMenu, int64, error)
	Insert(context.Context, ...*model.ResourceMenu) error
	Update(context.Context, *model.ResourceMenu) error
	DeleteByIDs(context.Context, []uint) error
}

type ResourceRouterRepo interface {
	SelectAll(context.Context) ([]model.Router, error)
	SelectByRoleIDs(context.Context, []string) ([]model.ResourceRouter, error)
	Update(context.Context, []model.ResourceRouter) error
	ClearByRoleIDs(context.Context, []string) error
}

type CaptchaRepo interface {
	Create(context.Context) (Captcha, error)
	Verify(context.Context, Captcha) bool
}

type JwtRepo interface {
	Create(context.Context, JwtUser) (*Token, error)
	IsInBlackList(context.Context, string) bool
	JoinInBlackList(context.Context, string) error
	ParseToken(context.Context, string) (*model.CustomClaims, error)
	GetSecretKey() string
	GetIssuer() string
}
