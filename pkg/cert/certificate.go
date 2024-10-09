// Package cert
// Date       : 2024/8/13 10:27
// Author     : Amu
// Description:
package cert

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	"github.com/go-acme/lego/v4/registration"
)

var _ ICertificate = (*Certificate)(nil)

type ICertificate interface {
	Load() error
	Generate() error
	Renew() error
	Expire() (int, error)
}

type Certificate struct {
	Domains               []string `json:"domains"`
	Domain                string   `json:"domain"`
	PrivateKey            []byte   `json:"private_key"`
	Certificate           []byte   `json:"certificate"`
	IssuerCertificate     []byte   `json:"issuer_certificate"`
	CSR                   []byte   `json:"csr"`
	DomainPath            string   `json:"domain_path"`
	PrivateKeyPath        string   `json:"private_key_path"`
	CertificatePath       string   `json:"certificate_path"`
	IssuerCertificatePath string   `json:"issuer_certificate_path"`
	CSRPath               string   `json:"csr_path"`
	UserPath              string   `json:"user_path"`
	RenewBefore           int
}

func NewCertificate(config *Config) *Certificate {
	cert := &Certificate{
		PrivateKeyPath:        filepath.Join(config.CacheDir, PrivateKeyFileName),
		CertificatePath:       filepath.Join(config.CacheDir, CertificateFileName),
		UserPath:              filepath.Join(config.CacheDir, UserFileName),
		DomainPath:            filepath.Join(config.CacheDir, DomainFileName),
		IssuerCertificatePath: filepath.Join(config.CacheDir, IssuerCertificateFileName),
		CSRPath:               filepath.Join(config.CacheDir, CSRFileName),
		Domains:               config.Domains,
		RenewBefore:           config.RenewBefore,
	}
	return cert
}

func (c *Certificate) Load() error {
	if _, err := FileExists(c.DomainPath); err != nil {
		slog.Error("Domain path not found", "domain path", c.DomainPath, "error", err)
		return err
	}
	domainData, err := os.ReadFile(c.DomainPath)
	if err != nil {
		slog.Error("Read domain file error", "domain path", c.DomainPath, "error", err)
		return err
	}
	c.Domain = strings.TrimSpace(string(domainData))

	if _, err := FileExists(c.PrivateKeyPath); err != nil {
		slog.Error("Private key path not found", "private key path", c.PrivateKeyPath, "error", err)
		return err
	}
	keyData, err := os.ReadFile(c.PrivateKeyPath)
	if err != nil {
		slog.Error("Read private key file error", "private key path", c.PrivateKeyPath, "error", err)
		return err
	}
	c.PrivateKey = keyData

	if _, err = FileExists(c.CertificatePath); err != nil {
		slog.Error("Certificate path not found", "certificate path", c.CertificatePath, "error", err)
		return err
	}
	certData, err := os.ReadFile(c.CertificatePath)
	if err != nil {
		slog.Error("Read certificate file error", "certificate path", c.CertificatePath, "error", err)
		return err
	}
	c.Certificate = certData

	if _, err := FileExists(c.IssuerCertificatePath); err != nil {
		slog.Error("Issuer certificate path not found", "issuer certificate path", c.IssuerCertificatePath, "error", err)
		return err
	}
	issuerData, err := os.ReadFile(c.IssuerCertificatePath)
	if err != nil {
		slog.Error("Read issuer certificate file error", "issuer certificate path", c.IssuerCertificatePath, "error", err)
		return err
	}
	c.IssuerCertificate = issuerData

	if _, err = FileExists(c.CSRPath); err != nil {
		slog.Error("CSR path not found", "csr path", c.CSRPath, "error", err)
		return err
	}
	csrData, err := os.ReadFile(c.CSRPath)
	if err != nil {
		slog.Error("Read csr file error", "csr path", c.CSRPath, "error", err)
		return err
	}
	c.CSR = csrData
	return nil
}

func (c *Certificate) Generate() error {
	u, err := c.getUser()
	if err != nil {
		slog.Error("Get user error", "error", err)
		return err
	}

	client, err := c.createClient(u)
	if err != nil {
		slog.Error("Create client error", "error", err)
		return err
	}
	request := certificate.ObtainRequest{
		Domains: c.Domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		slog.Error("Obtain certificate error", "error", err)
		return err
	}

	c.Domain = certificates.Domain
	c.Certificate = certificates.Certificate
	c.PrivateKey = certificates.PrivateKey
	c.CSR = certificates.CSR
	c.IssuerCertificate = certificates.IssuerCertificate
	// save certificate to file
	if err := c.save(); err != nil {
		slog.Error("Save certificate error", "error", err)
		return err
	}
	return nil
}

