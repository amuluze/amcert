// Package main
// Date       : 2024/9/2 19:43
// Author     : Amu
// Description:
package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
)

func runGenerate(args []string) error {
	err := db.Initialize("/etc/amcert/storage.db")
	if err != nil {
		slog.Error("db initialize failed", "error", err)
		return err
	}

	subCommand := args[0]
	switch subCommand {
	case "ssl":
		return generateSSL()
	case "tls":
		return generateTLS()
	default:
		return fmt.Errorf("unknown subcommand: %s", subCommand)
	}
}

func generateSSL() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter contact email: ")
	Email, _ := reader.ReadString('\n')
	Email = strings.TrimSpace(Email)

	fmt.Print("Please enter certificate save path: ")
	Path, _ := reader.ReadString('\n')
	Path = strings.TrimSpace(Path)

	fmt.Print("Please enter domains(split by comma): ")
	Domains, _ := reader.ReadString('\n')
	Domains = strings.TrimSpace(Domains)

	fmt.Println(" *************  ssl certificate generate info  ************* ")
	fmt.Printf("Email: %s, Path: %s, Domains: %s\n", Email, Path, Domains)
	fmt.Println(" *********************************************************** ")

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

	// ensure cache dir exists
	if err := os.MkdirAll(filepath.Dir(Path), os.ModePerm); err != nil {
		fmt.Printf("Error creating cache dir failed: %v\n", err)
		return err
	}

	if err := certificate.Generate(); err != nil {
		fmt.Printf("generate certificate failed: %v\n", err)
		return err
	}

	// save certificate generate info
	key := fmt.Sprintf("cert-%s", certificate.Domain)
	if err := db.PutJson(key, conf); err != nil {
		fmt.Printf("Failed to put certificate into database: %v\n", err)
		return err
	}
	return nil
}

func generateTLS() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter common name: ")
	CommonName, _ := reader.ReadString('\n')
	CommonName = strings.TrimSpace(CommonName)

	fmt.Print("Please enter cert and key storage path: ")
	CertPath, _ := reader.ReadString('\n')
	CertPath = strings.TrimSpace(CertPath)

	fmt.Printf(" *************  tls certificate generate info  ************* ")
	fmt.Printf("Common Name: %s, Cert Path: %s\n", CommonName, CertPath)
	fmt.Printf(" *********************************************************** ")

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Printf("Failed to generate private key: %s\n", err)
		return err
	}
	tpl := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: CommonName},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	derCert, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Printf("Failed to create tls certificate: %s\n", err)
		return err
	}

	certBuf := bytes.NewBuffer(nil)
	if err := pem.Encode(certBuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derCert,
	}); err != nil {
		fmt.Printf("Failed to encode tls certificate: %s\n", err)
		return err
	}

	keyBuf := bytes.NewBuffer(nil)
	if err := pem.Encode(keyBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		fmt.Printf("Failed to encode tls private key: %s\n", err)
		return err
	}
	pemCert := certBuf.Bytes()
	if err := os.WriteFile(filepath.Join(CertPath, "server.pem"), pemCert, 0600); err != nil {
		fmt.Printf("Failed to write pem cert: %s\n", err)
		return err
	}
	pemKey := keyBuf.Bytes()
	if err := os.WriteFile(filepath.Join(CertPath, "server.key"), pemKey, 0600); err != nil {
		fmt.Printf("Failed to write pem key: %s\n", err)
		return err
	}
	return nil
}
