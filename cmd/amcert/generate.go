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

	fmt.Printf(" *************  certificate generate info  ************* ")
	fmt.Printf("Email: %s, Path: %s, Domains: %s\n", Email, Path, Domains)
	fmt.Printf(" ******************************************************* ")

	// generate certificate
	domains := strings.Split(Domains, ",")
	conf := cert.Config{
		RenewBefore:   cert.RenewBefore,
		CheckInterval: cert.CheckInterval,
		ContactEmail:  Email,
		CacheDir:      Path,
		Domains:       domains,
	}
	certificate := cert.NewCertificate(&conf)

	err := certificate.Generate()
	if err != nil {
		slog.Error("generate certificate", "error", err)
		return err
	}

	// save certificate generate info
	key := fmt.Sprintf("cert-%s", certificate.Domain)
	err = db.PutJson(key, conf)
	if err != nil {
		slog.Error("Failed to put certificate into database", "error", err)
		return err
	}
	return nil
}
