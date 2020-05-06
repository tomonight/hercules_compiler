package executor

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/log"
	"strconv"
	"strings"
)

//ExecuteResult 执行结果
type ExecuteResult struct {
	Changed         bool //表示是否执行引起改变
	Successful      bool //实际动作是否成功
	Message         string
	ExecuteHasError bool
	StartTime       int64
	StopTime        int64
	ExitCode        int64
	RemoteStartTime int64
	RemoteStopTime  int64
	ResultData      map[string]string //获取信息类返回结果
}

//ExecutedStatus 执行状态
type ExecutedStatus struct {
	StartTime       int64
	StopTime        int64
	RemoteStartTime int64
	RemoteStopTime  int64
	ExitCode        int64
	Stdout          []string
	Stderr          []string
	ErrorMessage    string
}

//定义错误信息
const (
	ErrMsgUnknow       = "Unknow Error"
	ErrCmdNotFound     = "Command Not Found"
	ErrProcessNotFound = "Process Not Found"
)

//定义操作系统变量
const (
	Linux   = "linux"
	Windows = "windows"
	MacOS   = "macos"
)

// 定义context
const (
	ContextNameOSType               = "OSType"
	ContextNameTempDir              = "OSTempDir"
	ContextNamePathSeparator        = "OSPathSeparator"
	ContextNameVersion              = "version"
	ContextNameKernel               = "kernel"
	ContextNameArchitecture         = "architecture"
	ContextNameDist                 = "dist"
	ContextNamePsuedoname           = "psuedoname"
	ContextNameHostname             = "hostname"
	ContextNameMemoryInfo           = "memoryInfo"
	ContextNameMemorySize           = "memorySize"
	ContextNameCPUProcessorCount    = "cpuProcessorCount"
	ContextNameCPUPhysicalCoreCount = "cpuPhysicalCoreCount"
	ContextNameCPUModelName         = "cpuModelName"
)

//定义执行协议
const (
	OTRCE = "rce"
	OTSSH = "ssh"
)

//timeout
const (
	TimeOutStr = "timeout"
)

//Executor 执行器接口
type Executor interface {
	ExecShell(shellcmd string) (*ExecutedStatus, error)
	GetExecutorContext(contextName string) string
	SetWorkingPath(workingPath string)
	SetEnv(envName, envValue string)
	SetExecuteUser(username string)
	SetSudoEnabled(sudoEnabled bool)
	Close()
	SetTimeOut(timeout int)
	GetTimeOut() int
}

//ExecutorCmdParams 执行参数定义
type ExecutorCmdParams map[string]interface{}

//ExecutorCmdFunc 执行函数定义
type ExecutorCmdFunc func(Executor, *ExecutorCmdParams) ExecuteResult

var cmdRegistry map[string]ExecutorCmdFunc

//ExecutedStatus2ExecuteResult 转换执行状态到执行结果
func ExecutedStatus2ExecuteResult(er *ExecuteResult, es *ExecutedStatus) {
	er.StartTime = es.StartTime
	er.StopTime = es.StopTime
	er.RemoteStartTime = es.RemoteStartTime
	er.RemoteStopTime = es.RemoteStopTime
	er.ExitCode = es.ExitCode
}

//ErrorExecuteResult 返回错误的结果
func ErrorExecuteResult(err error) ExecuteResult {
	return ExecuteResult{Successful: false, Message: err.Error(), Changed: false, ExecuteHasError: true}
}

//SuccessulExecuteResult 返回成功的结果
func SuccessulExecuteResult(es *ExecutedStatus, changed bool, msg string) ExecuteResult {
	var er ExecuteResult
	ExecutedStatus2ExecuteResult(&er, es)
	er.Successful = true
	er.Changed = changed
	er.Message = msg
	return er
}

//SuccessulExecuteResultNoData 返回成功不带其他数据
func SuccessulExecuteResultNoData(msg string) ExecuteResult {
	return ExecuteResult{Successful: true, Changed: false, Message: msg}
}

//NotSuccessulExecuteResult 返回未成功的结果
func NotSuccessulExecuteResult(es *ExecutedStatus, errMsg string) ExecuteResult {
	var er ExecuteResult
	ExecutedStatus2ExecuteResult(&er, es)
	er.Successful = false
	er.Changed = false
	er.Message = errMsg
	return er
}

