// Package main
// Date       : 2024/9/2 19:43
// Author     : Amu
// Description:
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
	"log/slog"
	"os"
	"strings"
)

func runGenerate() error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("Please enter contact email: ")
	Email, _ := reader.ReadString('\n')
	
	fmt.Print("Please enter contact certificate path: ")
	Path, _ := reader.ReadString('\n')
	
	fmt.Print("Please enter contact domains(split by comma): ")
	Domains, _ := reader.ReadString('\n')
	fmt.Printf("Email：%s, Path: %s, Domains: %s", Email, Path, Domains)
	
	// generate certificate
	domains := strings.Split(Domains, ",")
	certificate := cert.NewCertificate(&cert.Config{
		RenewBefore:   cert.RenewBefore,
		CheckInterval: cert.CheckInterval,
		ContactEmail:  Email,
		CacheDir:      Path,
		Domains:       domains,
	})
	err := certificate.Generate()
	if err != nil {
		slog.Error("Error generating certificate", "error", err)
		return err
	}
	
	// save certificate info
	certConfig := cert.Config{
		ContactEmail: Email,
		CacheDir:     Path,
		Domains:      domains,
	}
	certString, _ := json.Marshal(certConfig)
	err = db.PutString("", string(certString))
	if err != nil {
		slog.Error("Failed to put certificate into database", "error", err)
		return err
	}
	return nil
}
