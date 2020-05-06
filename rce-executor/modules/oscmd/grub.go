package oscmd

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"strings"
)

//grub相关参数
const (
	ResultDataGrubFileKey = "grubFile"
	CmdNameGetGrubFile    = "GetGrubFile"
)

//GetGrubFile 获取GrubFile 文件路径
func GetGrubFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	er := LinuxDist(e, params)
	if !er.ExecuteHasError {
		version, ok := er.ResultData[ResultDataKeyVersion]
		if ok {
			(*params)[CmdParamVersion] = version
		}
	} else {
		return executor.ErrorExecuteResult(errors.New(er.Message))
	}

	iVer, err := GetLinuxDistVer(params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	knownGrubList := []string{}
	if iVer >= 7 {
		knownGrubList = []string{"/etc/grub2.cfg", "/etc/grub2-efi.cfg"}
	} else {
		knownGrubList = []string{"/etc/grub.conf"}
	}

	//判断文件是否存在
	existGrubList := []string{}
	for _, filePath := range knownGrubList {
		exist, err := FileExist(e, filePath)
		if err == nil {
			if exist {
				existGrubList = append(existGrubList, filePath)
			}
		} else {
			continue
		}
	}
	lastGrubList := []string{}
	//获取存在文件全路径
	for _, filePath := range existGrubList {
		fullPath, err := getFileFullPath(e, filePath)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		lastGrubList = append(lastGrubList, fullPath)
	}

	log.Debug("grub files list :%v", lastGrubList)
	result := new(executor.ExecuteResult)

	result.Changed = false
	result.Successful = true
	result.ExecuteHasError = false
	resultData := make(map[string]string)
	resultData[ResultDataGrubFileKey] = strings.Join(lastGrubList, ",")
	result.Message = "get grub file successful"
	result.ResultData = resultData
	log.Debug("result= %v", result)
	return *result

}

//FileExist 判断文件是否存在
func FileExist(e executor.Executor, filePath string) (bool, error) {
	cmdStr := fmt.Sprintf("ls %s", filePath)
	es, err := e.ExecShell(cmdStr)
	if err == nil {
		err = executor.GetExecResult(es)
		if err == nil {
			if len(es.Stdout) > 0 {
				if filePath == strings.TrimSpace(es.Stdout[0]) {
					return true, nil
				}
			}
		}
	}
	return false, err
}

//PathExist 判断目录是否存在
func PathExist(e executor.Executor, filePath string) (bool, error) {
	cmdStr := fmt.Sprintf("ls %s", filePath)
	es, err := e.ExecShell(cmdStr)
	if err == nil {
		if len(es.Stderr) == 0 {
			return true, nil
		}
	}
	return false, err
}

//getFileFullPath 获取文件全路径
func getFileFullPath(e executor.Executor, filePath string) (string, error) {
	cmdStr := fmt.Sprintf("readlink -f %s", filePath)
	es, err := e.ExecShell(cmdStr)
	if err == nil {
		err = executor.GetExecResult(es)
		if err == nil {
			if len(es.Stdout) > 0 {
				filePath = strings.TrimSpace(es.Stdout[0])
				return filePath, nil
			}
		}
	}
	return filePath, err
}
