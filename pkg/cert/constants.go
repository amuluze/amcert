// Package cert
// Date       : 2024/8/13 18:24
// Author     : Amu
// Description:
package cert

const (
	RenewBefore               = 30 * 24
	CheckInterval             = 1
	DefaultContactEmail       = "example@example.net"
	UserFileName              = "user.json"
	DomainFileName            = "domain"
	CertificateFileName       = "fullchain.pem"
	PrivateKeyFileName        = "privkey.pem"
	IssuerCertificateFileName = "issuer.cert.pem"
	CSRFileName               = "csr.pem"
	SecretID                  = "TENCENTCLOUD_SECRET_ID"
	SecretKey                 = "TENCENTCLOUD_SECRET_KEY"
)
