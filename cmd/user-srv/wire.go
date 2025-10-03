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



func InitializeApp(s *conf.Server, d *conf.Data) (*App, func(), error) {
	panic(wire.Build(
		NewApp,
		biz.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		sever.ProviderSet,
	))
}