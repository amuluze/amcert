// Package main
// Date: 2024/09/19 11:20:46
// Author: Amu
// Description:
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
)

func runDB(args []string) error {
	err := db.Initialize("/etc/amcert/storage.db")
	if err != nil {
		fmt.Printf("db initialize failed: %v", err)
		return err
	}

	subCommand := args[0]
	switch subCommand {
	case "keys":
		keys, err := db.GetPrefixKeys("cert-")
		if err != nil {
			fmt.Printf("get prefix keys failed: %v", err)
			return err
		}
		fmt.Printf("All keys: %#v\n", keys)
		return nil
	case "detail":
		key := args[1]
		var conf cert.Config
		if err := db.GetJson(key, &conf); err != nil {
			fmt.Printf("get certificate config failed: %v", err)
			return err
		}
		conf.RenewBefore = cert.RenewBefore
		conf.CheckInterval = time.Duration(cert.CheckInterval) * time.Hour
		fmt.Printf("Certificate config: %#v\n", conf)
		return nil
	case "delete":
		key := args[1]
		var conf cert.Config
		if err := db.GetJson(key, &conf); err != nil {
			fmt.Printf("get certificate config failed: %v", err)
			return err
		}
		fmt.Printf("Certificate config: %#v\n", conf)
		if err := db.DeleteKey(key); err != nil {
			fmt.Printf("delete key failed: %v", err)
			return err
		}
		if _, err := os.Stat(conf.CacheDir); os.IsNotExist(err) {
			fmt.Printf("%s does not exist", conf.CacheDir)
		} else {
			err := os.RemoveAll(conf.CacheDir)
			if err != nil {
				fmt.Printf("delete cache dir failed: %v", err)
				return err
			}
		}
	case "expire":
		key := args[1]
		var conf cert.Config
		if err := db.GetJson(key, &conf); err != nil {
			fmt.Printf("get certificate config failed: %v", err)
			return err
		}
		conf.RenewBefore = cert.RenewBefore
		conf.CheckInterval = time.Duration(cert.CheckInterval) * time.Hour
		certificate := cert.NewCertificate(&conf)
		if err := certificate.Load(); err != nil {
			fmt.Printf("load certificate failed: %v", err)
			return err
		}
		expire, err := certificate.Expire()
		if err != nil {
			fmt.Printf("get certificate expire failed: %v", err)
			return err
		}
		fmt.Printf("Certificate %s expire: %d\n", certificate.Domain, expire)
		return nil
	}
	return nil
}
