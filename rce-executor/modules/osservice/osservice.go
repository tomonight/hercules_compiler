/*
Package osservice 用于执行Linux操作系统服务相关的操作
*/
package osservice

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules"
	"regexp"
	"strconv"

	//	"strconv"
	"hercules_compiler/rce-executor/log"
	"strings"
)

// 模块名常量定义
const (
	OSServiceModuleName = "osservice"
)

// 函数名常量定义
const (
	CmdNameGenerateSystemdService    = "GenerateSystemdService"
	CmdNameSystemdServiceControl     = "SystemdServiceControl"
	CmdNameSystemdServiceStatus      = "SystemdServiceStatus"
	CmdNameSysServiceStatus          = "SysServiceStatus"
	CmdNameGenerateSysvService       = "GenerateSysVService"
	CmdNameSysvServiceControl        = "SysVServiceControl"
	CmdNameChkConfigControl          = "ChkConfigControl"
	CmdNameSystemdReload             = "SystemdReload"
	CmdNameSystemdServiceListControl = "SystemdServiceListControl"
	CmdNameDeleteServiceFile         = "DeleteServiceFile"
	CmdNameGetOsVersion              = "GetOsVersion"
)

// 命令参数常量定义
const (
	CmdParamServiceName    = "serviceName"
	CmdParamServiceType    = "serviceType"
	CmdParamServiceCmdLine = "serviceCmdLine"
	CmdParamServiceDesc    = "serviceDescription"
	CmdParamWorkingDir     = "workingDirectory"
	CmdParamServiceAction  = "serviceAction"
	CmdParamLevel          = "level"
	CmdParamStatus         = "status"
	CmdParamRunOrder       = "runOrder"
	CmdParamConfDir        = "confDir"
)

//GetLinuxDistVer get linux dist version
func GetLinuxDistVer(version string) (int, error) {
	var (
		iVer int
		err  error
	)

	verArray := strings.Split(version, ".")
	if len(verArray) > 0 {
		iVer, err = strconv.Atoi(verArray[0])
		if err != nil {
			return iVer, err
		}
	}
	return iVer, nil
}

//7以上的版本使用systemd

