// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
package service

import "github.com/amuluze/amcert/service/task"

func NewRenewTask(conf *Config) task.ITask {
	return task.NewRenewTask(conf)
}