func (c *Certificate) Renew() error {
	expire, err := c.Expire()
	if err != nil {
		slog.Error("Get certificate expire error", "error", err)
		return err
	}

	slog.Info("certificate expire", "certificate domain", c.Domain, "expire", expire)
	if c.RenewBefore >= expire {
		u, err := c.getUser()
		if err != nil {
			slog.Error("Get user error", "error", err)
			return err
		}
		client, err := c.createClient(u)
		if err != nil {
			slog.Error("Create client error", "error", err)
			return err
		}
		resource, err := client.Certificate.RenewWithOptions(certificate.Resource{
			Domain:            c.Domain,
			CertURL:           lego.LEDirectoryProduction,
			CertStableURL:     lego.LEDirectoryStaging,
			PrivateKey:        c.PrivateKey,
			Certificate:       c.Certificate,
			IssuerCertificate: c.IssuerCertificate,
			CSR:               c.CSR,
		}, &certificate.RenewOptions{})
		if err != nil {
			slog.Error("Renew certificate error", "error", err)
			return err
		}
		c.CSR = resource.CSR
		c.IssuerCertificate = resource.IssuerCertificate
		c.PrivateKey = resource.PrivateKey
		c.Certificate = resource.Certificate
		c.Domain = resource.Domain

		err = c.save()
		if err != nil {
			slog.Error("Save renew certificate error", "error", err)
			return err
		}
	}
	return nil
}

/**
 * TODO： 需要完善和补充 DNSProviderConfig
 */
func (c *Certificate) createClient(u *User) (lego.Client, error) {
	config := lego.NewConfig(u)
	config.CADirURL = lego.LEDirectoryProduction
	config.Certificate.KeyType = certcrypto.EC256

	client, err := lego.NewClient(config)
	if err != nil {
		slog.Error("Create lego client error", "error", err)
		return lego.Client{}, err
	}

	cfg := tencentcloud.NewDefaultConfig()
	cfg.SecretID = GetSecretID()
	cfg.SecretKey = GetSecretKey()
	provider, err := tencentcloud.NewDNSProviderConfig(cfg)
	if err != nil {
		return lego.Client{}, fmt.Errorf("Create DNS provider error: %s", err)
	}
	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return lego.Client{}, fmt.Errorf("Setting DNS provider error: %s", err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return lego.Client{}, fmt.Errorf("Registration error: %v", err)
	}

	u.Registration = reg
	c.saveUser(u)
	return *client, nil
}

func (c *Certificate) save() error {
	if err := os.WriteFile(c.DomainPath, []byte(c.Domain), 0600); err != nil {
		return fmt.Errorf("error writing domain file: %v", err)
	}
	if err := os.WriteFile(c.PrivateKeyPath, c.PrivateKey, os.ModePerm); err != nil {
		return fmt.Errorf("error creating private key: %s", err)
	}
	if err := os.WriteFile(c.CertificatePath, c.Certificate, os.ModePerm); err != nil {
		return fmt.Errorf("error creating certificate: %s", err)
	}
	if err := os.WriteFile(c.IssuerCertificatePath, c.IssuerCertificate, 0600); err != nil {
		return fmt.Errorf("error creating issuer certificate: %s", err)
	}
	if err := os.WriteFile(c.CSRPath, c.CSR, os.ModePerm); err != nil {
		return fmt.Errorf("error creating CSR: %s", err)
	}
	return nil
}

func (c *Certificate) Expire() (int, error) {
	var certDERBlock *pem.Block
	certDERBlock, _ = pem.Decode(c.Certificate)
	if certDERBlock.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			return -1, err
		}
		timeLeft := time.Until(cert.NotAfter)
		return int(timeLeft.Hours()), nil
	}
	return -1, fmt.Errorf("invalid certificate")
}

func (c *Certificate) getUser() (*User, error) {
	var user User
	b, err := os.ReadFile(c.UserPath)
	if b != nil && err != nil {
		err = json.Unmarshal(b, &user)
		if err != nil {
			return nil, err
		}
	} else {
		privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.EC256)
		if err != nil {
			return &user, fmt.Errorf("simplecert: failed to generate private key: %s", err)
		}
		user.PrivateKey = privateKey
		user.Email = DefaultContactEmail
	}
	return &user, nil
}

func (c *Certificate) saveUser(user *User) {
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Printf("error marshalling user: %v\n", err)
		return
	}
	err = os.WriteFile(c.UserPath, b, os.ModePerm)
	if err != nil {
		fmt.Printf("error saving user: %v\n", err)
		return
	}
}
