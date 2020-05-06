package zdata

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules"
)

func init() {
	modules.AddModule(ZDataModuleName)
	executor.RegisterCmd(ZDataModuleName, CmdNameSetYumSource, SetYumSource)
	executor.RegisterCmd(ZDataModuleName, CmdNameCmdDisableMonitotServices, DisableMonitorServices)
	executor.RegisterCmd(ZDataModuleName, CmdNameSetIBZtgParameter, SetIbZtgParameter)
	executor.RegisterCmd(ZDataModuleName, CmdNameCreateUdevRules, CreateUdevRules)
	executor.RegisterCmd(ZDataModuleName, CmdNameChangeLimitsConf, ChangeLimitsConf)
	executor.RegisterCmd(ZDataModuleName, CmdNameAddContent2Profile, AddContent2Profile)
	executor.RegisterCmd(ZDataModuleName, CmdNameAddContent2Login, AddContent2Login)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallIBDriver, InstallIBDriver)
	executor.RegisterCmd(ZDataModuleName, CmdNameStattIBDriver, StartIBService)
	executor.RegisterCmd(ZDataModuleName, CmdNameGenerateNodeRas, GenerateNodeRsa)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallService, InstallService)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallZcli, InstallZcli)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallFireman, InstallFireman)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallZMonAgent, InstallZMonAgent)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallCellManager, InstallCellManager)
	executor.RegisterCmd(ZDataModuleName, CmdNameInitStorageNode, InitStorageNode)
	executor.RegisterCmd(ZDataModuleName, CmdNameInitComputeNode, InitComputeNode)
	executor.RegisterCmd(ZDataModuleName, CmdNameInstallCxOracle, InstallCxOracle)
	executor.RegisterCmd(ZDataModuleName, CmdNameInitDirs, InitDirs)
	executor.RegisterCmd(ZDataModuleName, CmdNameRemoveModule, RmMpd)
	executor.RegisterCmd(ZDataModuleName, CmdNameGenerateClientConf, GenerateClientConf)
	executor.RegisterCmd(ZDataModuleName, CmdNameGenerateAgentConf, GenerateAgentConf)
}
