package infra

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewDiscovery,
	NewRegistrar,
	NewLogger,
	NewStorage,
)
