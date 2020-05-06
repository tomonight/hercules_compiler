package zdata

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"strings"
)

/*
  #ibDriverUrl --ibdriver源文件地址
  #ibDirverFileName --ibdriver文件名
  #ofedFileName -- ofed文件名
  # -- zdata配置文件路径
  execute oscmd.DownloadFile url="${ibDriverUrl}" outputFileName="/tmp/${ibDriverFilename}"
  execute oscmd.UnzipFile directory="/tmp" filename="/tmp/${ibDriverFilename}"
  set var zDataStorageDir="/zData_Storage"
  execute oscmd.UnzipFile directory="/tmp" filename="/tmp/${ibDriverFilename}/${zDataStorageDir}/${ofedFileName}

*/

//InstallIBDriver 安装 IB_Driver
func InstallIBDriver(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		cmdStr string
		es     *executor.ExecutedStatus
		err    error
		iVer   int
	)

	er := oscmd.LinuxDist(e, params)
	if !er.ExecuteHasError {
		version, ok := er.ResultData[oscmd.ResultDataKeyVersion]
		if ok {
			(*params)[oscmd.CmdParamVersion] = version
		}
	} else {
		return executor.ErrorExecuteResult(errors.New(er.Message))
	}
	//@dec 获取版本信息
	log.Debug("start get linux dist ver")
	iVer, err = oscmd.GetLinuxDistVer(params)
	if err != nil {
		log.Warn("GetLinuxDistVerfailed %v", err)
		return executor.ErrorExecuteResult(err)
	}

	//@desc 获取源文件路径
	path, err := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamPath)
	if err != nil {
		log.Warn("Get  zdata file path failed %v", err)
		return executor.ErrorExecuteResult(err)
	}

	//@desc 获取完整的文件名
	ibDriverFilename, err := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamFilename)
	if err != nil {
		log.Warn("Get  zdata file path failed %v", err)
		return executor.ErrorExecuteResult(err)
	}

	//@desc 切换目录
	cmdStr = fmt.Sprintf("cd /tmp")
	log.Debug("start change dir /tmp")
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	//@desc 解压文件
	//ibTar := "ZDATA_RHEL68_RHEL72_20170629.tar.gz"
	ibTar := path
	log.Debug("ibTart %s", ibTar)
	ibTarDir := strings.Split(ibTar, ".")[0]
	cmdStr = fmt.Sprintf("tar -xf %s -C /tmp", ibTar)
	log.Debug("start Decompress file %s", ibTar)
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	installDir := ""

	if iVer >= 7 {
		cmdStr = fmt.Sprintf("tar -xf %s/zData_Storage/%s -C /tmp", ibTarDir, ibDriverFilename)
	} else {
		cmdStr = fmt.Sprintf("tar -xf %s/zData_Computer/%s -C /tmp", ibTarDir, ibDriverFilename)
	}

	log.Debug("start Decompress file %s", ibTar)
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	//@desc 切换目录
	ibdFileName := strings.Split(ibDriverFilename, ".tar.gz")[0]
	if len(ibdFileName) == 0 {
		ibdFileName = strings.Split(ibDriverFilename, ".tgz")[0]
	}

	installDir = fmt.Sprintf("/tmp/%s/", ibdFileName)
	cmdStr = fmt.Sprintf("cd /tmp/%s", ibdFileName)
	log.Debug("start change dir /tmp/%s", ibdFileName)
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	//@desc start
	e.SetWorkingPath(installDir)
	cmdStr = "./install.pl -c zdata.conf"
	log.Debug("start ./install.pl -c zdata.conf")
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	return executor.SuccessulExecuteResult(es, true, "IB driver installed successful")
}

func getRmmodCommonError(moduleName string) []string {
	errStr := []string{fmt.Sprintf("Module %s is not currently loaded", moduleName),
		fmt.Sprintf("Module %s does not exist in /proc/modules", moduleName)}

	return errStr
}

func getRmmodError(moduleName string, err error) error {
	if err == nil {
		return nil
	}
	errStr := getRmmodCommonError(moduleName)
	for _, value := range errStr {
		if strings.Contains(err.Error(), value) {
			return nil
		}
	}
	return err
}

//modExist 判断内核是否存在
func modExist(e executor.Executor, modName string) (bool, error) {
	if e == nil {
		return false, errors.New("executor is nil")
	}
	//use lsmod |grep modName
	//cmdStr := fmt.Sprintf("lsmod |grep %s", modName)
	cmdStr := "lsmod"
	log.Debug("cmdstr = %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	log.Debug("es= %v err = %v", es, err)
	if err != nil {
		return true, err
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return true, err
	}
	log.Debug("stdout len  = %d and value  = %v", len(es.Stdout), es.Stdout)
	for index, value := range es.Stdout {
		log.Debug("stdout[%d] value=%v", index, value)

		tmpList := strings.Split(value, " ")
		log.Debug("tmList len=%d and value = %v", len(tmpList), tmpList)
	}
	if len(es.Stdout) > 0 {
		if es.Stdout[0] == modName {
			return true, nil
		}
	}

	return false, err
}

const (
	CmdParamsLinuxModName = "moduleName"
	CmdNameRemoveModule   = "removeModule"
)

//RmMod 卸载linux内核模块
func RmMpd(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	if e == nil {
		return executor.ErrorExecuteResult(errors.New("executor is nil"))
	}

	modName, err := executor.ExtractCmdFuncStringParam(params, CmdParamsLinuxModName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	exist, err := modExist(e, modName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !exist { //内核模块不存在，不执行卸载动作
		return executor.SuccessulExecuteResultNoData("mod " + modName + " not exist no need to remove")
	}

	cmdStr := fmt.Sprintf("rmmod %s", modName)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResult(es, true, "rmmod "+modName+" successful")
}

//StartIBService 启动IBService
func StartIBService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		cmdStr string
		es     *executor.ExecutedStatus
		err    error
	)
	cmdStr = "rmmod xprtrdma"
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	err = getRmmodError("xprtrdma", err)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	cmdStr = "rmmod ib_isert"
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	err = getRmmodError("ib_isert", err)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	cmdStr = "rmmod cxgb3i"
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	err = getRmmodError("cxgb3i", err)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	cmdStr = "rmmod cxgb4i"
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	err = getRmmodError("cxgb4i", err)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}

	cmdStr = "/etc/init.d/openibd start"
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		log.Warn("cmd %s failed %v", cmdStr, err)
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResult(es, true, "IB service started successful")
}
