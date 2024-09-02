// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
package service

import (
	"github.com/amuluze/amutool/logger"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(NewInjector)

type Injector struct {
	Config *Config
	Logger *logger.Logger
	Task   *TimedTask
}

func NewInjector(config *Config, task *TimedTask, logx *logger.Logger) (*Injector, error) {
	return &Injector{
		Config: config,
		Logger: logx,
		Task:   task,
	}, nil
}
