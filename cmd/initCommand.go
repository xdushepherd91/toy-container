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

/*
	xdushepherd 2019/11/19 18:33
	1. 读取pid记录文件，获取有toy-container启动的进程id
    2. 通过pid，将进程的命名空间存下来，以便后续程序使用
 */

func mountNamespace() error {

	////从指定目录获取pid
	//pidArray := util.GetPidFromFile("./run/default-id/child_pid.txt")
	//pid := util.ParsePid(pidArray)
	//
	//


	return nil
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
