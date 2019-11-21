package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"github.com/vishvananda/netlink/nl"
	"io"
	"os"
	"os/exec"
	"strings"
	"toy-container/message"

	//	"os"
	"io/ioutil"
	"strconv"
	"toy-container/util"
)

var execCommand = cli.Command{
	Name:  "exec",
	Usage: `重新进入一个容器命名空间，进行相关操作`,
	Action: func(context *cli.Context) {


		var id string
		if id = context.String("id");id=="" {
			fmt.Println("容器id不能为空")
		}

		pidPath := ToyContainerRunPath+"/" + id+"/pid.txt"

		//从指定目录获取pid病解析为int类型
		pidArray := util.GetPidFromFile(pidPath)
		pid := util.ParsePid(pidArray)
		// 初始化管道，用于toy-container run 进程与 toy-container init 进程之间进行通信

		parentInitPipe,childInitPipe,err := util.NewSocketPair("init")

		// xdushepherd 2019/11/18 14:11 为int程序设定管道环境变量，以便可以进行命名空间的初始化工作
		toyInit := exec.Command("/proc/self/exe", "init")

		// xdushepherd 2019/11/21 10:14	将init进程需要使用的管道复制给toyInit

		toyInit.ExtraFiles = append(toyInit.ExtraFiles,childInitPipe)

		toyInit.Env = append(toyInit.Env,fmt.Sprintf("_LIBCONTAINER_INITPIPE=%d",stdFdsCount+len(toyInit.ExtraFiles)-1))

		toyInit.Env = append(toyInit.Env,fmt.Sprintf("_TOYCONTAINER_EXEC=%d",1))

		// 读取命名空间文件
		nsPath := "/proc/" + strconv.FormatInt(int64(pid), 10) + "/ns"

		dir, err := ioutil.ReadDir(nsPath)
		if err != nil {
			fmt.Println(err)
		}

		var nsPaths []string

		for _, fi := range dir {
			if !fi.IsDir() {
				fmt.Println(fi.Name())
				nsPaths = append(nsPaths,fi.Name())
			}
		}


		nr := nl.NewNetlinkRequest(int(message.InitMsg), 0)


		//xdushepherd 2019/11/21 14:26	添加容器id，用于init进程设置容器hostname

		nr.AddData(&message.Bytemsg{
			Type:  message.NsPathsAttr,
			Value: []byte(strings.Join(nsPaths,",")),
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
