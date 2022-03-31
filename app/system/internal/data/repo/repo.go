package repo

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewLoginLogRepo,
	NewUserRepo,
	NewUserRoleRepo,
	NewRoleRepo,
	NewRoleResourceRepo,
	NewResourceMenuRepo,
	NewResourceRouterRepo,
	NewCaptchaRepo,
	NewJwtRepo,
)
