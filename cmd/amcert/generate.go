// Package main
// Date       : 2024/9/2 19:43
// Author     : Amu
// Description:
package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
)

func runGenerate() error {
	db.Init("/etc/amcert/storage.db")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter contact email: ")
	Email, _ := reader.ReadString('\n')
	Email = strings.TrimSpace(Email)

	fmt.Print("Please enter contact certificate path: ")
	Path, _ := reader.ReadString('\n')
	Path = strings.TrimSpace(Path)

	fmt.Print("Please enter contact domains(split by comma): ")
	Domains, _ := reader.ReadString('\n')
	Domains = strings.TrimSpace(Domains)

	fmt.Printf("Email: %s, Path: %s, Domains: %s\n", Email, Path, Domains)

	// generate certificate
	domains := strings.Split(Domains, ",")
	fmt.Printf("domains: %#v\n", domains)
	certificate := cert.NewCertificate(&cert.Config{
		RenewBefore:   cert.RenewBefore,
		CheckInterval: cert.CheckInterval,
		ContactEmail:  Email,
		CacheDir:      Path,
		Domains:       domains,
	})

	fmt.Printf("certificate: %#v\n", certificate)
	err := certificate.Generate()
	if err != nil {
		slog.Error("generate certificate", "error", err)
		return err
	}

	// save certificate info
	certConfig := cert.Config{
		ContactEmail: Email,
		CacheDir:     Path,
		Domains:      domains,
	}
	key := fmt.Sprintf("cert-%s", time.Now().Format("20060102150405"))
	err = db.PutJson(key, certConfig)
	if err != nil {
		slog.Error("Failed to put certificate into database", "error", err)
		return err
	}
	return nil
}
