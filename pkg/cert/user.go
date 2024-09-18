// Package cert
// Date       : 2024/8/13 18:18
// Author     : Amu
// Description:
package cert

import (
	"crypto"
	"crypto/rsa"

	"github.com/go-acme/lego/v4/registration"
)

type User struct {
	Email        string
	PrivateKey   *rsa.PrivateKey
	Registration *registration.Resource
}

func (c *User) GetEmail() string {
	return c.Email
}

func (c *User) GetPrivateKey() crypto.PrivateKey {
	return c.PrivateKey
}

func (c *User) GetRegistration() *registration.Resource {
	return c.Registration
}
