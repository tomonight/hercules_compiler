package zdata

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules/osconfig"
	"hercules_compiler/rce-executor/modules/osservice"
	"strconv"
	"strings"
)

// InitDirs 初始化目录
func InitDirs(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	homeDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamHomeDir)
	if err != nil {
		homeDir = defaultHomeDir
	}
	homeDir = strings.TrimRight(homeDir, "/")

	var (
		zManDir       = homeDir + "/zman"
		zMonDir       = homeDir + "/zmon"
		toolsDir      = homeDir + "/tools"
		docDir        = homeDir + "/doc"
		binDir        = homeDir + "/bin"
		confDir       = homeDir + "/conf"
		logDir        = homeDir + "/log"
		runDir        = homeDir + "/run"
		userScriptDir = homeDir + "/user_script"
		pyDir         = homeDir + "/python"
	)

	makeDirParams := []executor.ExecutorCmdParams{}
	makeDirParams = append(
		makeDirParams,
		executor.ExecutorCmdParams{oscmd.CmdParamPath: zManDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: zMonDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: toolsDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: docDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: binDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: confDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: logDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: runDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: userScriptDir},
		executor.ExecutorCmdParams{oscmd.CmdParamPath: pyDir},
	)
	for _, cmdParams := range makeDirParams {

		makeDirRes := oscmd.MakeDir(e, &cmdParams)
		if !makeDirRes.Successful {
			return makeDirRes
		}
	}
	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "zdata directories initialized successful")
}

