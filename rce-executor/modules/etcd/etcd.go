/*
Package etcd 用于执行etcd相关的操作
*/
package etcd

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules/osservice"
	"strings"
)

// 模块名常量定义
const (
	EtcdModuleName = "etcd"
)

// 函数名常量定义
const (
	CmdNameEtcdInstallService = "EtcdInstallService"
)

// 命令参数常量定义
const (
	CmdParamEtcdExecPath         = "etcdExecPath"
	CmdParamEtcdDataDir          = "etcdDataDirectory"
	CmdParamEtcdNodeName         = "etcdNodeName"
	CmdParamEtcdNodeIP           = "etcdNodeIP"
	CmdParamEtcdListenClientPort = "etcdListenClientPort"
	CmdParamEtcdListenPeerPort   = "etcdListenPeerPort"
	CmdParamEtcdClusterName      = "etcdClusterName"
	CmdParamEtcdNodeIPList       = "etcdNodeIPList"
	CmdParamEtcdNodeNameList     = "etcdNodeNameList"
)

// 默认值常量定义
const (
	EtcdDefaultExecPath      = "/usr/bin"
	EtcdDefaultDataDirectory = "/var/lib/etcd"
	EtcdDefaultClientPort    = 2379
	EtcdDefaultPeerPort      = 2380
)

// EtcdInstallService 生成Systemd服务, 目前仅支持3节点
func EtcdInstallService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	etcdExecPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamEtcdExecPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(etcdExecPath) == 0 {
		etcdExecPath = EtcdDefaultExecPath
	}
	if etcdExecPath[len(etcdExecPath)-1] != '/' {
		etcdExecPath = etcdExecPath + "/"
	}

	etcdNodeName, err := executor.ExtractCmdFuncStringParam(params, CmdParamEtcdNodeName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	etcdClusterName, err := executor.ExtractCmdFuncStringParam(params, CmdParamEtcdClusterName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	etcdNodeIP, err := executor.ExtractCmdFuncStringParam(params, CmdParamEtcdNodeIP)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	etcdNodeNameList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamEtcdNodeNameList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(etcdNodeNameList) == 0 {
		return executor.ErrorExecuteResult(fmt.Errorf("Etcd node name list can not be empty"))
	}

	etcdNodeIPList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamEtcdNodeIPList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(etcdNodeIPList) == 0 {
		return executor.ErrorExecuteResult(fmt.Errorf("Etcd node ip list can not be empty"))
	}

	if len(etcdNodeIPList) != len(etcdNodeNameList) {
		return executor.ErrorExecuteResult(fmt.Errorf("Etcd node ip list size must be same as node name list"))
	}

	etcdListenPeerPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamEtcdListenPeerPort)
	if err != nil {
		etcdListenPeerPort = EtcdDefaultPeerPort
	}

	etcdListenClientPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamEtcdListenClientPort)
	if err != nil {
		etcdListenClientPort = EtcdDefaultClientPort
	}

	etcdDataDirectory, err := executor.ExtractCmdFuncStringParam(params, CmdParamEtcdDataDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if len(etcdDataDirectory) == 0 {
		etcdDataDirectory = EtcdDefaultDataDirectory
	}
	etcdDataDirectory = strings.TrimRight(etcdDataDirectory, "/")

	workingDirectory := etcdDataDirectory + "/"

	initialCluster := ""

	for idx := 0; idx < len(etcdNodeNameList); idx++ {
		if idx > 0 {
			initialCluster = initialCluster + fmt.Sprintf(",%s=http://%s:%d", etcdNodeNameList[idx], etcdNodeIPList[idx], etcdListenPeerPort)
		} else {
			initialCluster = initialCluster + fmt.Sprintf("%s=http://%s:%d", etcdNodeNameList[idx], etcdNodeIPList[idx], etcdListenPeerPort)
		}
	}
	serviceCmdLineTemplate := `%s \\
	--name=%s \\
	--initial-advertise-peer-urls=http://%s:%d \\
	--listen-peer-urls=http://%s:%d \\
	--listen-client-urls=http://%s:%d,http://127.0.0.1:%d \\
	--advertise-client-urls=http://%s:%d \\
	--initial-cluster-token=%s \\
	--initial-cluster=%s \\
	--initial-cluster-state=new \\
	--data-dir=%s
`
	serviceCmdLine := fmt.Sprintf(serviceCmdLineTemplate, etcdExecPath+"etcd",
		etcdNodeName, etcdNodeIP, etcdListenPeerPort, etcdNodeIP, etcdListenPeerPort, etcdNodeIP, etcdListenClientPort, etcdListenClientPort,
		etcdNodeIP, etcdListenClientPort, etcdClusterName, initialCluster, etcdDataDirectory)

	cmdParams := executor.ExecutorCmdParams{}

	cmdParams[osservice.CmdParamServiceDesc] = "ETCD Service"
	cmdParams[osservice.CmdParamServiceCmdLine] = serviceCmdLine
	cmdParams[osservice.CmdParamWorkingDir] = workingDirectory
	cmdParams[osservice.CmdParamServiceType] = "notify"
	cmdParams[osservice.CmdParamServiceName] = "etcd"

	er := osservice.GenerateSystemdService(e, &cmdParams)

	if !er.Successful {
		return er
	}

	cmdParams[oscmd.CmdParamPath] = workingDirectory
	er2 := oscmd.MakeDir(e, &cmdParams)
	er2.StartTime = er.StartTime
	er2.RemoteStartTime = er.RemoteStartTime
	er2.Changed = true
	er2.Successful = true
	er2.Message = "etcd systemd service generate successful"
	return er2
}

func init() {
	executor.RegisterCmd(EtcdModuleName, CmdNameEtcdInstallService, EtcdInstallService)
}
