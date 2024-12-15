// Package main
// Date: 2024/09/19 11:20:46
// Author: Amu
// Description:
package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
)

func runDB(args []string) error {
	err := db.Initialize("/etc/amcert/storage.db")
	if err != nil {
		slog.Error("db initialize failed", "error", err)
		return err
	}

	subCommand := args[0]
	switch subCommand {
	case "keys":
		keys, err := db.GetPrefixKeys("cert-")
		if err != nil {
			slog.Error("get prefix keys failed", "error", err)
			return err
		}
		fmt.Printf("All keys: %#v\n", keys)
		return nil
	case "detail":
		key := args[1]
		var conf cert.Config
		if err := db.GetJson(key, &conf); err != nil {
			slog.Error("get certificate config failed", "error", err)
			return err
		}
		conf.RenewBefore = cert.RenewBefore
		conf.CheckInterval = time.Duration(cert.CheckInterval) * time.Hour
		fmt.Printf("Certificate config: %#v\n", conf)
		return nil
	case "expire":
		key := args[1]
		var conf cert.Config
		if err := db.GetJson(key, &conf); err != nil {
			slog.Error("get certificate config failed", "error", err)
			return err
		}
		conf.RenewBefore = cert.RenewBefore
		conf.CheckInterval = time.Duration(cert.CheckInterval) * time.Hour
		certificate := cert.NewCertificate(&conf)
		if err := certificate.Load(); err != nil {
			slog.Error("load certificate failed", "error", err)
			return err
		}
		expire, err := certificate.Expire()
		if err != nil {
			slog.Error("get certificate expire failed", "error", err)
			return err
		}
		fmt.Printf("Certificate %s expire: %d\n", certificate.Domain, expire)
		return nil
	}
	return nil
}
