package main

import (
	"github.com/urfave/cli"
	"os"
)


func main() {

	println("main init")
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
