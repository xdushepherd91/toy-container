package main

import (
	"github.com/urfave/cli"
	"os"
)

var (
	toyContainer = cli.NewApp()
)

func init() {
	println("toy-container init")
}

func main() {

	toyContainer.Usage = "toy-container"
	toyContainer.Name = "toy-container"

	toyContainer.Commands =[]cli.Command{
		initCommand,
		runCommand,
	}


	if err := toyContainer.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
