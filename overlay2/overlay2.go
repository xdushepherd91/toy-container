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

	xdushepherd 2019/11/15 17:17
	目前已经可以使用该函数搭建简易的联合文件系统
    后续需要做的事情：
    1. 向配置文件中增加lowerdir的配置数组
    2. 向配置文件中增加upperdir的配置项目
    3. 手动创建workdir目录
    4. 其他
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

	if err := unix.Mount(source,target,mountType,flag,data); err !=nil {
		fmt.Println(err)
	}

	return nil
}
