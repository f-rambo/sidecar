//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/f-rambo/cloud-copilot/sidecar/internal/biz"
	"github.com/f-rambo/cloud-copilot/sidecar/internal/conf"
	"github.com/f-rambo/cloud-copilot/sidecar/internal/interfaces"
	"github.com/f-rambo/cloud-copilot/sidecar/internal/server"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, biz.ProviderSet, interfaces.ProviderSet, newApp))
}
