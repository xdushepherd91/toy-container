package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

func GetPidFromFile(path string) []byte {

	fp,err := os.Open(path)
	if err !=nil{
		fmt.Println(err)
		return nil
	}

	pidArray ,err := ioutil.ReadAll(fp)
	if err!=nil{
		fmt.Println(err)
		return nil
	}


	return pidArray
}

func ParsePid( pids []byte) int  {
	var pid int

	pid_str := string(pids)

	pid, err := strconv.Atoi(pid_str)

	if err != nil{
		fmt.Println(err)
	}

	return pid
}

func MountNamespace(source,dest string) error {
	if err := syscall.Mount(source,dest,"bind",uintptr(0),"");err !=nil{
		fmt.Println(err)
		return err
	}
	return nil
}