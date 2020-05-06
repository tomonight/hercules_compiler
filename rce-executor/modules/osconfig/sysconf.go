package osconfig

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"strconv"
	"strings"
)

// 参数常量定义
const (
	CmdParamConfFile    = "confFile"
	CmdParamForceUpdate = "forceUpdate"
	CmdParamConfItem    = "confItem"
)

// SetSysConf 设置系统配置文件配置项，如`/etc/sysctl.conf`、`/etc/lvm/lvm.conf`等
func SetSysConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// confFile like: `/etc/sysctl.conf`
	confFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamConfFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// force update if item already configured or raise error
	// default to true
	forceUpdate, err := executor.ExtractCmdFuncBoolParam(params, CmdParamForceUpdate)
	if err != nil {
		forceUpdate = true
	}

	// confItem like: `vm.max_map_count = 262144`
	confItem, err := executor.ExtractCmdFuncStringParam(params, CmdParamConfItem)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	confItemKV := strings.Split(confItem, "=")
	confItemKey := strings.TrimSpace(confItemKV[0])

	cmdstr := fmt.Sprintf(`grep -c "%s" %s`, confItemKey, confFile)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if (es.ExitCode == 1 || es.ExitCode == 0) && len(es.Stderr) == 0 {
		keyCount, err := strconv.Atoi(es.Stdout[0])
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if keyCount == 0 {
			// not configured
			cmdstr = fmt.Sprintf("echo '%s' >> %s", confItem, confFile)
		} else {
			if forceUpdate {
				// allow force update
				escapedConfItem := strings.Replace(confItem, "/", "\\/", -1)
				cmdstr = fmt.Sprintf("sed -i 's/%s\\s*=.*/%s/g' %s", confItemKey, escapedConfItem, confFile)
				fmt.Print(cmdstr)
			} else {
				return executor.NotSuccessulExecuteResult(es, fmt.Sprintf(`"%s" already configured`, confItemKey))
			}
		}
		es, err := e.ExecShell(cmdstr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			return executor.SuccessulExecuteResult(es, true, confItem+" in "+confFile+" configured successful")
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
