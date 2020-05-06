package zdata

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"strings"
)

// InstallZMonAgent 安装zmon-agent
func InstallZMonAgent(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
		zMonDir      = homeDir + "/zmon"
		agentDir     = zMonDir + "/zmon_agent"
		agentSrcPath = installDir + "/src/zmon_agent/"

		agentConfDir = homeDir + "/conf/zmon_agent_conf"
		//agentConfFilePath = installDir + "/deploy/zmon_server_ipaddr_port.xml"

		agentLogDir = homeDir + "/log/zmon_agent_log"

		agentManageFile     = "/etc/init.d/zmonagentd"
		agentManageFilePath = installDir + "/deploy/zmonagentd"

		etcdConfDir = homeDir + "/conf/etcd_conf"
		//		etcdConfFile     = etcdConfDir + "/etcd_service.conf"
		//		etcdConfFilePath = installDir + "/deploy/etcd_service.conf.example"

		toolsDir          = homeDir + "/tools"
		etcdCheckFile     = toolsDir + "/check_etcd_state.py"
		etcdCheckFilePath = installDir + "/tools/check_etcd_state.py"

		itemsDir  = agentDir + "/items"
		sqliteDir = agentDir + "/sqlite_db"
	)

	makeDirParams := []executor.ExecutorCmdParams{}
	makeDirParams = append(makeDirParams, executor.ExecutorCmdParams{
		oscmd.CmdParamPath: agentLogDir,
	})
	makeDirParams = append(makeDirParams, executor.ExecutorCmdParams{
		oscmd.CmdParamPath: agentConfDir,
	})
	makeDirParams = append(makeDirParams, executor.ExecutorCmdParams{
		oscmd.CmdParamPath: etcdConfDir,
	})
	makeDirParams = append(makeDirParams, executor.ExecutorCmdParams{
		oscmd.CmdParamPath: zMonDir,
	})

	for _, cmdParams := range makeDirParams {
		makeDirRes := oscmd.MakeDir(e, &cmdParams)
		if !makeDirRes.Successful {
			return makeDirRes
		}
	}

	copyParams := []executor.ExecutorCmdParams{}
	//	copyParams = append(copyParams, executor.ExecutorCmdParams{
	//		oscmd.CmdParamSource: agentConfFilePath,
	//		oscmd.CmdParamTarget: agentConfDir,
	//	})
	//	copyParams = append(copyParams, executor.ExecutorCmdParams{
	//		oscmd.CmdParamSource: etcdConfFilePath,
	//		oscmd.CmdParamTarget: etcdConfFile,
	//	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource:        agentSrcPath,
		oscmd.CmdParamTarget:        zMonDir,
		oscmd.CmdParamRecursiveCopy: true,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: etcdCheckFilePath,
		oscmd.CmdParamTarget: etcdCheckFile,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: installDir + "/tools/ssacli",
		oscmd.CmdParamTarget: toolsDir,
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: installDir + "/deploy/zdata_logrotate.conf",
		oscmd.CmdParamTarget: homeDir + "/conf/",
	})
	copyParams = append(copyParams, executor.ExecutorCmdParams{
		oscmd.CmdParamSource: agentManageFilePath,
		oscmd.CmdParamTarget: agentManageFile,
	})

	for _, cmdParams := range copyParams {
		copyRes := oscmd.Copy(e, &cmdParams)
		if !copyRes.Successful {
			return copyRes
		}
	}

	changeModeParams := []executor.ExecutorCmdParams{}
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: agentManageFile,
	})
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: itemsDir + "/db_operate.py",
	})
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: itemsDir + "/db_operate_asm.py",
	})
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "755",
		oscmd.CmdParamFilenamePattern: itemsDir + "/db_operate_asm_dg_attr.py",
	})
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "777",
		oscmd.CmdParamFilenamePattern: sqliteDir,
	})
	changeModeParams = append(changeModeParams, executor.ExecutorCmdParams{
		oscmd.CmdParamModeExp:         "777",
		oscmd.CmdParamFilenamePattern: sqliteDir + "/asm.db",
	})

	for _, cmdParams := range changeModeParams {
		changeModeRes := oscmd.ChangeMode(e, &cmdParams)
		if !changeModeRes.Successful {
			return changeModeRes
		}
	}

	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "zmon-agent installed successful")
}