// GenerateSystemdService 生成Systemd服务
func GenerateSystemdService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	//	version := e.GetExecutorContext(executor.ContextNameVersion)
	//	iVersion, _ := GetLinuxDistVer(version)
	//	if iVersion < 7 {
	//		return GenerateSysVService(e, params)
	//	}

	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceDescription, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceDesc)
	if err != nil {
		serviceDescription = serviceName
	}
	serviceType, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceType)
	if err != nil {
		serviceType = "simple"
	}
	workingDirectory, err := executor.ExtractCmdFuncStringParam(params, CmdParamWorkingDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceCmdLine, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceCmdLine)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceTextTemplate := `
[Unit]
Description=%s
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=%s
WorkingDirectory=%s
ExecStart=%s
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`
	serviceText := fmt.Sprintf(serviceTextTemplate, serviceDescription, serviceType, workingDirectory, serviceCmdLine)

	filename := "/tmp/" + strings.ToLower(serviceName) + ".service"

	var cmdParams executor.ExecutorCmdParams
	cmdParams = executor.ExecutorCmdParams{}

	cmdParams[oscmd.CmdParamFilename] = filename
	cmdParams[oscmd.CmdParamOutText] = serviceText

	er := oscmd.TextToFile(e, &cmdParams)

	if !er.Successful {
		return er
	}
	cmdParams = executor.ExecutorCmdParams{}
	cmdParams[oscmd.CmdParamSource] = filename
	cmdParams[oscmd.CmdParamTarget] = "/etc/systemd/system/" + strings.ToLower(serviceName) + ".service"

	er = oscmd.Copy(e, &cmdParams)
	if !er.Successful {
		return er
	}

	var cmdstr string
	cmdstr = "systemctl daemon-reload"

	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//更改权限
	changeModeParams := executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: "/etc/systemd/system/" + strings.ToLower(serviceName) + ".service",
	}
	er = oscmd.ChangeMode(e, &changeModeParams)
	if !er.Successful {
		return er
	}

	//启动服务
	startServiceStr := "systemctl start " + strings.ToLower(serviceName) + ".service"
	_, err1 := e.ExecShell(startServiceStr)
	if err1 != nil {
		return executor.ErrorExecuteResult(err1)
	}

	//enable服务
	enableServiceStr := "systemctl enable " + strings.ToLower(serviceName) + ".service"

	_, err1 = e.ExecShell(enableServiceStr)
	if err1 != nil {
		return executor.ErrorExecuteResult(err1)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "systemd service "+serviceName+" generate successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// SystemdServiceControl 对Systemd Service进行操作，包括start、stop、restart
func SystemdServiceControl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//	version := e.GetExecutorContext(executor.ContextNameVersion)
	//	iVersion, _ := GetLinuxDistVer(version)
	//	if iVersion < 7 {
	//		return SysVServiceControl(e, params)
	//	}
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceAction, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceAction)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("systemctl %s %s", serviceAction, serviceName)
	log.Debug("commond = %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "systemd service "+serviceName+" "+serviceAction+" successful")
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//SystemdServiceListControl 通过Systemd Service批量管理服务
func SystemdServiceListControl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	fmt.Println("serviceName= ", serviceName)
	if serviceName == "" {
		return executor.ErrorExecuteResult(fmt.Errorf("service name not found"))
	}
	serviceName = strings.TrimSpace(serviceName)

	serviceAction, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceAction)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceAction = strings.TrimSpace(serviceAction)
	fmt.Println("serviceAction = ", serviceAction)
	if serviceAction == "" {
		return executor.ErrorExecuteResult(fmt.Errorf("service action not found"))
	}
	//多个服务名之间以;分格
	list := strings.Split(serviceName, ";")
	serviceName = strings.Join(list, " ")

	cmdstr := fmt.Sprintf("systemctl %s %s ", serviceAction, serviceName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "systemd service "+serviceName+" "+serviceAction+" successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
		if ServiceErrorCanIgnore(serviceAction, errMsg) {
			return executor.SuccessulExecuteResult(es, true, "systemd service "+serviceName+" "+serviceAction+" successful")
		}
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// SystemdServiceStatus 获取服务状态
func SystemdServiceStatus(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	cmdstr := fmt.Sprintf("systemctl status %s", serviceName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 || es.ExitCode == 3 {
		if len(es.Stdout) == 0 {
			return executor.NotSuccessulExecuteResult(es, "Empty stdout")
		}
		er := executor.SuccessulExecuteResult(es, true, "")
		if er.ResultData == nil {
			er.ResultData = make(map[string]string)
		}
		reLoaded, err := regexp.Compile(`^\s*Loaded:\s*(\w+)[\s\w]*`)
		if err != nil {
			return executor.NotSuccessulExecuteResult(es, "regexp 'loaded' compile error: "+err.Error())
		}
		reActive, err := regexp.Compile(`^\s*Active:\s*(\w+)[\s\w]*`)
		if err != nil {
			return executor.NotSuccessulExecuteResult(es, "regexp 'active' compile error: "+err.Error())
		}
		rePID, err := regexp.Compile(`^\s*Process:\s*(\d+)[\s\w]*`)
		if err != nil {
			return executor.NotSuccessulExecuteResult(es, "regexp 'PID' compile error: "+err.Error())
		}
		foundLoaded := false
		foundActive := false
		foundPID := false
		for _, line := range es.Stdout {
			if !foundLoaded {
				ss := reLoaded.FindStringSubmatch(line)
				if len(ss) >= 2 {
					foundLoaded = true
					er.ResultData["Loaded"] = ss[1]
				}
			}
			if !foundActive {
				ss := reActive.FindStringSubmatch(line)
				if len(ss) >= 2 {
					foundLoaded = true
					er.ResultData["Active"] = ss[1]
				}
			}
			if !foundPID {
				ss := rePID.FindStringSubmatch(line)
				if len(ss) >= 2 {
					foundLoaded = true
					er.ResultData["PID"] = ss[1]
				}
			}
		}
		er.Changed = false
		er.Message = "systemd service " + serviceName + " status get successful"
		return er
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// SystemdReload systemctl daemon-reload
func SystemdReload(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := "systemctl daemon-reload"
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "systemctl daemon-reload successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// GenerateSysVService 生成SysV服务
func GenerateSysVService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	fmt.Println("GenerateSysVService start")
	//	version := e.GetExecutorContext(executor.ContextNameVersion)
	//	iVersion, _ := GetLinuxDistVer(version)
	//	if iVersion >= 7 {
	//		return GenerateSystemdService(e, params)
	//	}
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceDescription, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceDesc)
	if err != nil {
		serviceDescription = serviceName
	}
	workingDirectory, err := executor.ExtractCmdFuncStringParam(params, CmdParamWorkingDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	workingDirectory = strings.TrimSuffix(workingDirectory, "/")
	serviceCmdLine, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceCmdLine)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	runOrder, err := executor.ExtractCmdFuncIntParam(params, CmdParamRunOrder)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	confDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamConfDir)
	if err != nil {
		confDir = ""
	}
	serviceTextTemplate := `
	#!/usr/bin/env bash
	# chkconfig:345 %d %d
	# description: %s.

	# Source function library.
	. /etc/rc.d/init.d/functions

	prog=%s
	workdir=%s
	exec="%s"
	pidfile=$workdir/$prog.pid
	logfile=$workdir/$prog.log
	lockfile=/var/lock/subsys/$prog

	start() {
		pid=$(pidof $prog)
		if [ -n "$pid" ];then
		   echo "program $prog already runing pid=$pid"
		   return
		fi

		echo -n $"Starting $prog: "
		nohup $exec > $logfile 2>&1 & echo $! > $pidfile
		sleep 2
		status -p $pidfile $prog
		retval=$?
		echo
		[ $retval -eq 0 ] && touch $lockfile && success || failure
	}
	stop() {
		pid=$(pidof $prog)
		if [ -z "$pid" ];then
		   echo "program $prog already stop runing"
		   return
		fi
		
		echo -n $"Stopping $prog: "
		killproc -p $pidfile $prog
		retval=$?
		echo
		[ $retval -eq 0 ] && rm -f $lockfile $pidfile
	}
	log() {
		tail -f -n 20 $logfile
	}
	case "$1" in
		start)
			start
			;;
		stop)
			stop
			;;
		restart)
			stop
			start
			;;
		status)
			status -p $pidfile $prog
			;;
		log)
			log
			;;
		*)
			echo "Usage: service $SERVICE {start|status|stop|restart}"
			exit 1
	esac
	`
	filename := "/etc/init.d/" + strings.ToLower(serviceName)
	if confDir == "" {

		serviceText := fmt.Sprintf(serviceTextTemplate, runOrder, runOrder, serviceDescription, serviceName, workingDirectory, serviceCmdLine)

		textToFileParams := executor.ExecutorCmdParams{
			oscmd.CmdParamFilename: filename,
			oscmd.CmdParamOutText:  serviceText,
		}
		er := oscmd.TextToFile(e, &textToFileParams)
		if !er.Successful {
			return er
		}
	}
	changeModeParams := executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: filename,
	}
	er := oscmd.ChangeMode(e, &changeModeParams)
	if !er.Successful {
		return er
	}

	//enable服务
	enableServiceStr := "chkconfig --add " + filename

	_, err1 := e.ExecShell(enableServiceStr)
	if err1 != nil {
		return executor.ErrorExecuteResult(err1)
	}

	//start服务
	startServiceStr := "service " + serviceName + " start"

	_, err1 = e.ExecShell(startServiceStr)
	if err1 != nil {
		return executor.ErrorExecuteResult(err1)
	}

	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "SysV service "+serviceName+" generated successful")
}

