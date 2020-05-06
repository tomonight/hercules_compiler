package osconfig

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules"
)

// 模块名常量定义
const (
	OSConfigModuleName = "osconfig"
)

// 函数名常量定义
const (
	CmdNameGetHostname  = "GetHostname"
	CmdNameSetHostname  = "SetHostname"
	CmdNameSetHostsFile = "SetHostsFile"
	CmdNameSetTimezone  = "SetTimezone"
	CmdNameGetDateTime  = "GetDateTime"
	CmdNameSetDateTime  = "SetDateTime"
	CmdNameSetSysConf   = "SetSysConf"
)

func init() {
	modules.AddModule(OSConfigModuleName)
	executor.RegisterCmd(OSConfigModuleName, CmdNameGetHostname, GetHostname)
	executor.RegisterCmd(OSConfigModuleName, CmdNameSetHostname, SetHostname)
	executor.RegisterCmd(OSConfigModuleName, CmdNameSetHostsFile, SetHostsFile)
	executor.RegisterCmd(OSConfigModuleName, CmdNameSetTimezone, SetTimezone)
	executor.RegisterCmd(OSConfigModuleName, CmdNameGetDateTime, GetDateTime)
	executor.RegisterCmd(OSConfigModuleName, CmdNameSetDateTime, SetDateTime)
	executor.RegisterCmd(OSConfigModuleName, CmdNameSetSysConf, SetSysConf)
}
