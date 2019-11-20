package main

import (
	"fmt"
	"github.com/urfave/cli"
//	"os"
	"io/ioutil"
	"strconv"
	"toy-container/util"
)

var execCommand = cli.Command{
	Name:  "exec",
	Usage: `重新进入一个容器命名空间，进行相关操作`,
	Action: func(context *cli.Context) {

		//从指定目录获取pid病解析为int类型
		pidArray := util.GetPidFromFile("/opt/toy-container/child_pid.txt")
		pid := util.ParsePid(pidArray)

		// 读取命名空间文件
		nsPath := "/proc/"+strconv.FormatInt(int64(pid),10)+"/ns"
/*		err := filepath.Walk(nsPath, func(path string, info os.FileInfo, err error) error {
			println(path)
			return nil
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

*/

    dir, err := ioutil.ReadDir(nsPath)
    if err != nil {
	fmt.Println(err)
    }

        for _, fi := range dir {
        if !fi.IsDir() {
		fmt.Println(fi.Name())
        }
    }

	},
}
