package registry

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewNacosRegistry,
)
