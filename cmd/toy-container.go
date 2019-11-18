package main

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	toyContainer = &cobra.Command{
		Use:   "toy-contianer",
		Short: "一个用来创建玩具容器的命令行工具",
		Long: "toy-container是一个仿照runc的轮子，可以讲一组程序运行在特定的namespace和cgroup之中，实现资源的隔离",
	}
)

func init() {
	println("toy-container init")
}

func main() {

	toyContainer.AddCommand(initCommand)
	toyContainer.AddCommand(runCommand)
	println("toy-container enter")
	if err := toyContainer.Execute(); err != nil {
		println("toy-container命令执行出错")
		os.Exit(1)
	}


}
