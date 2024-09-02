// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
package service

import (
	"encoding/json"
	"fmt"
	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/pkg/db"
	"github.com/amuluze/amutool/timex"
	"log/slog"
	"time"
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
	keys, err := db.GetPrefixKeys("cert-")
	if err != nil {
		return
	}
	for _, key := range keys {
		jsonString, err := db.GetJson(key)
		if err != nil {
			slog.Error("Get json error:", err)
			continue
		}
		var certConfig CertConfig
		err = json.Unmarshal([]byte(jsonString), &certConfig)
		if err != nil {
			slog.Error("Unmarshal json error:", err)
			continue
		}
		conf := &cert.Config{
			RenewBefore:   cert.RenewBefore,
			CheckInterval: cert.CheckInterval,
			CacheDir:      certConfig.Path,
			ContactEmail:  certConfig.Email,
			Domains:       certConfig.Domains,
		}
		certificate := cert.NewCertificate(conf)
		err = certificate.Load()
		if err != nil {
			slog.Error("Load certificate error:", err)
			continue
		}
		err = certificate.Renew()
		if err != nil {
			slog.Error("Renew certificate error:", err)
			continue
		}
	}
	return
}

func (r *TimedTask) Run() {
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
