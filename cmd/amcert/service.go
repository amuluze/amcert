// Package main
// Date       : 2024/8/30 17:36
// Author     : Amu
// Description:
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amuluze/amcert/service"
	"github.com/takama/daemon"
)

type Service struct {
	configFile string
	daemon     daemon.Daemon
}

// Start amcert bootstrap service non-blocking
func (s *Service) Start() {
	fmt.Printf("Starting amcert bootstrap service...\n")
}

// Run start amcert bootstrap service blocking and wait for exit signal
func (s *Service) Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	clearFunc, err := service.Run(s.configFile)
	if err != nil {
		return
	}

	for range interrupt {
		clearFunc()
		return
	}
}

// Stop amcert bootstrap service
func (s *Service) Stop() {
	fmt.Println("Stopping amcert bootstrap service...")
}

func (s *Service) manager(args []string) (string, error) {
	if len(args) > 0 {
		command := args[0]
		switch command {
		case "install":
			installArgs := args[1:]
			return s.daemon.Install(installArgs...)
		case "remove":
			return s.daemon.Remove()
		case "start":
			return s.daemon.Start()
		case "stop":
			return s.daemon.Stop()
		case "status":
			return s.daemon.Status()
		case "setup":
			fmt.Printf("initializing configuration files...\n")
			if err := runSetup(); err != nil {
				fmt.Printf("error initializing configuration files: %v\n", err)
				os.Exit(-1)
			} else {
				os.Exit(0)
			}
		case "generate":
			fmt.Printf("generating ssl certificate...\n")
			if err := runGenerate(args[1:]); err != nil {
				fmt.Printf("error generating ssl certificate: %v\n", err)
				os.Exit(-1)
			} else {
				os.Exit(0)
			}
		case "ssl":
			fmt.Printf("run db command...\n")
			if err := runDB(args[1:]); err != nil {
				fmt.Printf("error running db command: %v\n", err)
				os.Exit(-1)
			} else {
				os.Exit(0)
			}
		default:
			usage()
			return "", nil
		}
	}

	return s.daemon.Run(s)
}
