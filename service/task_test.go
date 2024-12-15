// Package service
// Date: 2024/12/15 14:39:41
// Author: Amu
// Description:
package service

import "testing"

func TestExecute(t *testing.T) {
	conf := Config{
		StoragePath: "/etc/amcert/storage.db",
	}
	t.Logf("conf: %#v", conf)
	task := NewTimedTask(&conf)
	t.Logf("task: %#v", task)
	task.Execute()
}
