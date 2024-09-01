// Package service
// Date: 2022/11/9 10:18
// Author: Amu
// Description:
package service

import "testing"

func TestConfig(t *testing.T) {
	config, err := NewConfig("/Users/amu/Desktop/github/amcert/config/config.yml")
	if err != nil {
		t.Fatalf("new config failed: %#v", err)
	}
	t.Logf("config: %#v", config)
}
