//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"github.com/kyson/e-shop-native/internal/user-srv/data"
	"github.com/kyson/e-shop-native/internal/user-srv/service"
	"github.com/kyson/e-shop-native/internal/user-srv/sever"
	"github.com/kyson/e-shop-native/internal/user-srv/biz"
)

func provideServerConfig(c *conf.Bootstrap) *conf.Server {
	return c.Server
}

func provideDataConfig(c *conf.Bootstrap) *conf.Data {
	return c.Data
}



func InitializeApp() (*App, func(), error) {
	panic(wire.Build(
		provideDataConfig,
		provideServerConfig,

		NewApp,
		LoadConfig,
		biz.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		sever.ProviderSet,
	))
}