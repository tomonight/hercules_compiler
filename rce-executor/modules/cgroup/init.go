package cgroup

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules"
)

// 模块名常量定义
const (
	CGroupModuleName = "cgroup"
)

// 函数名常量定义
const (
	CmdNameCreateNewCGroup       = "CreateNewCGroup"
	CmdNameSetCGroupParams       = "SetCGroupParams"
	CmdNameAddPIDIntoCGroup      = "AddPIDIntoCGroup"
	CmdNameExecCmdWithCGroup     = "ExecCmdWithCGroup"
	CmdNameMountCGroupFileSystem = "MountCGroupFileSystem"
)

func init() {
	modules.AddModule(CGroupModuleName)
	executor.RegisterCmd(CGroupModuleName, CmdNameCreateNewCGroup, CreateNewCGroup)
	executor.RegisterCmd(CGroupModuleName, CmdNameSetCGroupParams, SetCGroupParams)
	executor.RegisterCmd(CGroupModuleName, CmdNameAddPIDIntoCGroup, AddPIDIntoCGroup)
	executor.RegisterCmd(CGroupModuleName, CmdNameExecCmdWithCGroup, ExecCmdWithCGroup)
	executor.RegisterCmd(CGroupModuleName, CmdNameMountCGroupFileSystem, MountCGroupFileSystem)
}
