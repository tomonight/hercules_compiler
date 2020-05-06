//封装mysql备份命令
package mysql

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"path/filepath"
	"strconv"
	"strings"
)

// 函数名常量定义
const (
	CmdNameRecvBackupDataBySocat          = "RecvBackupDataBySocat"
	CmdNameSendBackupDataBySocat          = "SendBackupDataBySocat"
	CmdNameRecoverBackupData              = "RecoverBackupData"
	CmdNameCheckBackupData                = "CheckBackupData"
	CmdNameRecoverBackUpFiles             = "RecoverBackUpFiles"
	CmdNameRecoverOneIncrementBackUpData  = "RecoverOneIncrementBackUpData"
	CmdNameRenameSchema                   = "RenameSchema"
	CmdNameLogicalRecover                 = "LogicalRecover"
	CmdNameLogicalBackup                  = "LogicalBackup"
	CmdNameBinlogRecover                  = "BinlogRecover"
	CmdNameInterruptDatabaseRemoteService = "InterruptDatabaseRemoteService"
	CmdNameRecoverDatabaseRemoteService   = "RecoverDatabaseRemoteService"
)

// 命令参数常量定义
const (
	CmdParamSocatPath        = "socatPath"        //socat 路径
	CmdParamBytesPerTransfer = "bytesPerTransfer" //每次传送字节数
	CmdParamListenPort       = "listenPort"       //监听端口
	CmdParamXtrabackupPath   = "xtrabackupPath"   //xtrabackup软件安装地址
	CmdParamBackupUser       = "backupUser"       //备份用户名
	CmdParamBackupPassword   = "backupPassword"   //备份密码
	CmdParamRecoverHost      = "recoverHost"      //恢复主机host
	CmdParambackUpType       = "backUpType"       //备份类型
	CmdParambackUpLsn        = "lsn"              //备份LSN
	CmdParamHasIncrement     = "hasIncrement"     //是否包含增量备份文件
	CmdParamBackupDataPath   = "backupDataPath"   //备份数据目录
	CmdParamBackupDataSize   = "backupDataSize"   //备份文件大小CmdParamBackupDataPath
)

//定义默认每次传送字节数
const (
	DefaultBytesPerTransfer = 10485760
)

