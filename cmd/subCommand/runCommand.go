package subCommand

import (
	"github.com/spf13/cobra"
	"toy-container/config"
)
/*
	xdushepherd 2019/11/15 14:15
	run方法需要做的事情：
    1. 解析配置信息
    2. 根据配置信息的rootfs路径构建overlay2文件系统用于容器运行时
    3. 启动进程，进行命名空间的设置
    4. 根据配置信息获得最终的应用程序启动命令，启动应用程序
 */


var containerID string

func init() {
	RunCommand.Flags().StringVarP(&containerID,"container-id","d","default-id","容器id")
	runConfig := &config.Config{}
	runConfig.Id = containerID
}

var RunCommand = &cobra.Command{
	Use:                        "run [containerID]",
	Short:                      "运行一个容器",
	Long:                       "根据配置文件运行一个容器",
	Run: func(cmd *cobra.Command, args []string) {

	}}
