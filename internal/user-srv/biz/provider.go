package biz

import "github.com/google/wire"

// ProviderSet is a provider set for non-test builds.
var ProviderSet = wire.NewSet(NewUserUsecase, NewBcrypt)