// SysVServiceControl 对SysV服务进行操作，包括start, stop, restart等
func SysVServiceControl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	fmt.Println("do SysVServiceControl ")
	//	version := e.GetExecutorContext(executor.ContextNameVersion)
	//	iVersion, _ := GetLinuxDistVer(version)
	//	if iVersion >= 7 {
	//		return SystemdServiceControl(e, params)
	//	}
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	serviceAction, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceAction)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("service %s %s", serviceName, serviceAction)
	log.Debug("cmdStr = ", cmdstr)
	es, err := e.ExecShell(cmdstr)
	log.Debug("err ====", es.Stderr, "======out===", es.Stdout, "=========code=====", es.ExitCode)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "sysv service "+serviceName+" "+serviceAction+" successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// ChkConfigControl 使用chkconfig命令控制服务启动level
func ChkConfigControl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	level, err := executor.ExtractCmdFuncStringParam(params, CmdParamLevel)
	if err != nil {
		// default level
		level = "2345"
	}
	status, err := executor.ExtractCmdFuncStringParam(params, CmdParamStatus)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("chkconfig --level %s %s %s", level, serviceName, status)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "chkconfig "+serviceName+" "+status+" successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//SysServiceStatus  service serviceName status
func SysServiceStatus(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	cmdstr := fmt.Sprintf("service  %s status", serviceName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//ExitCode 0表示运行 3表示停止 2 表示异常
	if es.ExitCode == 0 {
		log.Debug(strings.Join(es.Stdout, ""))
		er := executor.SuccessulExecuteResult(es, true, strings.Join(es.Stdout, ""))
		if er.ResultData == nil {
			er.ResultData = make(map[string]string)
		}
		er.ResultData["Active"] = "active"
		return er
	} else if es.ExitCode == 3 {
		log.Debug(strings.Join(es.Stdout, ""))
		er := executor.SuccessulExecuteResult(es, true, strings.Join(es.Stdout, ""))
		if er.ResultData == nil {
			er.ResultData = make(map[string]string)
		}
		er.ResultData["Active"] = "inactive"
		return er

	} else {
		log.Debug(strings.Join(es.Stdout, ""))
		er := executor.SuccessulExecuteResult(es, true, strings.Join(es.Stdout, ""))
		if er.ResultData == nil {
			er.ResultData = make(map[string]string)
		}
		er.ResultData["Active"] = "abnormal"
		return er
	}

}

//DeleteServiceFile 删除文件名
func DeleteServiceFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	serviceName, err := executor.ExtractCmdFuncStringParam(params, CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//增加支持centOS6的删除
	version := e.GetExecutorContext(executor.ContextNameVersion)
	iVersion, err := GetLinuxDistVer(version)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//多个服务名之间以;分格
	serviceList := strings.Split(serviceName, ";")
	serviceFile := "/etc/systemd/system/"
	if iVersion < 7 {
		serviceFile = "/etc/init.d/"
	}
	serviceFileList := []string{}

	for _, name := range serviceList {
		serviceFileList = append(serviceFileList, serviceFile+name)
		serviceFileList = append(serviceFileList, serviceFile+name+".service")
		serviceFileList = append(serviceFileList, serviceFile+name+".service.d")
	}

	serviceFiles := strings.Join(serviceFileList, " ")
	cmdstr := fmt.Sprintf("rm -rf %s", serviceFiles)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "delete service file "+serviceFiles+" successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GetOsVersion 获取操作系统版本号
func GetOsVersion(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	version := e.GetExecutorContext(executor.ContextNameVersion)
	iVersion, err := GetLinuxDistVer(version)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	es, err := e.ExecShell("")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	er := executor.SuccessulExecuteResult(es, false, "target system version get successul")
	er.ResultData = make(map[string]string)
	er.ResultData["version"] = fmt.Sprintf("%d", iVersion)
	return er
}

func init() {
	modules.AddModule(OSServiceModuleName)
	executor.RegisterCmd(OSServiceModuleName, CmdNameGenerateSystemdService, GenerateSystemdService)
	executor.RegisterCmd(OSServiceModuleName, CmdNameSystemdServiceControl, SystemdServiceControl)
	executor.RegisterCmd(OSServiceModuleName, CmdNameSystemdServiceStatus, SystemdServiceStatus)
	executor.RegisterCmd(OSServiceModuleName, CmdNameSysServiceStatus, SysServiceStatus)
	executor.RegisterCmd(OSServiceModuleName, CmdNameGenerateSysvService, GenerateSysVService)
	executor.RegisterCmd(OSServiceModuleName, CmdNameSysvServiceControl, SysVServiceControl)
	executor.RegisterCmd(OSServiceModuleName, CmdNameChkConfigControl, ChkConfigControl)
	executor.RegisterCmd(OSServiceModuleName, CmdNameSystemdReload, SystemdReload)
	executor.RegisterCmd(OSServiceModuleName, CmdNameSystemdServiceListControl, SystemdServiceListControl)
	executor.RegisterCmd(OSServiceModuleName, CmdNameDeleteServiceFile, DeleteServiceFile)
	executor.RegisterCmd(OSServiceModuleName, CmdNameGetOsVersion, GetOsVersion)
}
