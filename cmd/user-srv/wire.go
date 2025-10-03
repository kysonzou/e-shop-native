//go:build wireinject

package main

import (
	"github.com/google/wire"
	//"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"github.com/kyson/e-shop-native/internal/user-srv/data"
	"github.com/kyson/e-shop-native/internal/user-srv/service"
	"github.com/kyson/e-shop-native/internal/user-srv/sever"
	"github.com/kyson/e-shop-native/internal/user-srv/biz"
)

func InitializeApp() (*App, func(), error) {
	panic(wire.Build(
		ProvideDataConfig,
		ProvideServerConfig,
		ProvideLogConfig,

		LoadConfig,
		NewApp,
		NewLogger,
		biz.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		sever.ProviderSet,
	))
}