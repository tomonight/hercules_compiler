package zdata

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	//	"hercules_compiler/rce-executor/utils"
	"strings"
)

// InstallZcli 安装zcli服务
func InstallZcli(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
		zcliLinkPath             = homeDir + "/zcli"
		zcliDir                  = homeDir + "/zman/zcli"
		zcliFile                 = zcliDir + "/zcli.py"
		zcliFilePath             = installDir + "/src/zcli/zcli.py"
		zcliConfDir              = homeDir + "/conf/comon_conf"
		zcliConfFile             = zcliConfDir + "/web_service.conf"
		zcliConfFilePath         = installDir + "/deploy/web_service.conf"
		userdefFlashConfDir      = homeDir + "/conf/userdef_flash_card"
		userdefFlashConfFile     = userdefFlashConfDir + "/userdef_flash_card.conf"
		userdefFlashConfFilePath = installDir + "/deploy/userdef_flash_card.conf"
	)
	log.Debug("zcliFilePath %s", zcliFilePath)
	makeDirParams := executor.ExecutorCmdParams{
		oscmd.CmdParamPath: zcliConfDir,
	}
	makeDirRes := oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	//zcliConfFileExist := utils.PathExist(zcliConfFile)
	zcliConfFileExist, _ := fileExist(e, zcliConfFile)
	if !zcliConfFileExist {
		copyParams := executor.ExecutorCmdParams{
			oscmd.CmdParamSource: zcliConfFilePath,
			oscmd.CmdParamTarget: zcliConfDir,
		}
		copyRes := oscmd.Copy(e, &copyParams)
		if !copyRes.Successful {
			return copyRes
		}
	}

	makeDirParams = executor.ExecutorCmdParams{
		oscmd.CmdParamPath: zcliDir,
	}
	makeDirRes = oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	//zcliFileExist := utils.PathExist(zcliFile)
	zcliFileExist, _ := fileExist(e, zcliFile)
	log.Debug("zcliFileExist: %d", zcliConfFileExist)
	if !zcliFileExist {
		copyParams := executor.ExecutorCmdParams{
			oscmd.CmdParamSource: zcliFilePath,
			oscmd.CmdParamTarget: zcliDir,
		}
		log.Debug("do copy file %s", zcliFile)
		copyRes := oscmd.Copy(e, &copyParams)
		if !copyRes.Successful {
			return copyRes
		}
	}

	changeModeParams := executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: zcliFile,
	}
	changeModeRes := oscmd.ChangeMode(e, &changeModeParams)
	if !changeModeRes.Successful {
		return changeModeRes
	}

	linkParams := executor.ExecutorCmdParams{
		oscmd.CmdParamSource: zcliFile,
		oscmd.CmdParamTarget: zcliLinkPath,
	}
	linkRes := oscmd.Link(e, &linkParams)
	if !linkRes.Successful {
		return linkRes
	}

	makeDirParams = executor.ExecutorCmdParams{
		oscmd.CmdParamPath: userdefFlashConfDir,
	}
	makeDirRes = oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	// 注意utils.PathExist此函数引起严重bug 远程文件不能在本地判断
	//userdefFlashConfFileExist := utils.PathExist(userdefFlashConfFile)
	userdefFlashConfFileExist, _ := fileExist(e, userdefFlashConfFile)
	if !userdefFlashConfFileExist {
		copyParams := executor.ExecutorCmdParams{
			oscmd.CmdParamSource: userdefFlashConfFilePath,
			oscmd.CmdParamTarget: userdefFlashConfDir,
		}
		copyRes := oscmd.Copy(e, &copyParams)
		if !copyRes.Successful {
			return copyRes
		}
	}
	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "zcli installed successful")
}
