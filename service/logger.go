// Package service
// Date       : 2024/8/30 18:17
// Author     : Amu
// Description:
package service

import (
	"github.com/amuluze/amcert/pkg/config"
	"github.com/amuluze/amutool/logger"
)

func NewLogger(config *config.Config) *logger.Logger {
	logx := logger.NewJsonFileLogger(
		logger.SetLogFile(config.Log.Output),
		logger.SetLogLevel(config.Log.Level),
		logger.SetLogFileRotationTime(config.Log.Rotation),
		logger.SetLogFileMaxAge(config.Log.MaxAge),
	)
	return logx
}
