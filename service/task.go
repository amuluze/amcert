// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
package service

import (
	"fmt"
	"github.com/amuluze/amcert/pkg/cert"
	"github.com/amuluze/amcert/service/task"
	"github.com/amuluze/amutool/timex"
	"time"
)

type RenewTask struct {
	tasks  []*task.Task
	ticker timex.Ticker
	stopCh chan struct{}
}

func NewRenewTask(conf *Config) *RenewTask {
	tk := timex.NewTicker(time.Duration(conf.CertificateConfigs[0].CheckInterval) * time.Hour)
	
	var renewTasks []*task.Task
	for _, cc := range conf.CertificateConfigs {
		ce := cert.NewCertificate(&cert.Config{
			RenewBefore:   cc.RenewBefore,
			CheckInterval: time.Duration(cc.CheckInterval) * time.Hour,
			ContactEmail:  cc.ContactEmail,
			Domains:       cc.Domains,
			CacheDir:      cc.Cert,
		})
		renewTasks = append(renewTasks, task.NewTask(ce))
	}
	return &RenewTask{
		ticker: tk,
		tasks:  renewTasks,
		stopCh: make(chan struct{}),
	}
}

func (r *RenewTask) Execute() {
	for _, rt := range r.tasks {
		fmt.Printf("task: %#v\n", rt)
	}
	return
}

func (r *RenewTask) Run() {
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

func (r *RenewTask) Stop() {
	close(r.stopCh)
}
