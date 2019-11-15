package subCommand

import (
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var InitCommand = &cobra.Command{
	Use:                        "init",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll("/home/node1/Desktop/"+strconv.Itoa(os.Getpid()),777) ;err !=nil{
			println("mkdir error")
			os.Exit(1)
		}
	},
}
