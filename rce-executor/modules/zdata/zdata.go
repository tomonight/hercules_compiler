package zdata

//定义 zdata 模块名
const (
	ZDataModuleName = "zdata"
)

//定义 zdata 命令
const (
	CmdNameSetYumSource              = "SetYumSource"
	CmdNameSetIBZtgParameter         = "SetIbZtgParameter"
	CmdNameCreateUdevRules           = "CreateUdevRules"
	CmdNameChangeLimitsConf          = "ChangeLimitsConf"
	CmdNameAddContent2Profile        = "AddContent2Profile"
	CmdNameAddContent2Login          = "AddContent2Login"
	CmdNameInstallIBDriver           = "InstallIBDriver"
	CmdNameStattIBDriver             = "StartIBService"
	CmdNameCmdDisableMonitotServices = "DisableMonitorServices"
	CmdNameGenerateNodeRas           = "GenerateNodeRsa"
	CmdNameInstallService            = "InstallService"
	CmdNameInstallZcli               = "InstallZcli"
	CmdNameInstallFireman            = "InstallFireman"
	CmdNameInstallZMonAgent          = "InstallZMonAgent"
	CmdNameInstallCellManager        = "InstallCellManager"
	CmdNameInitStorageNode           = "InitStorageNode"
	CmdNameInitComputeNode           = "InitComputeNode"
	CmdNameInstallCxOracle           = "InstallCxOracle"
	CmdNameInitDirs                  = "InitDirs"
	CmdNameGenerateClientConf        = "GenerateClientConf"
	CmdNameGenerateAgentConf         = "GenerateAgentConf"
)

//定义 zdata命令参数
const (
	CmdParamaSSHRsaPubFile = "SSHRsaPubFile"
	CmdParamaNodeList      = "NodeList"
	CmdParamIsMonitor      = "isMonitor"
	CmdParamHomeDir        = "homeDir"
	CmdParamInstallDir     = "installDir"
	CmdParamEtcdIpList     = "etcdIpList"
	CmdParamEtcdPortList   = "etcdPortList"
	CmdParamMonitorIp      = "monitorIp"
	CmdParamMonitorPort    = "monitorPort"
	CmdParamNodeType       = "nodeType"
)

const (
	defaultHomeDir    = "/opt/zdata"
	defaultInstallDir = "/tmp/zmanager"
)
