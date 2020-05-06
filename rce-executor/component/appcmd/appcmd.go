package appcmd

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"strings"
)

//定义业务命令模块名
const (
	AppModuleName = "appcmd"
)

//定义业务命令名
const (
	CmdNameServiceStart = "ServiceStart"
	CmdNameServiceStop  = "ServiceStop"
)

// ServiceStart 启动服务
func ServiceStart(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	executable, err := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamExecutable)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	processStatusParams := executor.ExecutorCmdParams{
		oscmd.CmdParamProcessName: executable,
	}
	processStatusResult := oscmd.ProcessStatus(e, &processStatusParams)
	if processStatusResult.Successful {
		// process exist
		if stat := processStatusResult.ResultData[oscmd.ResultDataKeyStat]; strings.HasPrefix(stat, "Z") {
			// process exist but is zombie stat
			// kill it first
			pid := processStatusResult.ResultData[oscmd.ResultDataKeyPID]
			killProcessByPIDParams := executor.ExecutorCmdParams{
				oscmd.CmdParamPID:       pid,
				oscmd.CmdParamForceKill: true,
			}
			killProcessByPIDResult := oscmd.KillProcessByPID(e, &killProcessByPIDParams)
			if !killProcessByPIDResult.Successful {
				return executor.ErrorExecuteResult(errors.New(killProcessByPIDResult.Message))
			}
		} else {
			// process exist and is running
			return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, false, "Process exist and running")
		}
	} else if !processStatusResult.Successful && processStatusResult.Message != executor.ErrMsgUnknow {
		// `ps` command failed
		return executor.ErrorExecuteResult(errors.New(processStatusResult.Message))
	}
	// process not exist (or killed because of zombie stat)
	logFile, err := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamLogFile)
	if err != nil {
		logFile = fmt.Sprintf("%s.log", executable)
	}

	pidFile, err := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamPIDFile)
	if err != nil {
		pidFile = fmt.Sprintf("%s.pid", executable)
	}

	nohupParams := executor.ExecutorCmdParams{
		oscmd.CmdParamExecutable: executable,
		oscmd.CmdParamLogFile:    logFile,
		oscmd.CmdParamPIDFile:    pidFile,
	}
	nohupResult := oscmd.Nohup(e, &nohupParams)
	if nohupResult.Successful {
		return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "Process started successful")
	}
	return executor.ErrorExecuteResult(errors.New(nohupResult.Message))
}

// ServiceStop 停止服务
func ServiceStop(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	executable, err := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamExecutable)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	processStatusParams := executor.ExecutorCmdParams{
		oscmd.CmdParamProcessName: executable,
	}
	processStatusResult := oscmd.ProcessStatus(e, &processStatusParams)
	if processStatusResult.Successful {
		pid := processStatusResult.ResultData[oscmd.ResultDataKeyPID]
		killProcessByPIDParams := executor.ExecutorCmdParams{
			oscmd.CmdParamPID:       pid,
			oscmd.CmdParamForceKill: true,
		}
		killProcessByPIDResult := oscmd.KillProcessByPID(e, &killProcessByPIDParams)
		if killProcessByPIDResult.Successful {
			return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "Process stopped successful")
		}
		return executor.ErrorExecuteResult(errors.New(killProcessByPIDResult.Message))
	} else if !processStatusResult.Successful && processStatusResult.Message == executor.ErrMsgUnknow {
		return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, false, "Process not running")
	} else {
		return executor.ErrorExecuteResult(errors.New(processStatusResult.Message))
	}
}

func init() {
	// executor.RegisterCmd(MODULE_NAME_APP, CMD_NAME_UPLOADFILE_DOMD5, UploadFileDoMd5)
	executor.RegisterCmd(AppModuleName, CmdNameServiceStart, ServiceStart)
	executor.RegisterCmd(AppModuleName, CmdNameServiceStop, ServiceStop)
}