func unInstallZData(e executor.Executor) error {
	cmdstr := "rpm -qa | egrep ^zdata"
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		return err
	}

	for _, value := range es.Stdout {
		err = rmpUnInstall(e, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func rmpUnInstall(e executor.Executor, moduleName string) error {
	cmdstr := fmt.Sprintf("rpm -e %s", moduleName)
	log.Debug("rmp uninstall module command = %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return err
	}
	err = executor.GetExecResult(es)
	return err
}

// InstallService 安装zdata服务
func InstallService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	// 首先检查是否已安装
	cmdstr := "rpm -qa | egrep ^zdata | wc -l"
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		count, _ := strconv.Atoi(es.Stdout[0])
		if count != 0 {
			err = unInstallZData(e)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			//return executor.NotSuccessulExecuteResult(es, "zdata already installed")
		}
		// zdata not installed
		downloadRes := oscmd.DownloadFile(e, params)
		outputFilename := downloadRes.ResultData[oscmd.ResultDataKeyoutputFilename]
		cmdstr := fmt.Sprintf("rpm -ivh %s", outputFilename)
		es, err := e.ExecShell(cmdstr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if es.ExitCode == 0 && len(es.Stderr) == 0 {
			return executor.SuccessulExecuteResult(es, true, "zdata installed successful")
		}
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// InstallCxOracle 安装cx_Oracle
func InstallCxOracle(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
		softDir               = "/soft"
		pyDir                 = homeDir + "/python"
		pyExec                = pyDir + "/bin/python"
		libDir                = homeDir + "/lib"
		libTarFile            = softDir + "/lib.tgz"
		cxOracleDir           = softDir + "/cx_Oracle-6.0b1"
		cxOracleTarFile       = softDir + "/cx_Oracle-6.0b1.tar.gz"
		cxOracleSetUpFilePath = cxOracleDir + "/setup.py"
	)

	unzipParams := executor.ExecutorCmdParams{
		oscmd.CmdParamFilename:  libTarFile,
		oscmd.CmdParamDirectory: softDir,
	}
	unzipRes := oscmd.UnzipFile(e, &unzipParams)
	log.Debug("unzipRes:%v", unzipRes)
	if !unzipRes.Successful {
		return unzipRes
	}

	unzipParams = executor.ExecutorCmdParams{
		oscmd.CmdParamFilename:  cxOracleTarFile,
		oscmd.CmdParamDirectory: softDir,
	}
	unzipRes = oscmd.UnzipFile(e, &unzipParams)
	log.Debug("unzipRes:%v", unzipRes)
	if !unzipRes.Successful {
		return unzipRes
	}

	envStr := fmt.Sprintf("export LD_LIBRARY_PATH=$ORACLE_HOME/lib:/lib:/usr/lib:%s/oracle/11.2/client64/lib", libDir)
	setSysConfParams := executor.ExecutorCmdParams{
		osconfig.CmdParamConfFile: "/etc/profile",
		osconfig.CmdParamConfItem: envStr,
	}
	setSysConfRes := osconfig.SetSysConf(e, &setSysConfParams)
	if !setSysConfRes.Successful {
		return setSysConfRes
	}

	cmdstr := fmt.Sprintf("cd %s;%s %s install", cxOracleDir, pyExec, cxOracleSetUpFilePath)
	log.Debug("cmdstr :%s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	log.Debug("exit code %d", es.ExitCode)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 /*&& len(es.Stderr) == 0*/ {
		return executor.SuccessulExecuteResult(es, true, "cx_Oracle installed successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

func getInstallError(errMsg string) error {
	ignoreMsg := []string{"already installed", "warning"}
	for _, msg := range ignoreMsg {
		if strings.Contains(errMsg, msg) {
			return nil
		}
	}
	return errors.New(errMsg)
}

// InitStorageNode 初始化存储节点
func InitStorageNode(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
		toolsDir         = homeDir + "/tools"
		lvmDir           = toolsDir + "/lvm2_rhel72"
		lvmTarFilePath   = toolsDir + "/lvm2_rhel72.tar.gz"
		binDir           = homeDir + "/bin"
		firemandFile     = "/etc/init.d/firemand"
		cellManagerdFile = "/etc/init.d/cell_managerd"
		zMonAgentdFile   = "/etc/init.d/zmonagentd"
		zManDir          = homeDir + "/zman"
	)

	unzipParams := executor.ExecutorCmdParams{
		oscmd.CmdParamFilename:  lvmTarFilePath,
		oscmd.CmdParamDirectory: toolsDir,
	}
	unzipRes := oscmd.UnzipFile(e, &unzipParams)
	if !unzipRes.Successful {
		return unzipRes
	}
	//忽略签名
	e.ExecShell("rpm --import /etc/pki/rpm-gpg/RPM* ")
	pkgs := []string{
		"lvm2-python-libs-2.02.166-1.el7_3.1.x86_64.rpm",
		"lvm2-libs-2.02.166-1.el7_3.1.x86_64.rpm",
		"dm/device-mapper-event-libs-1.02.135-1.el7_3.1.x86_64.rpm",
		"dm/device-mapper-event-1.02.135-1.el7_3.1.x86_64.rpm",
		"dm/device-mapper-1.02.135-1.el7_3.1.x86_64.rpm",
		"dm/device-mapper-libs-1.02.135-1.el7_3.1.x86_64.rpm",
		"lvm2-libs-2.02.166-1.el7_3.1.x86_64.rpm",
		"lvm2-2.02.166-1.el7_3.1.x86_64.rpm",
		"dm/device-mapper-persistent-data-0.6.3-1.el7.x86_64.rpm",
	}
	pkgNames := []string{}
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, lvmDir+"/"+pkg)
	}
	rpmInstallParams := executor.ExecutorCmdParams{
		oscmd.CmdParamFilename:    strings.Join(pkgNames, " "),
		oscmd.CmdParamInstallFlag: " -Uih",
	}
	rpmInstallRes := oscmd.RpmInstall(e, &rpmInstallParams)
	log.Debug("install rpm exitcode = %d", rpmInstallRes.ExitCode)
	if !rpmInstallRes.Successful {
		log.Debug("rmpInstallRes = %v", rpmInstallRes)
		err := getInstallError(rpmInstallRes.Message)
		if err != nil {
			return rpmInstallRes
		}
	}

	linkParams := []executor.ExecutorCmdParams{}
	linkParams = append(
		linkParams,
		executor.ExecutorCmdParams{
			oscmd.CmdParamSource: binDir + "/firemand",
			oscmd.CmdParamTarget: firemandFile,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamSource: zManDir + "/zcli/zcli.py",
			oscmd.CmdParamTarget: homeDir + "/zcli",
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamSource: zManDir + "/cell_manager/cell_manager.sh",
			oscmd.CmdParamTarget: cellManagerdFile,
		},
	)

	for _, cmdParams := range linkParams {
		linkRes := oscmd.Link(e, &cmdParams)
		if !linkRes.Successful {
			return linkRes
		}
	}

	copyParams := executor.ExecutorCmdParams{
		oscmd.CmdParamSource: homeDir + "/zmon/zmon_agent/zmonagentd",
		oscmd.CmdParamTarget: zMonAgentdFile,
	}
	copyRes := oscmd.Copy(e, &copyParams)
	if !copyRes.Successful {
		return copyRes
	}

	changeModeParams := []executor.ExecutorCmdParams{}
	changeModeParams = append(
		changeModeParams,
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "755",
			oscmd.CmdParamFilenamePattern: firemandFile,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "755",
			oscmd.CmdParamFilenamePattern: zMonAgentdFile,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "755",
			oscmd.CmdParamFilenamePattern: cellManagerdFile,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "+x",
			oscmd.CmdParamFilenamePattern: "/etc/rc.local",
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "+x",
			oscmd.CmdParamFilenamePattern: "/etc/rc.d/rc.local",
		},
	)

	for _, cmdParams := range changeModeParams {
		changeModeRes := oscmd.ChangeMode(e, &cmdParams)
		if !changeModeRes.Successful {
			return changeModeRes
		}
	}

	textToFileParams := executor.ExecutorCmdParams{
		oscmd.CmdParamFilename:  "/etc/rc.d/rc.local",
		oscmd.CmdParamOutText:   "nohup /etc/init.d/cell_managerd storage_boot &",
		oscmd.CmdParamOverwrite: false,
	}
	textToFileRes := oscmd.TextToFile(e, &textToFileParams)
	if !textToFileRes.Successful {
		return textToFileRes
	}

	serviceCtlParams := executor.ExecutorCmdParams{
		osservice.CmdParamServiceName:   "cell_managerd",
		osservice.CmdParamServiceAction: "zmon_agent boot_enable",
	}
	serviceCtlRes := osservice.SysVServiceControl(e, &serviceCtlParams)
	if !serviceCtlRes.Successful {
		return serviceCtlRes
	}

	serviceCtlParams = executor.ExecutorCmdParams{
		osservice.CmdParamServiceName:   "cell_managerd",
		osservice.CmdParamServiceAction: "zman_agent boot_enable",
	}
	serviceCtlRes = osservice.SysVServiceControl(e, &serviceCtlParams)
	if !serviceCtlRes.Successful {
		return serviceCtlRes
	}

	makeDirParams := executor.ExecutorCmdParams{
		oscmd.CmdParamPath: homeDir + "/log/cell_manager_log",
	}
	makeDirRes := oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	copyParams = executor.ExecutorCmdParams{
		oscmd.CmdParamSource: homeDir + "/conf/zdata_logrotate.conf",
		oscmd.CmdParamTarget: "/etc/logrotate.d/zdata",
	}
	copyRes = oscmd.Copy(e, &copyParams)
	if !copyRes.Successful {
		return copyRes
	}

	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "storage node initialized successful")
}

