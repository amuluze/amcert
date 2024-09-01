// Package service
// Date       : 2024/8/30 17:29
// Author     : Amu
// Description:
package service

import "log/slog"

func Run(configFile string) (func(), error) {
	injector, clearFunc, err := BuildInjector(configFile)
	if err != nil {
		slog.Error("build injector failed", "err", err)
		return nil, err
	}
	
	// 初始化日志
	slog.SetDefault(injector.Logger.Logger)
	
	// 定时任务
	//renewTask := injector.RenewTask
	//go renewTask.Run()
	
	return func() {
		clearFunc()
	}, nil
}
