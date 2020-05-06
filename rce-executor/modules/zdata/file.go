package zdata

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"strings"
)

// 定义写入模式
const (
	WriteModeCover  = 0
	WrtteModeAppend = 1
)

//backFile 备份文件
//@param e 执行器
//@param sourcePath 源文件地址
//@param backPath 备份文件地址
//@return error 错误信息
func backFile(e executor.Executor, soucePath, backPath string) error {
	log.Debug("start backup file %s", soucePath)
	cmdStr := fmt.Sprintf("ls %s", soucePath)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmdStr %s executor failed:%v ", cmdStr, err)
		return err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmdStr %s executor failed:%v ", cmdStr, err)
		return err
	}

	if len(es.Stdout) > 0 {
		if soucePath != strings.TrimSpace(es.Stdout[0]) {
			return errors.New(soucePath + "not exist")
		}
	}

	backExist := false
	cmdStr = fmt.Sprintf("ls %s", backPath)
	es, err = e.ExecShell(cmdStr)
	if err == nil {
		err = executor.GetExecResult(es)
		if err == nil {
			if len(es.Stdout) > 0 {
				if backPath == strings.TrimSpace(es.Stdout[0]) {
					backExist = true
				}
			}
		}
	}

	if !backExist {
		cmdStr = fmt.Sprintf("cp %s %s", soucePath, backPath)
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			log.Warn("cmdStr %s copy cmdstr executor failed:%v ", cmdStr, err)
			return err

			err = executor.GetExecResult(es)
			if err != nil {
				log.Warn("cmdStr %s copy cmdstr executor failed:%v ", cmdStr, err)
				return err
			}
			log.Debug("backup file %s success !", soucePath)
		}
	} else {
		log.Debug("no need to backup file %s", soucePath)
	}

	return nil
}

//WriteText 为某个文件写入新的文字
//@param e       执行器
//@param newText 需要写入的文本
//@param cmpText 需要对比的文本（有些特殊对比需求，若cmpText为空 默认用newText作为对比文本）
//@param mode 写入类型 0--重新写入 1--追加
//@bool  是否变化
//@error 错误信息
func WriteText(e executor.Executor, filePath, newText, cmpText string, mode int) (bool, error) {
	log.Debug("start get  %s text", filePath)
	cmdStr := fmt.Sprintf("cat %s", filePath)
	log.Debug("get text command %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		log.Warn(cmdStr, "executor failed ", err)
		return false, err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn(cmdStr, "executor failed ", err)
		return false, err
	}

	text := strings.Join(es.Stdout, "\n")
	log.Debug("get %s success text:%s", filePath, text)

	if cmpText == "" {
		cmpText = newText
	}
	log.Debug("source text:%s", cmpText)
	if strings.Contains(text, cmpText) {
		log.Debug("%s alreay contains do text", filePath)
		return false, nil
	}

	log.Debug("%s not contains do text , start write text", filePath)
	switch mode {
	case WriteModeCover:
		cmdStr = fmt.Sprintf("echo '%s' > %s", newText, filePath)
	case WrtteModeAppend:
		cmdStr = fmt.Sprintf("echo '%s' >> %s", newText, filePath)
	default:
		return false, errors.New("unsupported write mode")
	}

	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmdStr %s executor failed %v", cmdStr, err)
		return false, err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmdStr %s executor failed %v", cmdStr, err)
		return false, err
	}
	log.Debug("write text in file %s success !", filePath)
	return true, nil
}

