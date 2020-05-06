package zdata

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules/osservice"
	"strings"
)

// InstallCellManager 安装cell-manager
func InstallCellManager(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
		zManDir            = homeDir + "/zman"
		binDir             = homeDir + "/bin"
		confDir            = homeDir + "/conf"
		cellManagerConfDir = confDir + "/cell_manager_conf"
		logDir             = homeDir + "/log"
		cellManagerLogDir  = logDir + "/cell_manager_log"
		cellManagerdFile   = "/etc/init.d/cell_managerd"
	)

	makeDirParams := []executor.ExecutorCmdParams{}
	makeDirParams = append(makeDirParams, executor.ExecutorCmdParams{
		oscmd.CmdParamPath: cellManagerConfDir,
	})
	makeDirParams = append(makeDirParams, executor.ExecutorCmdParams{
		oscmd.CmdParamPath: cellManagerLogDir,
	})

	for _, cmdParams := range makeDirParams {
		makeDirRes := oscmd.MakeDir(e, &cmdParams)
		if !makeDirRes.Successful {
			return makeDirRes
		}
	}

	copyParams := []executor.ExecutorCmdParams{}
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource:        installDir + "/src/cell_manager",
		oscmd.CmdParamTarget:        zManDir,
		oscmd.CmdParamRecursiveCopy: true,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: installDir + "/deploy/cell_manager.conf",
		oscmd.CmdParamTarget: cellManagerConfDir + "/cell_manager.conf",
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: installDir + "/deploy/compute_init.sh",
		oscmd.CmdParamTarget: binDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: installDir + "/deploy/storage_init.sh",
		oscmd.CmdParamTarget: binDir,
	})

	for _, cmdParams := range copyParams {
		copyRes := oscmd.Copy(e, &cmdParams)
		if !copyRes.Successful {
			return copyRes
		}
	}

	linkParams := executor.ExecutorCmdParams{
		oscmd.CmdParamSource: zManDir + "/cell_manager/cell_manager.sh",
		oscmd.CmdParamTarget: cellManagerdFile,
	}
	linkRes := oscmd.Link(e, &linkParams)
	if !linkRes.Successful {
		return linkRes
	}

	textToFileParams := executor.ExecutorCmdParams{
		oscmd.CmdParamFilename: cellManagerLogDir + "/asm_check.log",
		oscmd.CmdParamOutText:  "",
	}
	textToFileRes := oscmd.TextToFile(e, &textToFileParams)
	if !textToFileRes.Successful {
		return textToFileRes
	}

	changeModeParams := []executor.ExecutorCmdParams{}
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: cellManagerdFile,
	})
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "666",
		oscmd.CmdParamFilenamePattern: cellManagerLogDir + "/asm_check.log",
	})

	for _, cmdParams := range changeModeParams {
		changeModeRes := oscmd.ChangeMode(e, &cmdParams)
		if !changeModeRes.Successful {
			return changeModeRes
		}
	}

	serviceCtrlParams := executor.ExecutorCmdParams{
		osservice.CmdParamServiceName:   "cell_managerd",
		osservice.CmdParamServiceAction: "all boot_disable",
	}
	serviceCtrlRes := osservice.SysVServiceControl(e, &serviceCtrlParams)
	if !serviceCtrlRes.Successful {
		return serviceCtrlRes
	}

	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "cell_manager installed successful")
}
