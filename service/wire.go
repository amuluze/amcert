// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
//go:build wireinject
// +build wireinject

package service

import (
	"github.com/google/wire"
)

func BuildInjector(configFile string) (*Injector, func(), error) {
	wire.Build(
		NewConfig,
		NewLogger,
		NewRenewTask,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
