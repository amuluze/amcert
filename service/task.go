// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
	"github.com/amuluze/amutool/timex"
)

type CertConfig struct {
	Email   string   `json:"email"`
	Path    string   `json:"path"`
	Domains []string `json:"domains"`
}

type TimedTask struct {
	ticker      timex.Ticker
	stopCh      chan struct{}
	storagePath string
}

func NewTimedTask(config *Config) *TimedTask {
	tk := timex.NewTicker(cert.CheckInterval * time.Hour)
	slog.Info("new timed task", "conf", config.StoragePath)
	return &TimedTask{
		ticker:      tk,
		stopCh:      make(chan struct{}),
		storagePath: config.StoragePath,
	}
}

func (r *TimedTask) Execute() {
	slog.Info("timed task execute", "storage path", r.storagePath)
	err := db.Initialize(r.storagePath)
	if err != nil {
		slog.Error("db initialize failed", "error", err)
		return
	}
	keys, err := db.GetPrefixKeys("cert-")
	if err != nil {
		slog.Error("Get prefix keys error", "error", err)
		return
	}
	slog.Info("cert keys", "keys", keys)
	for _, key := range keys {
		var conf cert.Config
		err := db.GetJson(key, &conf)
		if err != nil {
			slog.Error("Get json error", "err", err)
			continue
		}
		conf.RenewBefore = cert.RenewBefore
		conf.CheckInterval = time.Duration(cert.CheckInterval) * time.Hour
		certificate := cert.NewCertificate(&conf)
		slog.Info("cert config", "config", conf)
		err = certificate.Load()
		if err != nil {
			slog.Error("Load certificate error", "err", err)
			continue
		}

		err = certificate.Renew()
		if err != nil {
			slog.Error("Renew certificate error", "err", err)
			continue
		}
	}
}

func (r *TimedTask) Run() {
	slog.Info("timed task starting...")
	for {
		select {
		case <-r.ticker.Chan():
			go r.Execute()
		case <-r.stopCh:
			fmt.Println("task exit")
			return
		}
	}
}

func (r *TimedTask) Stop() {
	close(r.stopCh)
}
