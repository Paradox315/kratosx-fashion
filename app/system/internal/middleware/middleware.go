package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewJwtService,
	NewCasbinAuth,
	NewGlobalMiddleware,
)
