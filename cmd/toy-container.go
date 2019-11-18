package main

import (
	"github.com/urfave/cli"
	"os"
)



func init() {
	println("toy-container init")
}

func main() {

	toyContainer := cli.NewApp()
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
