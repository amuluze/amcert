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
	ticker timex.Ticker
	stopCh chan struct{}
}

func NewTimedTask() *TimedTask {
	tk := timex.NewTicker(24 * time.Hour)
	return &TimedTask{
		ticker: tk,
		stopCh: make(chan struct{}),
	}
}

func (r *TimedTask) Execute() {
	slog.Info("timed task execute")
	keys, err := db.GetPrefixKeys("cert-")
	if err != nil {
		return
	}
	for _, key := range keys {
		var conf cert.Config
		err := db.GetJson(key, &conf)
		if err != nil {
			slog.Error("Get json error", "err", err)
			continue
		}

		slog.Info("cert config", "config", conf)
		conf.RenewBefore = cert.RenewBefore
		conf.CheckInterval = cert.CheckInterval
		certificate := cert.NewCertificate(&conf)

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
