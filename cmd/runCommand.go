package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"toy-container/config"
	"toy-container/overlay2"
)
/*
	xdushepherd 2019/11/15 14:15
	run方法需要做的事情：
    1. 解析配置信息
    2. 根据配置信息的rootfs路径构建overlay2文件系统用于容器运行时
    3. 启动进程，进行命名空间的设置
    4. 根据配置信息获得最终的应用程序启动命令，启动应用程序
 */


var (
	containerID string
	runConfig config.Config
	)
func init() {
	Runcommand.Flags().StringVarP(&containerID,"container-id","d","default-id","容器id")
	runConfig.Id = containerID
}

var Runcommand = &cobra.Command{
	Use:                        "run [containerID]",
	Short:                      "运行一个容器",
	Long:                       "根据配置文件运行一个容器",
	Run: func(cmd *cobra.Command, args []string) {
		// 程序开始运行
		println("toy-container run go ")

		// 使用overlay2制作联合文件系统，等同于docker使用aufs基于镜像制作容器rootfs
		if err := overlay2.SetUpFS(runConfig);err !=nil{
			println("overlay2 run error")
		}
		toyInit := exec.Command("/proc/self/exe", "init")
		if err := toyInit.Start() ; err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := toyInit.Wait() ; err !=nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}}

