// Package main
// Date       : 2024/8/30 18:24
// Author     : Amu
// Description:
package main

import "github.com/amuluze/amcert/pkg/config"

// runSetup 生成 amcert 配置文件
func runSetup() error {
	cfg, err := config.Create(configFile)
	if err != nil {
		return err
	}
	if err := cfg.Save(); err != nil {
		return err
	}
	return nil
}
