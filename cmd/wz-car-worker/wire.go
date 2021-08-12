// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"wz-car-worker/internal/biz"
	"wz-car-worker/internal/conf"
	"wz-car-worker/internal/data"
	"wz-car-worker/internal/protocol"
	"wz-car-worker/internal/server"
	"wz-car-worker/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

func initConsumer(*conf.Data, log.Logger) (*protocol.Dispatcher, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, protocol.ProviderSet))
}
