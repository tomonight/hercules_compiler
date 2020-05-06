package zdata

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"strings"
)

// InstallFireman 安装fireman服务
func InstallFireman(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	homeDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamHomeDir)
	if err != nil {
		homeDir = defaultHomeDir
	}
	homeDir = strings.TrimRight(homeDir, "/")
	installDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamInstallDir)
	if err != nil {
		installDir = defaultInstallDir
	}
	installDir = strings.TrimRight(installDir, "/")

	var (
		zMonDir        = homeDir + "/zmon"
		toolsDir       = homeDir + "/tools"
		toolsDirPath   = installDir + "/tools"
		binDir         = homeDir + "/bin"
		confDir        = homeDir + "/conf"
		firemanDirPath = installDir + "/src/fireman"

		deployDirPath    = installDir + "/deploy"
		firemandFile     = binDir + "/firemand"
		firemandFilePath = deployDirPath + "/firemand"

		firemanConfDir = confDir + "/fireman_conf"
		commonConfDir  = confDir + "/common_conf"
	)

	makeDirParams := executor.ExecutorCmdParams{
		oscmd.CmdParamPath: firemanConfDir,
	}
	makeDirRes := oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	makeDirParams = executor.ExecutorCmdParams{
		oscmd.CmdParamPath: commonConfDir,
	}
	makeDirRes = oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	copyParams := []executor.ExecutorCmdParams{}
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource:        firemanDirPath,
		oscmd.CmdParamTarget:        zMonDir,
		oscmd.CmdParamRecursiveCopy: true,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: firemandFilePath,
		oscmd.CmdParamTarget: binDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: deployDirPath + "/fireman_uwsgi.ini",
		oscmd.CmdParamTarget: firemanConfDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: deployDirPath + "/fireman_gun.conf",
		oscmd.CmdParamTarget: firemanConfDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: deployDirPath + "/web_service.conf",
		oscmd.CmdParamTarget: commonConfDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: toolsDirPath + "/MegaCli64",
		oscmd.CmdParamTarget: toolsDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: toolsDirPath + "/check_volume_active.py",
		oscmd.CmdParamTarget: toolsDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: deployDirPath + "/checkvolumed",
		oscmd.CmdParamTarget: "/etc/init.d",
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: toolsDirPath + "/lvm2_rhel72.tar.gz",
		oscmd.CmdParamTarget: toolsDir,
	})

	for _, cmdParams := range copyParams {
		copyRes := oscmd.Copy(e, &cmdParams)
		if !copyRes.Successful {
			return copyRes
		}
	}

	changeModeParams := executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: "/etc/init.d/checkvolumed",
	}
	changeModeRes := oscmd.ChangeMode(e, &changeModeParams)
	if !changeModeRes.Successful {
		return changeModeRes
	}

	changeModeParams = executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "777",
		oscmd.CmdParamFilenamePattern: firemandFile,
	}
	changeModeRes = oscmd.ChangeMode(e, &changeModeParams)
	if !changeModeRes.Successful {
		return changeModeRes
	}

	linkParams := executor.ExecutorCmdParams{
		oscmd.CmdParamSource: firemandFile,
		oscmd.CmdParamTarget: "/etc/init.d/firemand",
	}
	linkRes := oscmd.Link(e, &linkParams)
	if !linkRes.Successful {
		return linkRes
	}

	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "fireman installed successful")
}
