package osconfig

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"strings"
)

// 参数常量定义
const (
	CmdParamHostname     = "hostname"
	CmdParamIPAddr       = "ipAddr"
	CmdParamHostnameList = "hostnameList"
	CmdParamIPList       = "ipList"
)

// 结果集键定义
const (
	ResultDataKeyHostname = "hostname"
)

// GetHostname 获取hostname
func GetHostname(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := "hostname"
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, false, "hostname get successful")
		er.ResultData = map[string]string{ResultDataKeyHostname: es.Stdout[0]}
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

// SetHostname 设置hostname
func SetHostname(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	hostname, err := executor.ExtractCmdFuncStringParam(params, CmdParamHostname)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("hostnamectl set-hostname %s", hostname)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "hostname "+hostname+" set successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Index(errMsg, "hostnamectl:") >= 0 && strings.Index(errMsg, "command not found") >= 0 {
			// `hostnamectl`` command not found, use `hostname` instead
			cmdstr := fmt.Sprintf("hostname %s", hostname)
			es, err := e.ExecShell(cmdstr)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			if es.ExitCode == 0 && len(es.Stderr) == 0 {
				// also write hostname to /etc/sysconfig/network
				// otherwise hostname will be restored after reboot
				// TODO: works only on RHEL
				cmdstr := fmt.Sprintf("sed -i 's/^HOSTNAME=.*/HOSTNAME=%s/g' /etc/sysconfig/network", hostname)
				es, err := e.ExecShell(cmdstr)
				if err != nil {
					return executor.ErrorExecuteResult(err)
				}
				if es.ExitCode == 0 && len(es.Stderr) == 0 {
					return executor.SuccessulExecuteResult(es, true, "hostname "+hostname+" set successful")
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

// SetHostsFile 设置`/etc/hosts`文件
// 指定IP地址和hostname，若IP已存在文件中，更新其hostname
func SetHostsFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	var (
		es           *executor.ExecutedStatus
		errMsg       string
		err          error
		ipList       []string
		hostnameList []string
	)
	// ipAddr, err := executor.ExtractCmdFuncStringParam(params, CmdParamIPAddr)
	// if err != nil {
	// 	return executor.ErrorExecuteResult(err)
	// }

	ipList, err = executor.ExtractCmdFuncStringListParam(params, CmdParamIPList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("iplist=%v", ipList)

	// hostname, err := executor.ExtractCmdFuncStringParam(params, CmdParamHostname)
	// if err != nil {
	// 	return executor.ErrorExecuteResult(err)
	// }

	hostnameList, err = executor.ExtractCmdFuncStringListParam(params, CmdParamHostnameList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("hostnameList=%v", hostnameList)

	if len(ipList) != len(hostnameList) {
		return executor.ErrorExecuteResult(errors.New("input param invalid"))
	}

	for index, ip := range ipList {
		hostname := hostnameList[index]
		//排除ip和hostname一样的时候报错
		if ip == hostname {
			es, err = e.ExecShell("hostname")
			continue
		}
		cmdstr := fmt.Sprintf(
			"grep -q %s /etc/hosts && sed -i 's/^%s.*/%s %s/g' /etc/hosts || echo '%s %s' >> /etc/hosts",
			ip, ip, ip, hostname, ip, hostname,
		)
		es, err = e.ExecShell(cmdstr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			log.Debug("ipaddress:hostname for " + ip + ":" + hostname + " set successful")
		} else {
			if len(es.Stderr) == 0 {
				errMsg = executor.ErrMsgUnknow
			} else {
				errMsg = es.Stderr[0]
			}
			if errMsg != "" {
				log.Debug("errMsg=%s", errMsg)
				return executor.NotSuccessulExecuteResult(es, errMsg)
			}
		}
	}
	ipString := strings.Join(ipList, ",")
	hostNameListtring := strings.Join(hostnameList, ",")
	return executor.SuccessulExecuteResult(es, true, "ipaddress:hostname for "+ipString+":"+hostNameListtring+" set successful")
}
