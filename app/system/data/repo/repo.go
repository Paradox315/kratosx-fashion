package repo

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/google/wire"
)

var codec = encoding.GetCodec("msgpack")

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
