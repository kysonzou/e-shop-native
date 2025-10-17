package validator

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewValidator)