//SetIbZtgParameter set_ib_ztg_parameter
func SetIbZtgParameter(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e != nil {
		var (
			iVer   int
			status bool
			es     executor.ExecutedStatus
			err    error
		)

		ibZtgConf := "/etc/modprobe.d/ib_ztg.conf"
		ibZtgConfBak := "/etc/modprobe.d/ib_ztg.conf.bak"

		exist, _ := oscmd.FileExist(e, ibZtgConf)
		log.Debug("exist: %v ", exist)

		if exist {
			err = backFile(e, ibZtgConf, ibZtgConfBak)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
		} else {
			cmdstr := fmt.Sprintf("%s %s", "touch", ibZtgConf)
			log.Debug("touch cmd string %s", cmdstr)
			_, err = e.ExecShell(cmdstr)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
		}

		er := oscmd.LinuxDist(e, params)
		if !er.ExecuteHasError {
			version, ok := er.ResultData[oscmd.ResultDataKeyVersion]
			if ok {
				(*params)[oscmd.CmdParamVersion] = version
			}
		} else {
			return executor.ErrorExecuteResult(errors.New(er.Message))
		}
		iVer, err = oscmd.GetLinuxDistVer(params)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		if iVer < 7 {
			text := "options ib_ztg cmd_sg_entries=128 ch_count=8 rdma_recv_tmo=19 fast_io_fail_tmo=1 reconnect_delay=10 dev_loss_tmo=16\n"
			status, err = WriteText(e, ibZtgConf, text, "", WrtteModeAppend)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			return executor.SuccessulExecuteResult(&es, status, "ib_ztg set successful")

		}
		return executor.SuccessulExecuteResult(&es, false, "ib_ztg set successful")
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//CreateUdevRules 命令
func CreateUdevRules(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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

		if iVer < 7 {
			zdata99Rules := "/etc/udev/rules.d/99-zdata.rules"

			rules := "#set device attr\n"
			rules += "ACTION==\"add\", BUS==\"scsi\", RUN+=\"'/bin/sh -c 'echo 4 >/sys$DEVPATH/device/timeout;echo 2 >/sys$DEVPATH/device/eh_timeout''\"\n"
			rules += "ACTION==\"add|change\", SUBSYSTEM==\"block\", ATTR{device/vendor}==\"ZDATA\", ATTR{queue/scheduler}=\"noop\", ATTR{queue/rq_affinity}=\"2\", ATTR{device/queue_depth}=\"31\"\n"
			rules += "ACTION==\"add|change\", KERNEL==\"dm-*\", PROGRAM=\"'/bin/bash -c 'cat /sys/block/$name/slaves/*/device/vendor | grep ZDATA''\", ATTR{queue/scheduler}=\"noop\", ATTR{queue/rq_affinity}=\"2\", ATTR{queue/add_random}=\"0\"\n"
			rules += "#set griddisk\n"
			rules += "ENV{DM_NAME}==\"ZDATA_FDISK*\", OWNER:=\"grid\", GROUP:=\"asmadmin\", MODE=:\"660\"\n"
			rules += "ENV{DM_NAME}==\"ZDATA_SDISK*\", OWNER:=\"grid\", GROUP:=\"asmadmin\", MODE=:\"660\"\n"

			log.Debug("start check zdata99Rules:%s exist or not ", zdata99Rules)
			fileExist := false
			cmdStr := fmt.Sprintf("ls %s", zdata99Rules)
			es, err := e.ExecShell(cmdStr)
			if err == nil {
				err = executor.GetExecResult(es)
				if err == nil {
					if len(es.Stdout) > 0 {
						if zdata99Rules == strings.TrimSpace(es.Stdout[0]) {
							log.Debug("zdata99Rules:%s already exist", zdata99Rules)
							fileExist = true
						}
					}
				}
			}
			if !fileExist {
				log.Debug("zdata99Rules:%s not exist  start create create file", zdata99Rules)
				//@desc 创建文件
				cmdStr = fmt.Sprintf("touch %s", zdata99Rules)
				es, err = e.ExecShell(cmdStr)
				if err != nil {
					log.Warn("cmdStr %s executor failed %v", cmdStr, err)
					return executor.ErrorExecuteResult(err)
				}
				err = executor.GetExecResult(es)
				if err != nil {
					log.Warn("cmdStr %s executor failed %v", cmdStr, err)
					return executor.ErrorExecuteResult(err)
				}
				log.Debug("create file success, start to write text")
			}

			//@desc 写入文件
			status := false
			cmpText := strings.Replace(rules, "'", "", -1)
			status, err = WriteText(e, zdata99Rules, rules, cmpText, WrtteModeAppend)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			return executor.SuccessulExecuteResult(es, status, "udev rules created successful")
		}
		log.Debug("linux version more than 7 no need to do this step")
		es := executor.ExecutedStatus{}
		return executor.SuccessulExecuteResult(&es, false, "")
	}
	return executor.ErrorExecuteResult(errors.New("executor is nil"))
}

//ChangeLimitsConf change_limits_conf命令
func ChangeLimitsConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		status bool
		es     executor.ExecutedStatus
		err    error
	)

	limitsConf := "/etc/security/limits.conf"
	limitsConfBak := "/etc/security/limits.conf.bak"
	err = backFile(e, limitsConf, limitsConfBak)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	moreOptions := "oracle soft stack 32768\n"
	moreOptions += "oracle hard stack 32768\n"
	moreOptions += "oracle soft nofile 131072\n"
	moreOptions += "oracle hard nofile 131072\n"
	moreOptions += "oracle soft nproc 131072\n"
	moreOptions += "oracle hard nproc 131072\n"
	moreOptions += "oracle soft core unlimited\n"
	moreOptions += "oracle hard core unlimited\n"
	moreOptions += "oracle soft memlock unlimited\n"
	moreOptions += "oracle hard memlock unlimited\n"
	moreOptions += "grid soft stack 10240\n"
	moreOptions += "grid hard stack 32768\n"
	moreOptions += "grid soft nofile 131072\n"
	moreOptions += "grid hard nofile 131072\n"
	moreOptions += "grid soft nproc 131072\n"
	moreOptions += "grid hard nproc 131072\n"
	moreOptions += "grid soft core unlimited\n"
	moreOptions += "grid hard core unlimited\n"
	moreOptions += "grid soft memlock 72000000\n"
	moreOptions += "grid hard memlock 72000000\n"

	status, err = WriteText(e, limitsConf, moreOptions, "", WrtteModeAppend)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResult(&es, status, "Limits conf changed successful")

}

