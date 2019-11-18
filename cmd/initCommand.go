package main


import (
	"github.com/urfave/cli"
	"os"
	"strconv"
	"syscall"
	_ "toy-container/namespace"
)


func init(){
	println( "initCommand begin")
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `initialize the namespaces and launch the process (do not call it outside of runc)`,
	Action: func(context *cli.Context) error {
		if err := os.MkdirAll("/home/node1/Desktop/"+strconv.Itoa(os.Getpid()),777) ;err !=nil{
			println("mkdir error")
			os.Exit(1)
		}

		if err := initContainer(); err != nil {
			println("init container error")
		}

		return nil
	},
}



func initContainer() error {
	/*
		xdushepherd 2019/11/18 10:11
		下面命令备用

		config := config.Config{}
		applicationCommand := config.ApplicationCmd
	 */

	println("in init command")

	applicationCommand := "echo";
	args := []string{
		"$HOSTNAME",">>","/opt/test.txt",
	}

	if err := syscall.Exec(applicationCommand, args, os.Environ()); err != nil {
		return err
	}

	return nil
}
