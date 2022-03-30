package repo

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewLoginLogRepo,
	NewUserRepo,
	NewUserRoleRepo,
	NewRoleRepo,
	NewRoleResourceRepo,
	NewResourceMenuRepo,
	NewResourceActionRepo,
	NewResourceRouterRepo,
	NewCaptchaRepo,
)
