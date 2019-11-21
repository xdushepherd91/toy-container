package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"github.com/vishvananda/netlink/nl"
	"io"
	"os"
	"os/exec"
	"toy-container/message"
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
		overlay2.SetUpFS(&runConfig)

		// 初始化管道，用于toy-container run 进程与 toy-container init 进程之间进行通信

		parentInitPipe,childInitPipe,err := util.NewSocketPair("init")

		// xdushepherd 2019/11/18 14:11 为int程序设定管道环境变量，以便可以进行命名空间的初始化工作
		toyInit := exec.Command("/proc/self/exe", "init")

		// xdushepherd 2019/11/21 10:14	将init进程需要使用的管道复制给toyInit

		toyInit.ExtraFiles = append(toyInit.ExtraFiles,childInitPipe)

		toyInit.Env = append(toyInit.Env,fmt.Sprintf("_LIBCONTAINER_INITPIPE=%d",stdFdsCount+len(toyInit.ExtraFiles)-1))

		toyInit.Env = append(toyInit.Env,fmt.Sprintf("_TOYCONTAINER_EXEC=%d",0))

		if err := toyInit.Start() ; err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		nr := nl.NewNetlinkRequest(int(message.InitMsg), 0)

		//xdushepherd 2019/11/21 13:01 添加cloneFlags参数
		nr.AddData(&message.Int32msg{
			Type:  message.CloneFlagsAttr,
			Value: uint32(runConfig.CloneFlags),
		})

		//xdushepherd 2019/11/21 14:07 添加pid路径，用于init进程存储容器进程数据
		nr.AddData(&message.Bytemsg{
			Type:  message.PidPath,
			Value: []byte(runConfig.PidPath),
		})


		//xdushepherd 2019/11/21 14:09 添加容器运行根目录数据，用于init进程进行数据挂以及chroot
		nr.AddData(&message.Bytemsg{
			Type:  message.ContainerPath,
			Value: []byte(runConfig.ContainerPath),
		})

		//xdushepherd 2019/11/21 14:26	添加容器id，用于init进程设置容器hostname

		nr.AddData(&message.Bytemsg{
			Type:  message.Id,
			Value: []byte(runConfig.Id),
		})


		fp := bytes.NewReader(nr.Serialize())

		if _,err := io.Copy(parentInitPipe,fp); err !=nil{
			fmt.Println(err)
			os.Exit(1)
		}

     	if err := toyInit.Wait() ; err !=nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

