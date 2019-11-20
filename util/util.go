package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"toy-container/config"
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

func saveConfig(config config.Config,path string) error {

	fp , err := os.OpenFile(path,os.O_CREATE,0)
	if err != nil {
		return handleError(err)
	}

	data ,err := json.Marshal(config)
	if err != nil {
		return handleError(err)
	}

	_,err  = fp.Write(data)

	return nil
}

func LoadConfig(path string)   (config.Config,error)  {

	var result config.Config
	/*
		xdushepherd 2019/11/20 17:51
		打开文件
	 */
	fp,err := os.Open(path)
	if err != nil{
		return result, handleError(err)
	}
	defer fp.Close()

	bytes,err := ioutil.ReadAll(fp)

	if err != nil {
		return result,handleError(err)
	}



	err = json.Unmarshal(bytes,&result)

	if err != nil {
		return result, handleError(err)
	}

	return result,nil
}

func IsDirOrFileExist(path string) bool {

	_,err := os.Stat(path)
	if err == nil {
		return true;
	}


	return false
}


func handleError(err error) error {
	fmt.Println(err)
	return err
}
