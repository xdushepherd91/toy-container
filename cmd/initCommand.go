package main


import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
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
	Action: func(context *cli.Context) {
		if err := os.MkdirAll("/home/node1/Desktop/"+strconv.Itoa(os.Getpid()),777) ;err !=nil{
			println("mkdir error")
			os.Exit(1)
		}

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


	applicationCommand , err := exec.LookPath("sh");
        if err != nil {
           panic(err)
        }
	args := []string{
		"/home/test.sh",
	}

	if err := syscall.Exec(applicationCommand, args, os.Environ()); err != nil {
		return err
	}

	return nil
}
