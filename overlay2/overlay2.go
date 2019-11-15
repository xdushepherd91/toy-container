package overlay2

import (
	"toy-container/config"

	"golang.org/x/sys/unix"
)

/*
	xdushepherd 2019/11/15 14:54
	setUpFS方法的主要目标是根据配置信息，构建一个overlay2文件系统，用于容器的运行时环境
 */
func setUpFS(config config.Config) {
	source := "overlay"
	mountType := "overlay"
	target := "/tmp/toy-contariner/target"
	flag := uintptr(0)
	data := "-o lowerdir=/tmp/toy-contariner/rootfs,upperdir=/tmp/toy-contariner/application,workdir=/tmp/toy-contariner/work"

	unix.Mount(source,target,mountType,flag,data)
}
