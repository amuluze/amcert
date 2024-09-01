// Package task
// Date       : 2024/8/30 17:39
// Author     : Amu
// Description:
package task

type ITask interface {
	Execute()
	Run()
	Stop()
}
