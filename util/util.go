package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
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

func FileExists(path string) bool {

	_,err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		if os.IsExist(err){
			return true
		}
		return false
	}

	return true
}

func IsDir(path string) bool  {

	s,err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return s.IsDir()
}


func ParsePid( pids []byte) int  {
	var pid int

	pid_str := strings.Replace(string(pids), "\n", "", -1)

	pid, err := strconv.Atoi(pid_str)

	if err != nil{
		fmt.Println(err)
	}

	return pid
}

//func MountNamespace(source,dest string) error {
//	if err := syscall.Mount(source,dest,"bind",uintptr(0),"");err !=nil{
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}
