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
	GetUid() uint
	GetUsername() string
	GetRoleIDs() []uint
	GetNickname() string
}

// Transaction 新增事务接口方法
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type UserRepo interface {
	Select(context.Context, uint) (*model.User, error)
	SelectByUsername(context.Context, string) (*model.User, error)
	SelectByMobile(context.Context, string) (*model.User, error)
	SelectByEmail(context.Context, string) (*model.User, error)
	SelectPasswordByUID(context.Context, uint) (*model.User, error)
	SelectPage(context.Context, int, int, *SQLOption) ([]*model.User, int64, error)
	SelectTokens(context.Context, uint) ([]string, error)
	Insert(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	DeleteByIDs(context.Context, []uint) error
	Verify(context.Context, uint) bool
}

type UserRoleRepo interface {
	SelectAllByUserID(context.Context, uint) ([]*model.UserRole, error)
	Insert(context.Context, ...*model.UserRole) error
	UpdateByUserID(context.Context, uint, []*model.UserRole) error
	DeleteByUserIDs(context.Context, []uint) error
	DeleteByRoleIDs(context.Context, []uint) error
}

type UserLogRepo interface {
	Select(context.Context, uint) (*model.UserLog, error)
	SelectLocation(context.Context, string) (*Location, error)
	SelectAgent(context.Context, string) (*Agent, error)
	SelectPageByUID(context.Context, uint, int, int) ([]*model.UserLog, int64, error)
	Insert(context.Context, *model.UserLog) error
	Delete(context.Context, uint) error
	DeleteByUserIDs(context.Context, []uint) error
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
	SelectByRoleID(context.Context, uint, ...model.ResourceType) ([]*model.RoleResource, error)
	SelectByResourceID(context.Context, string, ...model.ResourceType) ([]*model.RoleResource, error)
	Insert(context.Context, ...*model.RoleResource) error
	UpdateByRoleID(context.Context, uint, []*model.RoleResource) error
	DeleteByRoleIDs(context.Context, []uint) error
	DeleteByResourceIDs(context.Context, []string, model.ResourceType) error
}

type ResourceMenuRepo interface {
	Select(context.Context, uint) (*model.ResourceMenu, error)
	SelectAll(context.Context) ([]*model.ResourceMenu, error)
	SelectByIDs(context.Context, []uint) ([]*model.ResourceMenu, error)
	Insert(context.Context, ...*model.ResourceMenu) error
	Update(context.Context, *model.ResourceMenu) error
	DeleteByIDs(context.Context, []uint) error
}

type ResourceRouterRepo interface {
	SelectAll(context.Context) ([]*model.Router, error)
	SelectByRoleID(context.Context, string) ([]*model.ResourceRouter, error)
	Update(context.Context, []*model.ResourceRouter) error
	ClearByRoleIDs(context.Context, ...string) error
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
