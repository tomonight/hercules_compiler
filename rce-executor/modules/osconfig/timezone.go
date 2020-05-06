package osconfig

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"strings"
)

// 参数常量定义
const (
	CmdParamTimezone = "timezone"
	CmdParamDateTime = "dateTime"
)

// 结果集键定义
const (
	ResultDataKeyDateTime = "dateTime"
)

// SetTimezone 设置时区
// Works only on RHEL
func SetTimezone(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// timezone like: Asia/Shanghai
	timezone, err := executor.ExtractCmdFuncStringParam(params, CmdParamTimezone)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf(`timedatectl set-timezone "%s"`, timezone)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "Timezone "+timezone+" set successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Index(errMsg, "timedatectl:") >= 0 && strings.Index(errMsg, "command not found") >= 0 {
			// `timedatectl` command not found, change `/etc/sysconfig/clock` file instead
			cmdstr := fmt.Sprintf(`sed -i 's/^ZONE=.*/ZONE="%s"/g' /etc/sysconfig/clock`, strings.Replace(timezone, "/", "\\/", -1))
			es, err := e.ExecShell(cmdstr)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			if es.ExitCode == 0 && len(es.Stderr) == 0 {
				// also make `/usr/share/zoneinfo/{timezone}` soft link to `/etc/localtime`
				// TODO: works only on RHEL
				cmdstr := fmt.Sprintf("ln -sf /usr/share/zoneinfo/%s /etc/localtime", timezone)
				es, err := e.ExecShell(cmdstr)
				if err != nil {
					return executor.ErrorExecuteResult(err)
				}
				if es.ExitCode == 0 && len(es.Stderr) == 0 {
					return executor.SuccessulExecuteResult(es, true, "Timezone "+timezone+" set successful")
				}
			}

			if len(es.Stderr) == 0 {
				errMsg = executor.ErrMsgUnknow
			} else {
				errMsg = es.Stderr[0]
			}
		}
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// GetDateTime 使用`date +"%F %T"`命令查询时间日期
func GetDateTime(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := `date +"%F %T"`
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "Date time get successful")
		er.ResultData = map[string]string{
			ResultDataKeyDateTime: es.Stdout[0],
		}
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

// SetDateTime 使用`date -s`设置时间日期
func SetDateTime(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// dateTime: 2018-05-18 10:38:45
	dateTime, err := executor.ExtractCmdFuncStringParam(params, CmdParamDateTime)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf(`date -s "%s"`, dateTime)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		// also set hardware date time use `hwclock -w`
		cmdstr := fmt.Sprint("hwclock -w")
		es, err := e.ExecShell(cmdstr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			return executor.SuccessulExecuteResult(es, true, "Date time "+dateTime+" set successful")
		}
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}
