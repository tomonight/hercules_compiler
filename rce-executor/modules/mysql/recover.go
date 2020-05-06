//recover.go 用于mysql数据恢复所有操作

package mysql

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/utils"
	"path/filepath"
	"sort"
	"strings"
)

//RecoverCompleteBackUpData 恢复全量备份数据
func RecoverCompleteBackUpData(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	hasIncrement, _ := executor.ExtractCmdFuncBoolParam(params, CmdParamHasIncrement)

	//var xtrabackupCmd string
	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")

	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(true, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	completeBackUpPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	completeBackUpDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	completeBackUpPath = filepath.Join(completeBackUpDir, utils.GetFileName(completeBackUpPath))
	confFile := filepath.Join(completeBackUpPath, "backup-my.cnf")

	// xtrabackup --defaults-file=/data/zmysql/192.168.12.200_3306/testdb/testdb1/FullBackup/20181115231000/backup-my.cnf
	// --prepare --apply-log-only --target-dir=/data/zmysql/192.168.12.200_3306/testdb/testdb1/FullBackup/20181115231000
	var cmdStr string
	if hasIncrement {
		cmdStr = fmt.Sprintf("%s --defaults-file=%s --prepare --apply-log-only --target-dir=%s", xtrabackupCmd, confFile, completeBackUpPath)
	} else {
		cmdStr = fmt.Sprintf("%s --defaults-file=%s --prepare --target-dir=%s", xtrabackupCmd, confFile, completeBackUpPath)
	}
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

	binLogPath := filepath.Join(completeBackUpPath, "xtrabackup_binlog_info")
	er := executor.SuccessulExecuteResultNoData("recover backup data successful ")
	er.ResultData = map[string]string{
		"lastLogBinPath": binLogPath,
	}
	return er
}

//RecoverIncrementBackUpData 恢复增量备份数据
func RecoverIncrementBackUpData(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	hasIncrement, err := executor.ExtractCmdFuncBoolParam(params, CmdParamHasIncrement)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if !hasIncrement {
		return executor.SuccessulExecuteResultNoData("no increment backup data no need to recover ")
	}

	//var xtrabackupCmd string
	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")
	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(true, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	completeBackUpPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	completeBackUpDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	completeBackUpPath = filepath.Join(completeBackUpDir, utils.GetFileName(completeBackUpPath))

	incrementBackUpPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamIncrementBackUpPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	incrementBackUpDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamIncrementBackUpDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	incrementBackUpPathSlice := strings.Split(incrementBackUpPath, ",")
	incrementBackUpPathList := BackUpDataPathSlice{}
	for _, path := range incrementBackUpPathSlice {
		lastFilePath := utils.GetFileName(path)
		incrementBackUpPathList = append(incrementBackUpPathList, lastFilePath)
	}
	log.Debug("incrementBackUpPathList = %v", incrementBackUpPathList)
	incrementBackUpPathValid := incrementBackUpPathList.Valid()
	if !incrementBackUpPathValid {
		return executor.ErrorExecuteResult(fmt.Errorf("增量备份数据文件格式不合法"))
	}

	//进行从小到大的排序
	sort.Sort(incrementBackUpPathList)
	log.Debug("order incrementBackUpPathList = %v", incrementBackUpPathList)
	incrementBackUpPathCount := len(incrementBackUpPathList)

	lastBackUpPath := ""
	//根据时间对备份文件进行排序,找到最新的备份文件
	for index, path := range incrementBackUpPathList {
		// xtrabackup --defaults-file=/data/zmysql/192.168.12.200_3306/testdb/testdb1/IncrBackup/20181115231000/backup-my.cnf
		// --prepare --apply-log-only --target-dir=/data/zmysql/192.168.12.200_3306/testdb/testdb1/FullBackup/20181115231000
		// --incremental-dir=/data/zmysql/192.168.12.200_3306/testdb/testdb1/IncrBackup/20181115231000
		//lastFilePath := utils.GetFileName(path)

		path = filepath.Join(incrementBackUpDir, path)
		confFile := filepath.Join(path, "backup-my.cnf")
		cmdStr := ""
		if index+1 == incrementBackUpPathCount {
			cmdStr = fmt.Sprintf("%s --defaults-file=%s --prepare --target-dir=%s --incremental-dir=%s",
				xtrabackupCmd, confFile, completeBackUpPath, path)
			lastBackUpPath = path
		} else {
			cmdStr = fmt.Sprintf("%s --defaults-file=%s --prepare --apply-log-only --target-dir=%s --incremental-dir=%s",
				xtrabackupCmd, confFile, completeBackUpPath, path)
		}

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
	}

	if lastBackUpPath != "" {
		lastBackUpPath = filepath.Join(lastBackUpPath, "xtrabackup_binlog_info")
	}

	er := executor.SuccessulExecuteResultNoData("recover backup data successful ")
	er.ResultData = map[string]string{
		"lastLogBinPath": lastBackUpPath,
	}
	return er
}

//CopyFullBackUpDataToMySQLInstance 复制全量备份文件到实例数据目录
func CopyFullBackUpDataToMySQLInstance(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//var xtrabackupCmd string
	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")

	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(true, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	completeBackUpPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	completeBackUpDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	completeBackUpPath = filepath.Join(completeBackUpDir, utils.GetFileName(completeBackUpPath))

	confFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlParameterFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// xtrabackup --defaults-file=/zmysql/db/yyy1/yyy102/conf/my.cnf --move-back
	// --target-dir=/data/zmysql/192.168.12.200_3306/testdb/testdb1/FullBackup/20181115231000
	cmdStr := fmt.Sprintf("%s --defaults-file=%s --move-back --target-dir=%s", xtrabackupCmd, confFile, completeBackUpPath)
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
	return executor.SuccessulExecuteResult(es, true, "copy full backup to mysql instance successful ")
}

//RecoverOneIncrementBackUpData 恢复一个单独的增量备份数据
func RecoverOneIncrementBackUpData(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//var xtrabackupCmd string
	//xtrabackupPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}
	//xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")

	xtrabackupCmd, err := getBackupOrRecoverExecuteCommand(true, params)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	completeBackUpPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	completeBackUpDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamCompleteBackUpDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	completeBackUpPath = filepath.Join(completeBackUpDir, utils.GetFileName(completeBackUpPath))

	incrementBackUpPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamIncrementBackUpPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	incrementBackUpDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamIncrementBackUpDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	lastFilePath := utils.GetFileName(incrementBackUpPath)

	isLast, _ := executor.ExtractCmdFuncBoolParam(params, "isLast")

	lastBackUpPath := ""

	// xtrabackup --defaults-file=/data/zmysql/192.168.12.200_3306/testdb/testdb1/IncrBackup/20181115231000/backup-my.cnf
	// --prepare --apply-log-only --target-dir=/data/zmysql/192.168.12.200_3306/testdb/testdb1/FullBackup/20181115231000
	// --incremental-dir=/data/zmysql/192.168.12.200_3306/testdb/testdb1/IncrBackup/20181115231000

	path := filepath.Join(incrementBackUpDir, lastFilePath)
	confFile := filepath.Join(path, "backup-my.cnf")
	cmdStr := ""
	if isLast {
		cmdStr = fmt.Sprintf("%s --defaults-file=%s --prepare --target-dir=%s --incremental-dir=%s",
			xtrabackupCmd, confFile, completeBackUpPath, path)
		lastBackUpPath = path
	} else {
		cmdStr = fmt.Sprintf("%s --defaults-file=%s --prepare --apply-log-only --target-dir=%s --incremental-dir=%s",
			xtrabackupCmd, confFile, completeBackUpPath, path)
	}

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

	if lastBackUpPath != "" {
		lastBackUpPath = filepath.Join(lastBackUpPath, "xtrabackup_binlog_info")
	}

	er := executor.SuccessulExecuteResultNoData("recover backup data successful ")
	er.ResultData = map[string]string{
		"lastLogBinPath": lastBackUpPath,
	}
	return er
}

// RenameSchema 执行rename schema
func RenameSchema(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	schemaList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamSchemaList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	renameSchemaList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamRenameSchemaList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if len(schemaList) != len(renameSchemaList) {
		log.Error("len(schemaList) != len(renameSchemaList)")
		return executor.ErrorExecuteResult(fmt.Errorf("逻辑恢复参数不合法"))
	}

	tableList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamTableList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	deleteOldSchema, _ := executor.ExtractCmdFuncBoolParam(params, CmdParamDeleteOldSchema)

	for index, schema := range schemaList {
		renameSchema := renameSchemaList[index]
		err = renameOneSchema(schema, renameSchema, host, mysqlUser, mysqlPassword, tableList, port, deleteOldSchema)
		if err != nil {
			log.Error("renameOneSchema failed ", err)
			return executor.ErrorExecuteResult(err)
		}
	}

	return executor.SuccessulExecuteResultNoData("rename schema execute successful")
}

func stringSliceEmpty(inputParams []string) (invalid bool) {
	for _, text := range inputParams {
		if strings.TrimSpace(text) == "" {
			invalid = true
			break
		}
	}
	return
}

func backupTableExist(backupTable string, currentTableList []string) (exist bool) {
	for _, table := range currentTableList {
		if backupTable == table {
			exist = true
			return
		}
	}
	return
}

func renameOneSchema(schema, backupSchemaName, host, username, password string, tables []string, port int, deleteOldSchema bool) (err error) {
	var renameTableList []string
	nowTables, err := getOneSchemaTables(host, username, password, schema, port)
	if err != nil {
		return
	}

	if stringSliceEmpty(tables) {
		renameTableList = nowTables
	} else {
		//check backup tables exist or not in current time
		for _, backupTable := range tables {
			if backupTableExist(backupTable, nowTables) {
				renameTableList = append(renameTableList, backupTable)
			}
		}
	}

	createDBSQL := fmt.Sprintf("create database if not exists %s", backupSchemaName)
	sqlList := []string{createDBSQL}

	for _, table := range renameTableList {
		renameTableSQL := fmt.Sprintf("rename table %s.%s to %s.%s", schema, table, backupSchemaName, table)
		sqlList = append(sqlList, renameTableSQL)
	}

	if deleteOldSchema {
		deleteSQL := fmt.Sprintf("drop database %s", schema)
		sqlList = append(sqlList, deleteSQL)
	}
	err = ExecuteMySQLCommand(host, username, password, strings.Join(sqlList, ";"), port)
	return
}

func getOneSchemaTables(host, username, password, schema string, port int) (tables []string, err error) {
	sql := fmt.Sprintf("select TABLE_NAME from information_schema.tables where TABLE_SCHEMA='%s';", schema)
	err = ExecuteMySQLCommandQuery(host, username, password, sql, port, true, &tables)
	return
}

// LogicalRecover 执行逻辑恢复
func LogicalRecover(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	myloaderPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMyLoaderPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	logicalBackupFilePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamLogicalBackupFilePath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//myloader --host 192.168.11.175 --port 10000 --user root --password 123qwe \
	//-e -d /home/backup/ms01/ms0101/2002111439
	command := fmt.Sprintf("%s --host %s --port %d --user %s --password %s -e -o -d %s", myloaderPath, host, port, mysqlUser, mysqlPassword, logicalBackupFilePath)
	log.Info("logical recover command = ", command)
	es, err := e.ExecShell(command)
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

func LogicalBackup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	mydumperPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMyDumperPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	logicalBackupFilePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamLogicalBackupFilePath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	workInfo, err := executor.ExtractCmdFuncStringParam(params, CmdParamMyWorkInfo)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	backupData, err := executor.ExtractCmdFuncBoolParam(params, CmdParamBackupData)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//myloader --host 192.168.11.175 --port 10000 --user root --password 123qwe \
	//-e -d /home/backup/ms01/ms0101/2002111439
	command := fmt.Sprintf("%s --host %s --port %d --user %s --password %s -e %s -o %s", mydumperPath, host, port, mysqlUser, mysqlPassword, workInfo, logicalBackupFilePath)
	if !backupData {
		command = fmt.Sprintf("%s --host %s --port %d --user %s --password %s -e  -d %s -o %s", mydumperPath, host, port, mysqlUser, mysqlPassword, workInfo, logicalBackupFilePath)
	}
	log.Info("logical recover command = ", command)
	es, err := e.ExecShell(command)
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
	return executor.SuccessulExecuteResult(es, true, "backup data successful ")
}

// BinlogRecover 执行binlog恢复
func BinlogRecover(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	binlogPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamBinlogPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	binlogBackupFilePath, err := executor.ExtractCmdFuncStringParam(params, CmdParamBinlogBackupFilePath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	//if err != nil {
	//	return executor.ErrorExecuteResult(err)
	//}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	stopDateTime, _ := executor.ExtractCmdFuncStringParam(params, CmdParamStopDateTime)

	//检查binlog文件是否存在
	checkCommand := fmt.Sprintf("ls %s", binlogBackupFilePath)
	log.Info("checkCommand = %s", checkCommand)
	es, err := e.ExecShell(checkCommand)
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

	var command string
	if stopDateTime != "" {
		//mysqlbinlog --stop-datetime='2020-02-24 10:11:40' mysql-bin.000025|-h127.0.0.1 -P3306 -uroot -p123456
		command = fmt.Sprintf("%s --stop-datetime='%s' %s|mysql -h%s -P%d -u%s -p%s", binlogPath, stopDateTime, binlogBackupFilePath, "127.0.0.1", port, mysqlUser, mysqlPassword)
	} else {
		//mysqlbinlog mysql-bin.000019|mysql -h127.0.0.1 -P3306 -uroot -p123456
		command = fmt.Sprintf("%s %s|mysql -h%s -P%d -u%s -p%s", binlogPath, binlogBackupFilePath, "127.0.0.1", port, mysqlUser, mysqlPassword)
	}

	log.Info("binlog recover command = %s", command)
	es, err = e.ExecShell(command)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Info("binlog recover result =%v   ", es)
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.SuccessulExecuteResult(es, true, "binlog recover backup data successful ")
}

//InterruptDatabaseRemoteService interrupt the database remote service
func InterruptDatabaseRemoteService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	confFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlParameterFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//sed -i -r '/^bind-address|^[[:space:]]+bind_address/a bind-address=127.0.0.1' my.cnf
	command := fmt.Sprintf("sed -i -r '/^bind-address|^[[:space:]]+bind_address/a bind-address=127.0.0.1' %s", confFile)
	//sed -i -r '/^\[mysqld\]|^[[:space:]]+\[mysqld\]/a bind-address=127.0.0.1' my.cnf
	command += fmt.Sprintf(` ; sed -i -r '/^\[mysqld\]|^[[:space:]]+\[mysqld\]/a bind-address=127.0.0.1' %s`, confFile)
	log.Info("interrupt the database remote service command = %s", command)
	es, err := e.ExecShell(command)
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
	return executor.SuccessulExecuteResult(es, true, "interrupt the database remote service successful ")
}

//RecoverDatabaseRemoteService recover the database remote service
func RecoverDatabaseRemoteService(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	confFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlParameterFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//sed -i -r '/bind-address=127.0.0.1/d' my.cnf
	command := fmt.Sprintf(`sed -i -r '/bind-address=127.0.0.1/d' %s`, confFile)
	log.Info("recover the database remote service command = %s", command)
	es, err := e.ExecShell(command)
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
	return executor.SuccessulExecuteResult(es, true, "recover the database remote service successful ")
}
