// Package db
// Date: 2024/09/03 23:52:04
// Author: Amu
// Description:
package db

import (
	"testing"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDB(t *testing.T) {
	err := Initialize("./storage.db")
	if err != nil {
		t.Fatalf("init error: %s", err)
	}
	user := User{
		Name: "jhon",
		Age:  12,
	}
	err = PutJson("cert-12345", user)
	if err != nil {
		t.Fatalf("put json error: %s", err)
	}

	var u User
	err = GetJson("cert-12345", &u)
	if err != nil {
		t.Fatalf("get json error: %s", err)
	}
	t.Logf("u: %#v", u)
}
