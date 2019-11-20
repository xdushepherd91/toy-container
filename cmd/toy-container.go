package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"toy-container/util"
)

const (
	configPath = "config"
	toyContainerRootPath = "/var/lib/toyContainer/overlay2"
	mountType = "overlay2"
	mountSource = "overlay2"
)


func main() {

	if !util.IsDirOrFileExist(toyContainerRootPath) {
		err := os.MkdirAll(toyContainerRootPath,741)
		if err != nil {
			fmt.Println("创建容器运行时根目录失败")
			os.Exit(1)
		}
	}


	toyContainer := cli.NewApp()
	toyContainer.Usage = "toy-container"
	toyContainer.Name = "toy-container"

	toyContainer.Commands =[]cli.Command{
		initCommand,
		execCommand,
		runCommand,

	}

	if err := toyContainer.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
