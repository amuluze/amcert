// Package task
// Date       : 2024/8/30 17:50
// Author     : Amu
// Description:
package task

import (
	"fmt"
	"github.com/amuluze/amcert/service"
	"github.com/amuluze/amutool/timex"
	"github.com/patrickmn/go-cache"
	"time"
)

var _ ITask = (*RenewTask)(nil)

type RenewTask struct {
	ticker timex.Ticker
	cache  *cache.Cache
	stopCh chan struct{}
}

func NewRenewTask(conf *service.Config) *RenewTask {
	tk := timex.NewTicker(time.Duration(interval) * time.Hour)
	return &RenewTask{
		ticker: tk,
		cache:  cache.New(5*time.Minute, 60*time.Second),
		stopCh: make(chan struct{}),
	}
}

func (r *RenewTask) Execute() {
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