//ExtractCmdFuncStringParam 获取string类型的参数
func ExtractCmdFuncStringParam(params *ExecutorCmdParams, paramName string) (string, error) {
	var v interface{}
	var ok bool
	v, ok = (*params)[paramName]
	if !ok {
		return "", fmt.Errorf("parameter '%s' does not exists", paramName)
	}
	value, found := v.(string)
	if found {
		return value, nil
	}
	return "", fmt.Errorf("type of parameter '%s' is not string", paramName)
}

//ExtractCmdFuncIntParam 获取int类型的参数
func ExtractCmdFuncIntParam(params *ExecutorCmdParams, paramName string) (int, error) {
	var v interface{}
	var ok bool
	v, ok = (*params)[paramName]
	if !ok {
		return 0, fmt.Errorf("parameter '%s' does not exists", paramName)
	}
	value, found := v.(int)
	if found {
		return value, nil
	}
	str, found := v.(string)
	if found {
		if value, err := strconv.Atoi(str); err == nil {
			return value, nil
		}
	}
	return 0, fmt.Errorf("type of parameter '%s' is not int", paramName)
}

//ExtractCmdFuncBoolParam 获取bool类型的参数
func ExtractCmdFuncBoolParam(params *ExecutorCmdParams, paramName string) (bool, error) {
	var v interface{}
	var ok bool
	v, ok = (*params)[paramName]
	if !ok {
		return false, fmt.Errorf("parameter '%s' does not exists", paramName)
	}
	value, found := v.(bool)
	if found {
		return value, nil
	}
	str, found := v.(string)
	if found {
		str = strings.ToLower(str)
		if str == "yes" || str == "true" {
			return true, nil
		} else if str == "no" || str == "false" {
			return false, nil
		}
	}

	return false, fmt.Errorf("type of parameter '%s' is not bool", paramName)
}

//ExtractCmdFuncStringListParam 获取[]string 类型的参数
func ExtractCmdFuncStringListParam(params *ExecutorCmdParams, paramName string, seperator string) ([]string, error) {
	var v interface{}
	var ok bool
	v, ok = (*params)[paramName]
	if !ok {
		return nil, fmt.Errorf("parameter '%s' does not exists", paramName)
	}
	value, found := v.(string)
	if found {
		return strings.Split(value, seperator), nil
	}
	return nil, fmt.Errorf("type of parameter '%s' is not string", paramName)
}

func init() {
	cmdRegistry = make(map[string]ExecutorCmdFunc)
}

//GetCmdByModuleAndName 获取命令执行函数
//warning: not concurrent safe
func GetCmdByModuleAndName(moduleName string, cmdName string) (ExecutorCmdFunc, bool) {
	name := fmt.Sprintf("%s.%s", moduleName, cmdName)
	v, ok := cmdRegistry[name]
	return v, ok
}

//RegisterCmd 注册命令
//warning: not concurrent safe
func RegisterCmd(moduleName string, cmdName string, cmdFunc ExecutorCmdFunc) {
	name := fmt.Sprintf("%s.%s", moduleName, cmdName)
	cmdRegistry[name] = cmdFunc
}

//ExecCmd 执行命令
func ExecCmd(executor Executor, cmdType, cmd string, params ExecutorCmdParams) (map[string]string, error) {
	res := map[string]string{}
	if executor == nil {
		return res, errors.New("executor is nil")
	}
	cmdFunc, ok := GetCmdByModuleAndName(cmdType, cmd)
	if ok {
		er := cmdFunc(executor, &params)
		if er.Successful {
			res = er.ResultData
			return res, nil
		}
		return res, errors.New(er.Message)

	}
	return res, errors.New(ErrCmdNotFound)
}

//GetExecResult 通过执行结果判断执行结果
//es 执行结果
func GetExecResult(es *ExecutedStatus) error {
	if es.ErrorMessage != "" {
		return errors.New(es.ErrorMessage)
	}
	if es.ExitCode != 0 {
		log.Debug("es.ExitCode=%d", es.ExitCode)
		if len(es.Stderr) == 0 {
			return errors.New(ErrMsgUnknow)
		}
		errMsg := strings.Join(es.Stderr, ",")
		return errors.New(errMsg)
	}
	return nil
}