//RecvBackupDataBySocat 通过socat接收备份数据
func RecvBackupDataBySocat(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var socatCmd, xtrabackupCmd string

	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xbstream")

	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(false, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	socatPath, _ := executor.ExtractCmdFuncStringParam(params, CmdParamSocatPath)
	if socatPath != "" {
		socatCmd = filepath.Join(socatPath, "socat")
	} else {
		socatCmd = "socat"
	}
	//指定接收网卡，做分流处理
	recoverHost, err := executor.ExtractCmdFuncStringParam(params, "recoverHost")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	bytesPerTransfer, _ := executor.ExtractCmdFuncIntParam(params, CmdParamBytesPerTransfer)
	if bytesPerTransfer <= 0 {
		bytesPerTransfer = DefaultBytesPerTransfer
	}

	listenPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamListenPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	dataDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlDataDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//socat -b 10485760 -u TCP-LISTEN:${backListenPort},reuseaddr stdio | /usr/local/xtrabackup/bin/xbstream -xv
	cmdStr := fmt.Sprintf("%s -b %d -u TCP-LISTEN:%d,bind=%s,reuseaddr,rcvtimeo=60 stdio | %s -xv -C %s",
		socatCmd, bytesPerTransfer, listenPort, recoverHost, xtrabackupCmd, dataDir)
	log.Debug("execute command %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.SuccessulExecuteResult(es, true, "recv backup data by socat successful ")
}

//SendBackupDataBySocat 通过socat传送备份数据
func SendBackupDataBySocat(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var socatCmd, xtrabackupCmd string

	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")

	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(true, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	socatPath, _ := executor.ExtractCmdFuncStringParam(params, CmdParamSocatPath)
	if socatPath != "" {
		socatCmd = filepath.Join(socatPath, "socat")
	} else {
		socatCmd = "socat"
	}

	bytesPerTransfer, _ := executor.ExtractCmdFuncIntParam(params, CmdParamBytesPerTransfer)
	if bytesPerTransfer <= 0 {
		bytesPerTransfer = DefaultBytesPerTransfer
	}

	listenPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamListenPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	backUpType, err := executor.ExtractCmdFuncStringParam(params, CmdParambackUpType)
	if err != nil {
		backUpType = "1"
	}

	recoverHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamRecoverHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPort, err := executor.ExtractCmdFuncStringParam(params, CmdParamMySQLPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	lsn, _ := executor.ExtractCmdFuncStringParam(params, CmdParambackUpLsn)
	if lsn == "0" && backUpType == "2" {
		return executor.ErrorExecuteResult(errors.New("lsn can not be an empty string"))
	}

	paramFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlParameterFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	bkUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamBackupUser)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	bkPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamBackupPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	incrStr := ""
	if backUpType == "2" {
		incrStr = fmt.Sprintf(" --incremental-lsn=%s", lsn)
	}

	//oldversion  usr/local/xtrabackup/bin/xtrabackup --defaults-file=${masterMysqlParameterFile} --backup --stream=xbstream --socket=${masterSocketFile} --user=bkuser --password=123456 | socat -b 10485760 -u stdio TCP:${slaveHost}:${backListenPort}"
	//newversion  xtrabackup --defaults-file=/home/zmysql/db/mgr_tigten/mgr_tigten01/conf/my.cnf  --backup --stream=xbstream --host=127.0.0.1 --port=8454 --user=bkuser --password=123456 | socat -b 10485760 -u stdio TCP:192.168.11.189:13558
	cmdStr := fmt.Sprintf("%s --defaults-file=%s --backup --stream=xbstream --host=127.0.0.1 --port=%s --user=%s --password=%s %s | %s -b %d -u stdio TCP:%s:%d",
		xtrabackupCmd, paramFile, mysqlPort, bkUser, bkPassword, incrStr, socatCmd, bytesPerTransfer, recoverHost, listenPort)
	log.Debug("execute command = %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	//避免未备份成功，然后未返回错误的情况
	completeFlag := false
	if len(es.Stderr) > 0 && strings.Contains(es.Stderr[len(es.Stderr)-1], "completed OK!") {
		completeFlag = true
	}
	if !completeFlag {
		errMsg := strings.Join(es.Stderr, "\n")
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.SuccessulExecuteResult(es, true, "send backup data by socat successful ")
}

//RecoverBackupData 恢复备份数据
func RecoverBackupData(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var xtrabackupCmd string
	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")

	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(true, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	dataDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlDataDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// /usr/local/xtrabackup/bin/xtrabackup --apply-log ${slaveDataBaseDir}/data
	// /usr/local/xtrabackup/bin/xtrabackup --prepare --target-dir=${slaveDataBaseDir}/data
	// cmdStr := fmt.Sprintf("%s --apply-log %s", xtrabackupCmd, dataDir)
	cmdStr := fmt.Sprintf("%s --prepare --target-dir=%s", xtrabackupCmd, dataDir)
	log.Debug("execute command %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.SuccessulExecuteResult(es, true, "recover backup data successful ")
}

//CheckBackupData 检查备份数据是否有效
func CheckBackupData(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//取出辈分文件的路径
	backupTaskID, err := executor.ExtractCmdFuncStringParam(params, CmdParambackUpType)
	if err != nil {
		log.Error(err.Error())
		return executor.ErrorExecuteResult(err)
	}
	result := checkData(e, params)
	es, err := e.ExecShell("")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	er := executor.SuccessulExecuteResult(es, true, "backup data enable")

	er.ResultData = make(map[string]string)
	er.ResultData["status"] = fmt.Sprintf("%t", result)
	er.ResultData["backupID"] = backupTaskID
	return er
}

func checkData(e executor.Executor, params *executor.ExecutorCmdParams) (result bool) {
	backupDataPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamBackupDataPath)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	//检查备份文件是否存在
	log.Debug("check path is %s", backupDataPath)
	exsit, err := oscmd.PathExist(e, backupDataPath)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if !exsit {
		err = errors.New("file not exsit")
		log.Error(err.Error())
		return false
	}

	//检查备份文件大小
	//取出备份文件的大小
	log.Debug("get backup file size")
	backupDataSize, err := executor.ExtractCmdFuncIntParam(params, CmdParamBackupDataSize)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	res, err := e.ExecShell(fmt.Sprintf("du -s %s", backupDataPath))
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if res.ExitCode != 0 {
		var errMsg string
		if len(res.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(res.Stderr, "\n")
		}
		log.Error(errMsg)
		return false
	}

	value := strings.Split(res.Stdout[0], "\t")
	num, err := strconv.Atoi(value[0])
	if err != nil {
		log.Error(err.Error())
		return false
	}
	log.Debug("database recode size is %d and check size is %d", num, backupDataSize)
	if num < backupDataSize {
		return false
	}
	return true
}

//RecoverBackUpFiles 修复备份文件，避免恢复数据库失败
func RecoverBackUpFiles(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	backupDataPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamBackupDataPath)
	if err != nil {
		log.Error(err.Error())
		return executor.ErrorExecuteResult(err)
	}

	//检查备份文件是否存在
	log.Debug("check path is %s", backupDataPath)
	exist, err := oscmd.PathExist(e, backupDataPath)
	if err != nil {
		log.Error(err.Error())
		return executor.ErrorExecuteResult(err)
	}
	if !exist {
		err = errors.New("file not exist")
		log.Error(err.Error())
		return executor.ErrorExecuteResult(err)
	}

	changeDirCommand := fmt.Sprintf("cd %s", backupDataPath)

	//binlog_name=$(awk '{print $1}' xtrabackup_binlog_info)
	getBinLogCommand := `binlog_name=$(awk '{print $1}' xtrabackup_binlog_info)`

	//binlog_index_name=$(echo $binlog_name|awk -F '.' '{print $1}')'.index'
	getIndexNameCommand := `binlog_index_name=$(echo $binlog_name|awk -F '.' '{print $1}')'.index'`

	//echo './'${binlog_name} > ${binlog_index_name}
	resetIndexFileContentCommand := `echo './'${binlog_name} > ${binlog_index_name}`

	command := fmt.Sprintf("%s;%s;%s;%s", changeDirCommand, getBinLogCommand, getIndexNameCommand, resetIndexFileContentCommand)
	es, err := e.ExecShell(command)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("recover backup data files %s successful", backupDataPath))
}
