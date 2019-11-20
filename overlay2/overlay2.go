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
	/*
		xdushepherd 2019/11/20 19:48
		生成容器根目录，用于容纳overlay2相关目录
	 */
	containerRoot := config.ContainerPath
	if err := os.Mkdir(containerRoot, 777); err !=nil {
		fmt.Println(err)
		return err
	}

	flag := uintptr(0)

	bases := []string{
		"rootfs",
		"testfs",
	}

	lowerDir := "lowerdir="
	for _,dir := range bases {
		lowerDir =  lowerDir + containerRoot+"/" + dir
	}

	applicationDir := containerRoot + "/app/"
	upperDir := " upperdir="+"/"+applicationDir
	workdir := containerRoot + "/work"
	merged := containerRoot + "/merged"

	data := lowerDir+ "," +upperDir + "," + workdir

	if err := unix.Mount(config.MountSource,merged,config.MountType,flag,data); err !=nil {
		fmt.Println(err)
	}

	println("target文件系统mount成功")
	return nil
}
