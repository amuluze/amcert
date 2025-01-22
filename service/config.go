// Package service
// Date       : 2024/8/30 18:19
// Author     : Amu
// Description:
package service

import (
	"log/slog"

	"github.com/amuluze/amcert/pkg/config"
	"github.com/spf13/viper"
)

func NewConfig(configFile string) (*config.Config, error) {
	config := &config.Config{}
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("read config error", "err", err)
		return nil, err
	}
	if err := viper.Unmarshal(config); err != nil {
		slog.Error("parse config error", "error", err)
		return nil, err
	}
	return config, nil
}
