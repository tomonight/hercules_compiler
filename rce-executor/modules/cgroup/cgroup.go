package cgroup

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/utils"
	"strings"
)

// 参数常量定义
const (
	CmdParamCGroupName          = "cgroupName"
	CmdParamUserName            = "userName"
	CmdParamGroupName           = "groupName"
	CmdParamControllers         = "controllers"
	CmdParamCGroupParamKey      = "cgroupParamKey"
	CmdParamCGroupParamValue    = "cgroupParamValue"
	CmdParamPID                 = "pid"
	CmdParamCmd                 = "cmd"
	CmdParamMountAt             = "mountAt"
	CmdParamFileSystemName      = "fileSystemName"
	CmdParamSubFileSystemNames  = "subFileSystemNames"
	CmdParamSubFileSystemSubDir = "subFileSystemSubDir"
)

// MountCGroupFileSystem 挂载cgroup文件系统
func MountCGroupFileSystem(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// /cgroup
	mountAt, err := executor.ExtractCmdFuncStringParam(params, CmdParamMountAt)
	if err != nil {
		mountAt = "/cgroup"
	}

	fileSystemName, err := executor.ExtractCmdFuncStringParam(params, CmdParamFileSystemName)
	if err != nil {
		fileSystemName = "cgroup_root"
	}

	// cpuset,memory
	subFileSystemNames, err := executor.ExtractCmdFuncStringParam(params, CmdParamSubFileSystemNames)
	if err != nil {
		subFileSystemNames = "cpuset,memory"
	}

	// cpuset_and_mem
	subFileSystemSubDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamSubFileSystemSubDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("mkdir -p %s; mount -t tmpfs %s %s", mountAt, fileSystemName, mountAt)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		cmdstr = fmt.Sprintf(
			"mkdir -p %s/%s; mount -t cgroup -o %s %s %s/%s",
			mountAt, subFileSystemSubDir, subFileSystemNames, subFileSystemSubDir,
			mountAt, subFileSystemSubDir,
		)
		es, err = e.ExecShell(cmdstr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			return executor.SuccessulExecuteResult(es, true, "cgroup filesystem mounted successful")
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

// CreateNewCGroup 创建新的资源控制组
func CreateNewCGroup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// cg1
	cgroupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamCGroupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	userName, err := executor.ExtractCmdFuncStringParam(params, CmdParamUserName)
	if err != nil {
		userName = "root"
	}

	groupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamGroupName)
	if err != nil {
		groupName = "root"
	}

	// cpuset,memory
	controllers, err := executor.ExtractCmdFuncStringParam(params, CmdParamControllers)
	if err != nil {
		controllers = "cpuset,memory"
	}

	// cgcreate -t root:root -a root:root -g cpuset,memory:/cg1
	cmdstr := fmt.Sprintf("cgcreate -t %s:%s -a %s:%s -g %s:/%s", userName, groupName, userName, groupName, controllers, cgroupName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "CGroup added successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Index(errMsg, "cgcreate:") >= 0 && strings.Index(errMsg, "Cgroup is not mounted") >= 0 {
			mountParams := executor.ExecutorCmdParams{
				CmdParamSubFileSystemNames:  controllers,
				CmdParamSubFileSystemSubDir: strings.Replace(controllers, ",", "_", -1),
			}
			er := MountCGroupFileSystem(e, &mountParams)
			if er.Successful {
				es, err = e.ExecShell(cmdstr)
				if err != nil {
					return executor.ErrorExecuteResult(err)
				}
				if es.ExitCode == 0 && len(es.Stderr) == 0 {
					return executor.SuccessulExecuteResult(es, true, "CGroup added successful")
				}

				if len(es.Stderr) == 0 {
					errMsg = executor.ErrMsgUnknow
				} else {
					errMsg = es.Stderr[0]
				}
			} else {
				errMsg = er.Message
			}
		}
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// SetCGroupParams 设置控制组参数
func SetCGroupParams(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cgroupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamCGroupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// cpuset.cpus
	paramKey, err := executor.ExtractCmdFuncStringParam(params, CmdParamCGroupParamKey)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// 0-1
	paramValue, err := executor.ExtractCmdFuncStringParam(params, CmdParamCGroupParamValue)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// cgset -r cpuset.cpus=0-1 cg1
	cmdstr := fmt.Sprintf("cgset -r %s=%s %s", paramKey, paramValue, cgroupName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "CGroup parameter set successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// AddPIDIntoCGroup 将进程添加进控制组
func AddPIDIntoCGroup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// cg1
	cgroupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamCGroupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// cpuset,memory
	controllers, err := executor.ExtractCmdFuncStringParam(params, CmdParamControllers)
	if err != nil {
		controllers = "cpuset,memory"
	}

	// 支持多个pid：1121 1122 1123
	pid, err := executor.ExtractCmdFuncStringParam(params, CmdParamPID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("cgset -r cpuset.mems=0 %s && cgclassify -g %s:/%s %s", cgroupName, controllers, cgroupName, pid)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "Add process into cgroup successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// ExecCmdWithCGroup 使用控制组执行命令
func ExecCmdWithCGroup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// cg1
	cgroupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamCGroupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// cpuset,memory
	controllers, err := executor.ExtractCmdFuncStringParam(params, CmdParamControllers)
	if err != nil {
		controllers = "cpuset,memory"
	}

	cmd, err := executor.ExtractCmdFuncStringParam(params, CmdParamCmd)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	cmd = utils.EscapeShellCmd(cmd)

	cmdstr := fmt.Sprintf("cgexec -g %s:/%s %s", controllers, cgroupName, cmd)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "Execute cmd with cgroup successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}