//AddContent2Profile do AddContent2Profile
func AddContent2Profile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		es     executor.ExecutedStatus
		err    error
		status bool
	)

	profile := "/etc/profile"
	profileBak := "/etc/profile.bak"
	err = backFile(e, profile, profileBak)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	scriptSnippet := `if [ $USER = "oracle" ] || [ $USER = "grid" ] || [ $USER = "root" ]; then`
	scriptSnippet += "\n"
	scriptSnippet += `        if  [ $SHELL = "/bin/ksh" ];  then`
	scriptSnippet += "\n"
	scriptSnippet += `              ulimit -p 16384`
	scriptSnippet += "\n"
	scriptSnippet += `              ulimit -n 65536`
	scriptSnippet += "\n"
	scriptSnippet += `        else`
	scriptSnippet += "\n"
	scriptSnippet += `              ulimit -u 16384 -n 65536`
	scriptSnippet += "\n"
	scriptSnippet += `        fi`
	scriptSnippet += "\n"
	scriptSnippet += `fi`
	scriptSnippet += "\n"

	status, err = WriteText(e, profile, scriptSnippet, "", WrtteModeAppend)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResult(&es, status, "Content added to profile file successful")
}

//AddContent2Login do AddContent2Login
func AddContent2Login(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		es     executor.ExecutedStatus
		err    error
		status bool
	)

	login := "/etc/pam.d/login"
	loginBak := "/etc/pam.d/login.bak"
	err = backFile(e, login, loginBak)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	newLine := `session    required     pam_limits.so`
	status, err = WriteText(e, login, newLine, "", WrtteModeAppend)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResult(&es, status, "Content added to login file successful")
}

//GenerateClientConf 创建etcd client 配置文件
func GenerateClientConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	mkPath := "/opt/zdata/conf/etcd_conf/"
	filePath := "/opt/zdata/conf/etcd_conf/etcd_service.conf"
	text := ""
	etcdIpList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamEtcdIpList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("ectdIpList:=%v", etcdIpList)
	etcdPortList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamEtcdPortList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("etcdPortList:=%v", etcdPortList)

	if len(etcdIpList) != len(etcdPortList) {
		return executor.ErrorExecuteResult(errors.New("etcd cluster information invalid"))
	}

	for index, value := range etcdIpList {
		tagLine := fmt.Sprintf("[host%d]\n", index)
		text += tagLine
		hostLine := fmt.Sprintf("host=%s\n", value)
		text += hostLine
		portLine := fmt.Sprintf("port=%s\n", etcdPortList[index])
		text += portLine
		text += "\n"
	}
	authLine := fmt.Sprintf("[auth]\n")
	text += authLine
	userLine := fmt.Sprintf("username=%s\n", "root")
	text += userLine
	password := fmt.Sprintf("password=%s\n", "d940f20e68c04c")
	text += password
	cmdStr := fmt.Sprintf("mkdir -p %s", mkPath)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//touch file
	cmdStr = fmt.Sprintf("touch %s", filePath)
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	changed, err := WriteText(e, filePath, text, "", WriteModeCover)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if changed {
		return executor.SuccessulExecuteResultNoData("generate client conf " + text + "successful ")
	}
	return executor.SuccessulExecuteResultNoData("")
}

//GenerateAgentConf 创建客户端配置文件
func GenerateAgentConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	mkPath := "/opt/zdata/conf/zmon_agent_conf"
	filePath := "/opt/zdata/conf/zmon_agent_conf/zmon_server_ipaddr_port.xml"

	text := ""
	monitorIp, err := executor.ExtractCmdFuncStringParam(params, CmdParamMonitorIp)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	monitorPort, err := executor.ExtractCmdFuncStringParam(params, CmdParamMonitorPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	nodeType, err := executor.ExtractCmdFuncIntParam(params, CmdParamNodeType)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	tagLine := fmt.Sprintf("<server_info>\n")
	text += tagLine
	monitorIpLine := fmt.Sprintf("<ipaddr>%s</ipaddr>\n", monitorIp)
	text += monitorIpLine
	portLine := fmt.Sprintf("<port>%s</port>\n", monitorPort)
	text += portLine
	typeLine := fmt.Sprintf("<node_type>%d</node_type>\n", nodeType)
	text += typeLine
	tagEndLine := fmt.Sprintf("</server_info>\n")
	text += tagEndLine
	text += "\n"

	cmdStr := fmt.Sprintf("mkdir -p %s", mkPath)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//touch file
	cmdStr = fmt.Sprintf("touch %s", filePath)
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	changed, err := WriteText(e, filePath, text, "", WriteModeCover)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if changed {
		return executor.SuccessulExecuteResultNoData("generate client conf " + text + "successful ")
	}
	return executor.SuccessulExecuteResultNoData("")
}

//fileExist 判断文件是否存在
//@param e 执行器
//@param filePath 源文件地址
//@return error 错误信息
func fileExist(e executor.Executor, filePath string) (bool, error) {
	cmdStr := fmt.Sprintf("ls %s", filePath)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmdStr %s executor failed:%v ", cmdStr, err)
		return false, err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmdStr %s executor failed:%v ", cmdStr, err)
		return false, err
	}

	if len(es.Stdout) > 0 {
		if filePath != strings.TrimSpace(es.Stdout[0]) {
			return false, errors.New(filePath + "not exist")
		}
	}
	return true, nil
}
