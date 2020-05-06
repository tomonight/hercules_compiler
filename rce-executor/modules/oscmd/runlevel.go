package oscmd

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"strconv"
	"strings"
)

//定义runlevel 命令
const (
	//系统自动选择runlevel命令
	CmdNameGetCurrentRunLevel = "GetCurRunLevel"
	CmdNameSetRunLevel        = "SetRunLevel"
	CmdNameIsolateRunLevel    = "IsolateRunLevel"

	//系统运行级别 使用systemctl命令
	CmdNameGetCurrentRunLevelSctrl = "GetCurRunLevelSctrl"
	CmdNameSetRunLevelSctrl        = "SetRunLevelSctrl"
	CmdNameIsolateRunLevelSctrl    = "IsolateRunLevelSctrl"

	//系统运行级别 使用init命令
	CmdNameGetCurrentRunLevelInit = "GetCurRunLevelInit"
	CmdNameSetRunLevelInit        = "SetRunLevelInit"
	CmdNameIsolateRunLevelInit    = "IsolateRunLevelInit"
)

//runlevel 参数
const (
	CmdParamRunLevel = "runLevel"
)

//runlevel 返回结果
const (
	ResultDataKeyStatus = "status"
)

var (
	gCentOs7RunLevelList = map[int]string{
		0: "poweroff.target",
		1: "rescue.target",
		2: "multi-user.target",
		3: "multi-user.target",
		4: "multi-user.target",
		5: "graphical.target",
		6: "reboot.target"}
)

func getRunLevelIndex(runLevel string) string {
	switch runLevel {
	case "poweroff.target":
		return "0"
	case "rescue.target":
		return "1"
	case "multi-user.target":
		return "3"
	case "graphical.target":
		return "5"
	case "reboot.target":
		return "6"
	}
	return ""
}

func getRunLevelIndexFromCmdInit(runLevel string) string {
	replaceFlag := "N "
	levelStr := strings.Replace(runLevel, replaceFlag, "", -1)
	if levelStr != "" {
		_, err := strconv.Atoi(levelStr)
		if err == nil {
			return levelStr
		}
	}
	return ""
}

//GetCurRunLevelSctrl centos7 获取当前runlevel 命令
func GetCurRunLevelSctrl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		//@remark 获取当前runlevel 不需要参数
		cmdStr := "systemctl get-default"
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			er := executor.SuccessulExecuteResult(es, false, "Runlevel get successful")
			var value string
			if len(es.Stdout) == 1 {
				value = es.Stdout[0]
				value = getRunLevelIndex(value)
			}

			resultData := map[string]string{ResultDataKeyStatus: value}
			er.ResultData = resultData
			return er
		}

		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//SetRunLevelSctrl centos7 设置runlevel 命令
func SetRunLevelSctrl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		level, err := executor.ExtractCmdFuncIntParam(params, CmdParamRunLevel)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		var (
			runLevelStr string
			ok          bool
		)

		runLevelStr, ok = gCentOs7RunLevelList[level]
		if !ok {
			return executor.ErrorExecuteResult(errors.New("unrecognized param"))
		}

		cmdStr := fmt.Sprintf("%s %s %s", "systemctl", "set-default", runLevelStr)
		log.Debug("cmdStr=%s", cmdStr)
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		log.Debug("SetRunLevel Result %v", es)
		log.Debug("SetRunLevel es.Stdout %v", es.Stdout)
		log.Debug("SetRunLevel es.Stderr %v", es.Stderr)
		if es.ExitCode == 0 /*&& len(es.Stderr) == 0*/ {
			er := executor.SuccessulExecuteResult(es, true, fmt.Sprintf("Runlevel %d set successful", level))
			return er
		}

		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//IsolateRunLevelSctrl centos7 切换到指定的运行级别
func IsolateRunLevelSctrl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {

		level, err := executor.ExtractCmdFuncIntParam(params, CmdParamRunLevel)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		var (
			runLevelStr string
			ok          bool
		)

		runLevelStr, ok = gCentOs7RunLevelList[level]
		if !ok {
			return executor.ErrorExecuteResult(errors.New("unrecognized param"))
		}

		cmdStr := fmt.Sprintf("%s %s %s", "systemctl", "isolate", runLevelStr)
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			er := executor.SuccessulExecuteResult(es, true, "")
			return er
		}

		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//GetCurRunLevelInit  通过init命令获取runlevel
func GetCurRunLevelInit(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		cmdStr := "runlevel"
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			er := executor.SuccessulExecuteResult(es, false, "get current run level init successful")
			var value string
			if len(es.Stdout) == 1 {
				value = es.Stdout[0]
				value = getRunLevelIndexFromCmdInit(value)
			}
			resultData := map[string]string{ResultDataKeyStatus: value}
			er.ResultData = resultData
			return er
		}

		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//SetRunLevelInit 低于centos7 设置runlevel 命令
func SetRunLevelInit(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {

		//sed s/id:5:initdefault:/id:3:initdefault:/g
		initPath := "/etc/inittab"

		exist, err := fileExistCharacter(e, initPath, "id:3:initdefault")
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if exist {
			result := new(executor.ExecuteResult)
			result.Changed = false
			result.Successful = true
			result.ExecuteHasError = false
			result.Message = "no need to set runlevel to 3"
			return *result
		}

		cmdStr := fmt.Sprintf("%s %s %s", "sed -i ", " s/id:5:initdefault:/id:3:initdefault:/g", initPath)
		fmt.Println("cmdStr=", cmdStr)
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			er := executor.SuccessulExecuteResult(es, true, "set run level init successful")
			return er
		}

		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//IsolateRunLevelInit init N   立即切换runlevel
func IsolateRunLevelInit(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		level, err := executor.ExtractCmdFuncIntParam(params, CmdParamRunLevel)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if level < 0 || level > 6 {
			return executor.ErrorExecuteResult(errors.New("unrecognized param"))
		}

		cmdStr := "init " + strconv.Itoa(level)
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			er := executor.SuccessulExecuteResult(es, true, fmt.Sprintf("isolate run level init to %d successful", level))
			return er
		}

		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//GetCurRunLevel 获取当前runlevel 命令
func GetCurRunLevel(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	return GetCurRunLevelInit(e, params)
}

//SetRunLevel 设置runlevel 命令 自动选择版本
func SetRunLevel(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
	log.Debug("Current Version %d", iVer)
	if iVer >= 7 {
		return SetRunLevelSctrl(e, params)
	}
	return SetRunLevelInit(e, params)
}

//IsolateRunLevel 切换到指定的运行级别
func IsolateRunLevel(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	return IsolateRunLevelInit(e, params)
}

//fileExistCharacter 判断某个字符串是存在某个文件中
func fileExistCharacter(e executor.Executor, fileName, text string) (bool, error) {
	containFlag := "contain"
	notContainFlag := "not contain"
	cmdStr := fmt.Sprintf("grep -wq %s %s && echo %s || echo %s ", text, fileName, containFlag, notContainFlag)
	es, err := e.ExecShell(cmdStr)
	if err == nil {
		err = executor.GetExecResult(es)
		if err == nil {
			if len(es.Stdout) > 0 {
				if containFlag == strings.TrimSpace(es.Stdout[0]) {
					return true, nil
				}
			}
		}
	}
	return false, err
}
