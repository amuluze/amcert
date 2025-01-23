// Package main
// Date       : 2024/8/30 17:28
// Author     : Amu
// Description:
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/takama/daemon"
)

const (
	name        = "amcert"
	description = "certificate generate and renew"
)

var dependencies = []string{""}

var (
	configFile string
)

func usage() {
	fmt.Println("Description: \n\t", description)
	fmt.Println("Usage: \n\t", os.Args[0], " install | remove | start | stop | status | setup")
	fmt.Println("       \t", os.Args[0], " generate ssl")
	fmt.Println("       \t", os.Args[0], " ssl keys")
	fmt.Println("       \t", os.Args[0], " ssl detail key")
	fmt.Println("       \t", os.Args[0], " ssl expire key")
	fmt.Println("Flags: ")
}

func parseConfig() []string {
	flag.StringVar(&configFile, "conf", "/etc/amcert/config.yml", "config file path")
	flag.Parse()
	return flag.Args()
}

// TODO: need to opt output and log

func main() {
	flag.Usage = usage
	args := parseConfig()

	var kind daemon.Kind
	if runtime.GOOS == "darwin" {
		kind = daemon.GlobalDaemon
	} else if runtime.GOOS == "linux" {
		kind = daemon.SystemDaemon
	}

	src, err := daemon.New(name, description, kind, dependencies...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	service := &Service{daemon: src, configFile: configFile}
	status, err := service.manager(args)
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
