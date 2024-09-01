// Package service
// Date       : 2024/8/30 18:19
// Author     : Amu
// Description:
package service

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Log                Log                 `yaml:"log"`
	CertificateConfigs []CertificateConfig `yaml:"certificates"`
}

func NewConfig(configFile string) (*Config, error) {
	config := &Config{}

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

type Log struct {
	Output   string `yaml:"output"`
	Level    string `yaml:"level"`
	Rotation int    `yaml:"rotation"`
	MaxAge   int    `yaml:"max_age"`
}

type CertificateConfig struct {
	RenewBefore   int      `yaml:"renew_before"`
	CheckInterval int      `yaml:"check_interval"`
	ContactEmail  string   `yaml:"contact_email"`
	Cert          string   `yaml:"cert"`
	Domains       []string `yaml:"domains"`
}
