package overlay2

import (
	"fmt"
	"os"
	"toy-container/config"

	"golang.org/x/sys/unix"
)

/*
	xdushepherd 2019/11/15 14:54
	setUpFS方法的主要目标是根据配置信息，构建一个overlay2文件系统，用于容器的运行时环境
 */
func SetUpFS(config config.Config) error{
	source := "overlay"
	mountType := "overlay"
	target := "/tmp/toy-container/"+config.Id+"/target"
	if err := os.MkdirAll(target, 777); err !=nil {
		fmt.Println(err)
		return err
	}
	flag := uintptr(0)
	data := "lowerdir=/tmp/toy-container/rootfs,upperdir=/tmp/toy-container/application,workdir=/tmp/toy-container/work"

	return unix.Mount(source,target,mountType,flag,data)
}
