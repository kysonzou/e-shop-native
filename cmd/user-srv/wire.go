//go:build wireinject

//go:generate wire //把下面go:build注释掉就可以使用，但是会冲突
package main

import (
	"github.com/google/wire"
	//"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"github.com/kyson/e-shop-native/internal/user-srv/auth"
	"github.com/kyson/e-shop-native/internal/user-srv/biz"
	"github.com/kyson/e-shop-native/internal/user-srv/data"
	"github.com/kyson/e-shop-native/internal/user-srv/service"
	"github.com/kyson/e-shop-native/internal/user-srv/server"
	"github.com/kyson/e-shop-native/internal/user-srv/validator"
)

func InitializeApp() (*App, func(), error) {
	panic(wire.Build(
		ProvideDataConfig,
		ProvideServerConfig,
		ProvideLogConfig,
		ProvideAuthConfig,

		LoadConfig,
		NewApp,
		NewLogger,
		biz.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		auth.ProviderSet,
		validator.ProviderSet,
	))
}
