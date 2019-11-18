package subCommand


import (
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"syscall"
	_ "toy-container/namespace"
)


func init(){
	println( "initCommand begin")
}

var InitCommand = &cobra.Command{
	Use:                        "init",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll("/home/node1/Desktop/"+strconv.Itoa(os.Getpid()),777) ;err !=nil{
			println("mkdir error")
			os.Exit(1)
		}

		if err := initContainer(); err != nil {
			println("init container error")
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

	applicationCommand := "echo";
	args := []string{
		"$HOSTNAME",">>","/opt/test.txt",
	}

	if err := syscall.Exec(applicationCommand, args, os.Environ()); err != nil {
		return err
	}

	return nil
}
