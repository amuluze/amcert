// Package service
// Date       : 2024/8/30 18:17
// Author     : Amu
// Description:
package service

import "github.com/amuluze/amutool/logger"

func NewLogger() *logger.Logger {
	return logger.NewJsonFileLogger(
		logger.SetLogFile("/etc/amcert/cert.log"),
		logger.SetLogLevel("info"),
	)
}
