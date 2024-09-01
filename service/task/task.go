// Package task
// Date       : 2024/8/30 17:39
// Author     : Amu
// Description:
package task

import (
	"github.com/amuluze/amcert/pkg/cert"
	"github.com/patrickmn/go-cache"
	"time"
)

type ITask interface {
}

type Task struct {
	cert  *cert.Certificate
	cache *cache.Cache
}

func NewTask(cert *cert.Certificate) *Task {
	return &Task{
		cert:  cert,
		cache: cache.New(5*time.Minute, 60*time.Second),
	}
}
