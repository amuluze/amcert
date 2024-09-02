// Package config
// Date       : 2024/8/30 18:21
// Author     : Amu
// Description:
package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

func Create(filename string) (*Config, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer func(fp *os.File) {
		err := fp.Close()
		if err != nil {
			return
		}
	}(fp)
	
	cfg := new(Config)
	cfg.filename = filename
	if err := cfg.loadDefault(); err != nil {
		return nil, err
	}
	
	if data, err := io.ReadAll(fp); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

type Config struct {
	filename string
	// 配置项
	Log struct {
		Output   string `yaml:"output"`
		Level    string `yaml:"level"`
		Rotation int    `yaml:"rotation"`
		MaxAge   int    `yaml:"max_age"`
	} `yaml:"log"`
	StoragePath string `yaml:"storage_path"`
}

func (c *Config) loadDefault() error {
	c.Log.Output = "/etc/amcert/cert.log"
	c.Log.Level = "info"
	c.Log.Rotation = 1
	c.Log.MaxAge = 7
	c.StoragePath = "/etc/amcert/storage.db"
	
	return nil
}

func (c *Config) Save() error {
	if c.filename == "" {
		return errors.New("empty filename")
	}
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(c.filename, out, 0644); err != nil {
		return err
	}
	return nil
}
