// Package main
// Date       : 2024/9/2 19:43
// Author     : Amu
// Description:
package main

import (
	"github.com/amuluze/amcert/pkg/db"
	"log/slog"
	"strings"
)

func runGenerate() error {
	configPath := strings.Trim(configFile, "/config.yml")
	err := db.Init(configPath)
	if err != nil {
		slog.Error("Error initializing database", "err", err)
		return err
	}
	return nil
}
