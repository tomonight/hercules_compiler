package zdata

import (
	"errors"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
)

//DisableMonitorServices 停止monitor的某些服务
func DisableMonitorServices(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		er := oscmd.LinuxDist(e, params)
		if !er.ExecuteHasError {
			version, ok := er.ResultData[oscmd.ResultDataKeyVersion]
			if ok {
				(*params)[oscmd.CmdParamVersion] = version
			}
		} else {
			return executor.ErrorExecuteResult(errors.New(er.Message))
		}
		iVer, err := oscmd.GetLinuxDistVer(params)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		log.Debug("iver %d", iVer)
		monitorFlag, err := executor.ExtractCmdFuncStringParam(params, CmdParamIsMonitor)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		isMonitor := false
		if monitorFlag == "true" {
			isMonitor = true
		}

		var cmdList []string
		if iVer >= 7 {
			cmdList = []string{
				"systemctl disable postfix",
				"systemctl stop postfix",
				"systemctl disable NetworkManager",
				"systemctl stop NetworkManager",
				"systemctl disable irqbalance",
				"systemctl stop irqbalance",
				"systemctl disable tuned",
				"systemctl stop tuned",
				"systemctl disable srpd",
				"systemctl stop srpd",
			}
			if !isMonitor {
				cmdList = append(cmdList, "systemctl disable ntpd")
				cmdList = append(cmdList, "systemctl stop ntpd")
			}
		} else {

			cmdList = []string{
				"chkconfig irqbalance off",
				"/etc/init.d/irqbalance stop",
				"chkconfig cpuspeed off",
				"/etc/init.d/cpuspeed stop",
				"chkconfig bluetooth off",
				"/etc/init.d/bluetooth stop",
				"chkconfig postfix off",
				"/etc/init.d/postfix stop",
				"chkconfig trace-cmd off",
				"/etc/init.d/trace-cmd stop",
				"chkconfig tuned off",
				"/etc/init.d/tuned stop",
				"chkconfig ktune off",
				"/etc/init.d/ktune stop",
				"chkconfig NetworkManager off",
				"/etc/init.d/NetworkManager stop",
				"chkconfig libvirt-guests off",
				"/etc/init.d/libvirt-guests stop",
				"chkconfig netfs off",
				"/etc/init.d/netfs stop",
				"chkconfig portreserve off",
				"/etc/init.d/portreserve stop",
				"chkconfig rpcbind off",
				"/etc/init.d/rpcbind stop",
				"chkconfig nfslock off",
				"/etc/init.d/nfslock stop",
				"chkconfig rpcgssd off",
				"/etc/init.d/rpcgssd stop",
				"chkconfig abrtd off",
				"/etc/init.d/abrtd stop",
				"chkconfig auditd off",
				"/etc/init.d/auditd stop"}

			if !isMonitor {
				cmdList = append(cmdList, "chkconfig ntpd off")
				cmdList = append(cmdList, "/etc/init.d/ntpd stop")
			}
		}

		var (
			es *executor.ExecutedStatus
		)
		for _, cmd := range cmdList {
			es, err = e.ExecShell(cmd)
			if err != nil {
				continue
				//return executor.ErrorExecuteResult(err)
			}
			err = executor.GetExecResult(es)
			if err != nil {
				continue
				//return executor.ErrorExecuteResult(err)
			}
		}

		return executor.SuccessulExecuteResult(es, true, "Monitor service disabled successful")

	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}