// InitComputeNode 初始化计算节点
func InitComputeNode(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
		libDir              = homeDir + "/lib"
		zMonAgentFile       = "/etc/init.d/zmonagentd"
		zMonAgentFilePath   = homeDir + "/zmon/zmon_agent/zmonagentd"
		cellManagerFile     = "/etc/init.d/cell_managerd"
		cellManagerFilePath = homeDir + "/zman/cell_manager/cell_manager.sh"
		cellManagerLogDir   = homeDir + "/log/cell_manager_log"
	)

	envStr := fmt.Sprintf("export LD_LIBRARY_PATH=$ORACLE_HOME/lib:/lib:/usr/lib:%s/oracle/11.2/client64/lib", libDir)
	setSysConfParams := executor.ExecutorCmdParams{
		osconfig.CmdParamConfFile: "/etc/profile",
		osconfig.CmdParamConfItem: envStr,
	}
	setSysConfRes := osconfig.SetSysConf(e, &setSysConfParams)
	if !setSysConfRes.Successful {
		return setSysConfRes
	}

	// 需要安装oracle数据库
	// setSysConfParams = executor.ExecutorCmdParams{
	// 	osconfig.CmdParamConfFile: "/home/grid/.bash_profile",
	// 	osconfig.CmdParamConfItem: envStr,
	// }
	// setSysConfRes = osconfig.SetSysConf(e, &setSysConfParams)
	// if !setSysConfRes.Successful {
	// 	return setSysConfRes
	// }

	copyParams := executor.ExecutorCmdParams{
		oscmd.CmdParamSource: zMonAgentFilePath,
		oscmd.CmdParamTarget: zMonAgentFile,
	}
	copyRes := oscmd.Copy(e, &copyParams)
	if !copyRes.Successful {
		return copyRes
	}

	linkPrams := executor.ExecutorCmdParams{
		oscmd.CmdParamSource: cellManagerFilePath,
		oscmd.CmdParamTarget: cellManagerFile,
	}
	linkRes := oscmd.Link(e, &linkPrams)
	if !linkRes.Successful {
		return linkRes
	}

	text := fmt.Sprintf("nohup %s compute_boot &", cellManagerFile)
	textToFileParams := executor.ExecutorCmdParams{
		oscmd.CmdParamFilename: "/etc/rc.d/rc.local",
		oscmd.CmdParamOutText:  text,
	}
	textToFileRes := oscmd.TextToFile(e, &textToFileParams)
	if !textToFileRes.Successful {
		return textToFileRes
	}

	serviceCtlParams := executor.ExecutorCmdParams{
		osservice.CmdParamServiceName:   "cell_managerd",
		osservice.CmdParamServiceAction: "zmon_agent boot_enable",
	}
	serviceCtlRes := osservice.SysVServiceControl(e, &serviceCtlParams)
	if !serviceCtlRes.Successful {
		return serviceCtlRes
	}

	makeDirParams := executor.ExecutorCmdParams{
		oscmd.CmdParamPath: cellManagerLogDir,
	}
	makeDirRes := oscmd.MakeDir(e, &makeDirParams)
	if !makeDirRes.Successful {
		return makeDirRes
	}

	touchParams := executor.ExecutorCmdParams{
		oscmd.CmdParamPath: cellManagerLogDir + "/asm_check.log",
	}
	touchRes := oscmd.Touch(e, &touchParams)
	if !touchRes.Successful {
		return touchRes
	}

	touchParams = executor.ExecutorCmdParams{
		oscmd.CmdParamPath: cellManagerLogDir + "/auto_multipath.log",
	}
	touchRes = oscmd.Touch(e, &touchParams)
	if !touchRes.Successful {
		return touchRes
	}

	copyParams = executor.ExecutorCmdParams{
		oscmd.CmdParamSource: homeDir + "/conf/zdata_logrotate.conf",
		oscmd.CmdParamTarget: "/etc/logrotate.d/zdata",
	}
	copyRes = oscmd.Copy(e, &copyParams)
	if !copyRes.Successful {
		return copyRes
	}

	changeModeParams := []executor.ExecutorCmdParams{}
	changeModeParams = append(
		changeModeParams,
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "755",
			oscmd.CmdParamFilenamePattern: zMonAgentFile,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "755",
			oscmd.CmdParamFilenamePattern: cellManagerFile,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "777",
			oscmd.CmdParamFilenamePattern: cellManagerLogDir,
			oscmd.CmdParamRecursiveChange: true,
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "+x",
			oscmd.CmdParamFilenamePattern: "/etc/rc.local",
		},
		executor.ExecutorCmdParams{
			oscmd.CmdParamModeExp:         "+x",
			oscmd.CmdParamFilenamePattern: "/etc/rc.d/rc.local",
		},
	)

	for _, cmdParams := range changeModeParams {
		changeModeRes := oscmd.ChangeMode(e, &cmdParams)
		if !changeModeRes.Successful {
			return changeModeRes
		}
	}

	return executor.SuccessulExecuteResult(&executor.ExecutedStatus{}, true, "compute node initialized successful")
}
