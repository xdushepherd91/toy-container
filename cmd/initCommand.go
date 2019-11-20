package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"syscall"
	_ "toy-container/namespace"
)


var initCommand = cli.Command{
	Name:  "init",
	Usage: `initialize the namespaces and launch the process (do not call it outside of runc)`,
	Action: func(context *cli.Context) {
		if err := initContainer(); err != nil {
			fmt.Println(err)
		}

	},
}

func initContainer() error {
	/*
		xdushepherd 2019/11/18 10:11
		下面命令备用

		config := config.Config{}
		applicationCommand := config.ApplicationCmd
	 */
	println("init container")
	applicationCommand , err := exec.LookPath("sh");
        if err != nil {
           panic(err)
        }
	args := []string{
		"sh",
		"/home/test.sh",
		"&",
	}

	if err := syscall.Exec(applicationCommand, args[0:], os.Environ()); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
