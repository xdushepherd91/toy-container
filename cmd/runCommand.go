package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"toy-container/overlay2"
	"toy-container/util"
)
/*
    xdushepherd 2019/11/15 14:15
    run方法需要做的事情：
    1. 解析配置信息
    2. 根据配置信息的rootfs路径构建overlay2文件系统用于容器运行时
    3. 启动进程，进行命名空间的设置
    4. 根据配置信息获得最终的应用程序启动命令，启动应用程序
 */



var runCommand = cli.Command{
	Name:  "run",
	Usage: `根据指定的容器id创建容器并运行`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "id,d",
			Usage:       "容器名称",
			Value:       "default-id",
		},
	},
	Action: func(context *cli.Context) {

		var id string
		if id = context.String("id");id=="" {
			fmt.Println("容器id不能为空")
		}
		// 从配置文件中读取配置信息
		runConfig,err := util.LoadConfig(configPath)

		if err != nil  {
			fmt.Println(err)
			os.Exit(1)

		}


		runConfig.Id = id
		if runConfig.MountSource == "" {
			runConfig.MountSource = mountSource
		}
		if runConfig.MountType == "" {
			runConfig.MountType = mountType
		}
		runConfig.ContainerPath = toyContainerRootPath+"/"+id


		// 使用overlay2制作联合文件系统，等同于docker使用aufs基于镜像制作容器rootfs
		if err := overlay2.SetUpFS(runConfig);err !=nil{
			println("overlay2 run error")
		}


		// xdushepherd 2019/11/18 14:11 为int程序设定管道环境变量，以便可以进行命名空间的初始化工作

		toyInit := exec.Command("/proc/self/exe", "init")

		toyInit.Env = append(toyInit.Env,fmt.Sprint("_LIBCONTAINER_INITPIPE=1"))

		if err := toyInit.Start() ; err != nil {
			fmt.Println(err)
			os.Exit(1)
		}


		println("init calling")

/*		if err := toyInit.Wait() ; err !=nil {
			fmt.Println(err)
			os.Exit(1)
		}

		*/
		println("everything is ok");
	},
}

