// Package cert
// Date       : 2024/8/13 18:24
// Author     : Amu
// Description:
package cert

const (
	RenewBefore               = 30 * 24
	CheckInterval             = 24
	DefaultContactEmail       = "wangjialong89@yeah.net"
	UserFileName              = "user.json"
	DomainFileName            = "domain"
	CertificateFileName       = "fullchain.pem"
	PrivateKeyFileName        = "privkey.pem"
	IssuerCertificateFileName = "issuer.cert.pem"
	CSRFileName               = "csr.pem"
	LocalCertificateFileName  = "localcert.pem"
	LocalPrivateKeyFileName   = "localkey.pem"
	SecretID                  = "TENCENTCLOUD_SECRET_ID"
	SecretKey                 = "TENCENTCLOUD_SECRET_KEY"
)
