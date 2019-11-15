package overlay2

import (
	"toy-container/config"
	"fmt"
	"golang.org/x/sys/unix"
)

/*
	xdushepherd 2019/11/15 14:54
	setUpFS方法的主要目标是根据配置信息，构建一个overlay2文件系统，用于容器的运行时环境
 */
func SetUpFS(config config.Config) error{
	source := "overlay"
	mountType := "overlay"
	target := "/tmp/toy-container/target"
	flag := uintptr(0)
	data := "lowerdir=/tmp/toy-container/rootfs,upperdir=/tmp/toy-container/application,workdir=/tmp/toy-container/work"

	if err := unix.Mount(source,target,mountType,flag,data); err !=nil {
		fmt.Println(err)
	}

	return nil

}
