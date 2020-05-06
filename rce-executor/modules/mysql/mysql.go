package mysql

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules/osservice"
	"hercules_compiler/rce-executor/utils"
	"hercules_compiler/rce-executor/modules"
	"regexp"
	"strings"
	"time"
)

// 模块名常量定义
const (
	MysqlModuleName = "mysql"
)

// 函数名常量定义
const (
	CmdNameInitMysqlInstance                 = "InitializeMySQLInstance"
	CmdNameStartupMysqlInstance              = "StartupMySQLInstance"
	CmdNameMysqlInstanceReadiness            = "MySQLInstanceReadiness"
	CmdNameMysqlInstanceAlive                = "MySQLInstanceAlive"
	CmdNameMysqlInstanceAliveByMyClient      = "MySQLInstanceAliveByMyClient"
	CmdNameInstanceAliveWithCode             = "InstanceAliveWithCode"
	CmdNameMysqlCmdSQL                       = "MySQLCmdSQL"
	CmdNameShutDownMySQLInstance             = "ShutDownMySQLInstance"
	CmdNameStartMySQLService                 = "StartMySQLService"
	CmdNameCheckMySQLInstanceAlive           = "CheckMySQLInstanceAlive"
	CmdNameCheckMySQLUpgrade                 = "MySQLUpgrade"
	CmdNameCopyFullBackUpDataToMySQLInstance = "CopyFullBackUpDataToMySQLInstance"
	CmdNameSetSalveDelay                     = "SetSalveDelay"
	CmdNameAddMemberToProxySQL               = "AddMemberToProxySQL"
	CmdNameMGRAddMemberToProxySQL            = "MGRAddMemberToProxySQL"
	CmdNameStopMysqlInstance                 = "StopMysqlInstance"
	CmdNameGetVersion                        = "GetVersion"
	CmdNameMySQLClusterHasMaster             = "ClusterHasMaster"
	CmdNameSelectMySQLClusterMasterInfo      = "SelectMySQLClusterMasterInfo"
	CmdNameCheckMySQLCommunityClientExist    = "CheckMySQLCommunityClientExist"
)

// 命令参数常量定义
const (
	CmdParamMysqlPath             = "mysqlPath"
	CmdParamMysqlRunWay           = "mysqlRunWay"
	CmdParamRootPassword          = "rootPassword"
	CmdParamPort                  = "port"
	CmdParamPortList              = "ports"
	CmdParamDataBaseDir           = "dataBaseDir"
	CmdParamMysqlRunningUser      = "user"
	CmdParamMysqlServerID         = "serverId"
	CmdParamMysqlDataDir          = "mysqlDataDir"
	CmdParamSocketFile            = "socketFile"
	CmdParamMysqlParameterFile    = "mysqlParameterFile"
	CmdParamMysqlUser             = "mysqlUser"
	CmdParamMysqlPassword         = "mysqlPassword"
	CmdParamClusterName           = "clusterName"
	CmdParamHost                  = "host"
	CmdParamHostList              = "hosts"
	CmdParamHostName              = "hostname"
	CmdParamMysqlCmdSQL           = "cmdSql"
	CmdParamTimeout               = "timeout"
	CmdParamLogFileSize           = "logSize"
	CmdParamMySQLPort             = "mysqlPort"
	CmdParamDoProxySQL            = "doProxySQL"
	CmdParamDoSlaveDelay          = "doSlaveDelay"
	CmdParamDelayTime             = "delayTime"
	CmdParamMiddlewareType        = "middlewareType"
	CmdParamCompleteBackUpPath    = "completeBackUpPath"
	CmdParamCompleteBackUpDir     = "completeBackUpDir"
	CmdParamIncrementBackUpPath   = "incrementBackUpPath"
	CmdParamIncrementBackUpDir    = "incrementBackUpDir"
	CmdParamSchemaList            = "schemaList"
	CmdParamRenameSchemaList      = "renameSchemaList"
	CmdParamDeleteOldSchema       = "deleteOldSchema"
	CmdParamTableList             = "tableList"
	CmdParamMyLoaderPath          = "myloaderPath"
	CmdParamMyDumperPath          = "mydumperPath"
	CmdParamMyWorkInfo            = "workInfo"
	CmdParamBackupData            = "backupData"
	CmdParamLogicalBackupFilePath = "logicalBackupFilePath"
	CmdParamBinlogPath            = "binlogPath"
	CmdParamBinlogBackupFilePath  = "binlogBackupFilePath"
	CmdParamStopDateTime          = "stopDateTime"
)

//定义默认值
const (
	MinDefaultTimeout = 60 * 2
	MaxDefaultTimeout = 10 * MinDefaultTimeout
	MinDefaultLogSize = 1
	MaxDefaultLogSize = 8 * MinDefaultLogSize
	MasterText        = "1" //主节点
	SlaveText         = "2" //从节点
)

//InstanceInfo mysql instance info
type InstanceInfo struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewInstanceInfo(host, username, password string, port int) *InstanceInfo {
	return &InstanceInfo{Host: host, Port: port, Username: username, Password: password}
}

// InitializeMySQLInstance 初始化MySQL实例
func InitializeMySQLInstance(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	/*rootPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamRootPassword)
	if err != nil {
		rootPassword = ""
	}*/
	runningUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlRunningUser)
	if err != nil {
		runningUser = "mysql"
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)

	if err != nil {
		port = 3306
	}

	serverID, err := executor.ExtractCmdFuncIntParam(params, CmdParamMysqlServerID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	dataBaseDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamDataBaseDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlDataDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlDataDir)
	if err != nil {
		mysqlDataDir = dataBaseDir + "/data"
	}
	ibdataFilename := mysqlDataDir + "/ibdata1"

	logFileSize, err := executor.ExtractCmdFuncIntParam(params, CmdParamLogFileSize)
	if err != nil {
		logFileSize = MinDefaultLogSize
	}

	if logFileSize > MaxDefaultLogSize {
		logFileSize = MaxDefaultLogSize
	}

	es, err := e.ExecShell("file " + ibdataFilename)

	if err != nil {
		return executor.ErrorExecuteResult(fmt.Errorf("can not detect mysql file '%s',error=%s", ibdataFilename, err.Error()))
	}
	msg := strings.Join(es.Stdout, "\n") + strings.Join(es.Stderr, "\n")

	//如果没有这个确认文件不存在的信息，就报错
	if (strings.Index(msg, "cannot open") == -1) || (strings.Index(msg, "No such file or directory") == -1) {
		return executor.ErrorExecuteResult(fmt.Errorf("mysql file '%s' exists or Permission denied", ibdataFilename))
	}
	//创建必要的目录
	execParams := executor.ExecutorCmdParams{}
	execParams["path"] = mysqlDataDir

	er := oscmd.MakeDir(e, &execParams)

	if !er.Successful {
		er.Message = "mkdir mysql data directory error:" + er.Message
		return er
	}
	execParams["path"] = dataBaseDir + "/logs"
	er = oscmd.MakeDir(e, &execParams)
	if !er.Successful {
		er.Message = "mkdir mysql logs directory error:" + er.Message
		return er
	}
	execParams["path"] = dataBaseDir + "/tmp"
	er = oscmd.MakeDir(e, &execParams)
	if !er.Successful {
		er.Message = "mkdir mysql tmp directory error:" + er.Message
		return er
	}
	execParams["path"] = dataBaseDir + "/run"
	er = oscmd.MakeDir(e, &execParams)
	if !er.Successful {
		er.Message = "mkdir mysql run directory error:" + er.Message
		return er
	}

	execParams = executor.ExecutorCmdParams{}
	execParams["filenamePattern"] = dataBaseDir + "/{data,logs,tmp,run}"
	execParams["own"] = "mysql"
	execParams["group"] = "mysql"

	er = oscmd.ChangeOwnAndGroup(e, &execParams)
	if !er.Successful {
		return er
	}

	socketFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamSocketFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlParameterFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlParameterFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//	cmdStr := fmt.Sprintf("%s/bin/mysqld --innodb_undo_tablespaces=3 --basedir=%s --initialize-insecure --innodb_log_file_size=%dG --user=%s  --server-id=%d --datadir=%s --port=%d --socket=%s --log-bin=%s --log-error=%s",
	//		mysqlPath, mysqlPath, logFileSize, runningUser, serverID, mysqlDataDir,
	//		port, socketFile, dataBaseDir+"/logs/mysql-bin", dataBaseDir+"/logs/mysql-error-log.err")

	cmdStr := fmt.Sprintf("%s/bin/mysqld --defaults-file=%s --basedir=%s --initialize-insecure --user=%s  --server-id=%d --datadir=%s --port=%d --socket=%s --log-bin=%s --log-error=%s",
		mysqlPath, mysqlParameterFile, mysqlPath, runningUser, serverID, mysqlDataDir,
		port, socketFile, dataBaseDir+"/logs/mysql-bin", dataBaseDir+"/logs/mysql-error-log.err")
	log.Debug("--111-------------------%s", cmdStr)
	es, err = e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(es.Stderr) == 0 && len(es.Stdout) == 0 {
		return executor.SuccessulExecuteResult(es, true, "mysql instance "+dataBaseDir+" initialize successful")
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
	return executor.SuccessulExecuteResult(es, true, "mysql instance "+dataBaseDir+" initialize successful")
}

// StartupMySQLInstance 启动MySQL实例
func StartupMySQLInstance(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	mysqlRunWay, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlRunWay)
	if err != nil {
		mysqlRunWay = "mysqld"
	}
	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	runningUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlRunningUser)
	if err != nil {
		runningUser = "mysql"
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)

	if err != nil {
		port = 3306
	}

	mysqlDataDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlDataDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlParameterFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlParameterFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	socketFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamSocketFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//先判断数据库是否已经启动
	cmdParams := executor.ExecutorCmdParams{}
	for k, v := range *params {
		cmdParams[k] = v
	}
	//cmdParams[CmdParamSocketFile] = socketFile
	er := MySQLInstanceAlive(e, &cmdParams)
	log.Debug("startup mysql check alive:%s\n", er.Message)

	if er.Successful {
		er.Message = "MySQL Instance not need startup , is running"
		er.Changed = false
		return er
	}

	cmdStr := fmt.Sprintf("%s/bin/%s --defaults-file=%s  --basedir=%s --user=%s --datadir=%s --port=%d --socket=%s",
		mysqlPath, mysqlRunWay, mysqlParameterFile, mysqlPath, runningUser, mysqlDataDir,
		port, socketFile)
	log.Debug("startup mysql instance cmd:%s\n", cmdStr)
	er.Successful = true
	er.Message = "MySQL Instance need to startup"
	er.ResultData = make(map[string]string)
	er.ResultData["cmdStr"] = cmdStr
	return er
}

//ShutDownMySQLInstance 关闭mysql实例
func ShutDownMySQLInstance(e executor.Executor, params *executor.ExecutorCmdParams) (er executor.ExecuteResult) {
	if (*params)["dbName"] == nil && (*params)["instanceName"] == nil {
		log.Debug("Can Not Shut Down Mysql Instance")
		er.Message = "Can Not Shut Down Mysql Instance"
		er.Changed = false
		return
	}
	serviceName := GetServiceName(params)
	version := e.GetExecutorContext(executor.ContextNameVersion)
	iVersion, _ := osservice.GetLinuxDistVer(version)
	if iVersion < 7 {
		//start服务
		startServiceStr := "service " + serviceName + " stop"

		e, err := e.ExecShell(startServiceStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if e.ExitCode != 0 {
			return executor.ErrorExecuteResult(fmt.Errorf(strings.Join(e.Stdout, "")))
		}

	} else {
		//停止服务
		//serviceName = serviceName + ".service"
		startServiceStr := "systemctl stop " + strings.ToLower(serviceName) + ".service"
		e, err := e.ExecShell(startServiceStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		if e.ExitCode != 0 {
			return executor.ErrorExecuteResult(fmt.Errorf(strings.Join(e.Stdout, "")))
		}
	}

	er.Message = "ShutDown Mysql Instance successful"
	er.Successful = true
	er.Changed = true
	return
}

//StopMysqlInstance 关闭mysql实例
func StopMysqlInstance(e executor.Executor, params *executor.ExecutorCmdParams) (er executor.ExecuteResult) {
	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	mysqlHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	mysqlPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMySQLPort)
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

	cmdStr := fmt.Sprintf("%s/bin/mysqladmin -h%s -P%d -u%s -p%s shutdown 2>/dev/null", mysqlPath, mysqlHost, mysqlPort, mysqlUser, mysqlPassword)
	//cmdStr = "ps -a|grep mysql"

	log.Debug("stop mysql instance cmd = %s", cmdStr)

	statusCode, err := e.ExecShell(cmdStr)
	if statusCode.ExitCode != 0 {
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		return executor.ErrorExecuteResult(fmt.Errorf("Stop Mysql Instance Failed"))
	}

	er.Message = "Stop Mysql Instance successful"
	er.Successful = true
	er.Changed = true
	return
}

// StartMySQLService 启动mysql服务
func StartMySQLService(e executor.Executor, params *executor.ExecutorCmdParams) (er executor.ExecuteResult) {
	if (*params)["dbName"] == nil && (*params)["instanceName"] == nil {
		log.Debug("Can Not Start Mysql Service")
		er.Message = "Can Not Start Mysql Service"
		er.Changed = false
		return
	}
	serviceName := GetServiceName(params)
	version := e.GetExecutorContext(executor.ContextNameVersion)
	iVersion, _ := osservice.GetLinuxDistVer(version)
	if iVersion < 7 {
		//start服务
		startServiceStr := "service " + serviceName + " start"

		_, err := e.ExecShell(startServiceStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

	} else {
		//停止服务
		//serviceName = serviceName + ".service"
		startServiceStr := "systemctl start " + strings.ToLower(serviceName) + ".service"
		_, err := e.ExecShell(startServiceStr)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
	}
	er.Message = "Start Mysql Instance successful"
	er.Successful = true
	er.Changed = true
	return
}

// GetServiceName 获取服务名
func GetServiceName(params *executor.ExecutorCmdParams) (serviceName string) {
	if (*params)["dbName"] != nil && (*params)["instanceName"] != nil {
		serviceName = "zmysql_" + (*params)["dbName"].(string) + "_" + (*params)["instanceName"].(string)
	} else {
		if (*params)["dbName"] != nil {
			serviceName = "zmysql_" + (*params)["dbName"].(string)
		} else if (*params)["instanceName"] != nil {
			serviceName = "zmysql_" + (*params)["instanceName"].(string)
		} else {
			serviceName = "mysql_service"
		}

	}
	return
}

//创建service自启动
func createService(e executor.Executor, params *executor.ExecutorCmdParams, cmdStr string, baseDir string) executor.ExecuteResult {
	version := e.GetExecutorContext(executor.ContextNameVersion)
	iVersion, _ := osservice.GetLinuxDistVer(version)
	(*params)["serviceName"] = GetServiceName(params)
	(*params)["workingDirectory"] = baseDir
	(*params)["serviceCmdLine"] = cmdStr
	if iVersion < 7 {
		//默认runorder为30
		(*params)["runOrder"] = 30
		return osservice.GenerateSysVService(e, params)
	} else {
		return osservice.GenerateSystemdService(e, params)
	}
}

// MySQLInstanceReadiness 通过连接测试，判断MySQL实例是否准备就绪
func MySQLInstanceReadiness(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port := int(0)
	host := ""

	socketFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamSocketFile)

	if err != nil {
		if host, err = executor.ExtractCmdFuncStringParam(params, CmdParamHost); err != nil || host == "" {
			host = "localhost"
		}
		if port, err = executor.ExtractCmdFuncIntParam(params, CmdParamPort); err != nil || port == 0 {
			return executor.ErrorExecuteResult(err)
		}
		socketFile = ""
	}
	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		mysqlUser = "root"
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		mysqlPassword = ""
	}

	passwordStr := ""
	if mysqlPassword != "" {
		passwordStr = "-p" + strings.TrimSpace(mysqlPassword)
	}

	hostStr := ""
	if socketFile != "" {
		hostStr = "-S " + socketFile
	} else {
		hostStr = fmt.Sprintf("-h %s -P %d ", host, port)
	}

	userStr := ""
	if mysqlUser != "" {
		userStr = "-u" + mysqlUser
	}

	selectStr := `" select 1 " `
	cmdStr := fmt.Sprintf("%s/bin/mysql --no-defaults %s %s %s  -E -e %s ", mysqlPath, hostStr, userStr, passwordStr, selectStr)
	//log.Printf("mysql readiness detect,cmd=%s\n", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//log.Printf("mysql readiness: exit code=%d msg=%s\n", es.ExitCode, strings.Join(es.Stderr, "\n"))
	if es.ExitCode != 0 {
		var errMsg string

		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
			//尝试以port连接
			//ERROR 2002 (HY000): Can't connect to local MySQL server through socket '/tmp/mysql3306.sock' (2)
			if socketFile != "" && host != "" && port > 0 && strings.Index(errMsg, "ERROR 2002") > -1 {
				hostStr = fmt.Sprintf("-h %s -P %d ", host, port)
				cmdStr := fmt.Sprintf("%s/bin/mysql --no-defaults %s %s %s  -E  -e %s ", mysqlPath, hostStr, userStr, passwordStr, selectStr)
				//log.Printf("check mysql instance: cmd='%s'\n", cmdStr)
				es2, err := e.ExecShell(cmdStr)
				if err != nil {
					return executor.ErrorExecuteResult(err)
				}
				if es2.ExitCode != 0 {
					var errMsg2 string
					if len(es.Stderr) == 0 {
						errMsg2 = executor.ErrMsgUnknow
					} else {
						errMsg2 = strings.Join(es.Stderr, "\n")
					}
					return executor.NotSuccessulExecuteResult(es2, errMsg2)
				}
				return executor.SuccessulExecuteResult(es2, false, fmt.Sprintf("mysql readiness successful"))
			}
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}

	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("mysql readiness successful"))
}

//MySQLInstanceAliveByMyclient 使用golang的客户端实现判断MySQL是否可用
func MySQLInstanceAliveByMyClient(_ executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		host string
		port int
		err  error
	)
	es := new(executor.ExecutedStatus)
	if host, err = executor.ExtractCmdFuncStringParam(params, CmdParamHost); err != nil || host == "" {
		host = "localhost"
	}
	if port, err = executor.ExtractCmdFuncIntParam(params, CmdParamPort); err != nil || port == 0 {
		return executor.ErrorExecuteResult(err)
	}
	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		mysqlUser = "root"
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		mysqlPassword = ""
	}
	log.Debug("mysql readiness detect,cmd=%s........\n", host, port, mysqlPassword, mysqlUser)
	err = ExecuteMySQLPing(host, mysqlUser, mysqlPassword, port)
	if err != nil {
		if strings.Index(err.Error(), "Error 1045") > -1 {
			return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("mysql readiness successful"))
		}
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("mysql readiness successful"))
}

// MySQLInstanceAlive 判断MySQL实例是否可用
func MySQLInstanceAlive(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var er executor.ExecuteResult
	socketFile, _ := executor.ExtractCmdFuncStringParam(params, CmdParamSocketFile)
	mysqlPassword, _ := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if socketFile == "" && mysqlPassword != "" {
		er = MySQLInstanceAliveByMyClient(e, params)
	} else {
		er = MySQLInstanceReadiness(e, params)
	}

	if er.Successful {
		er.Message = "mysql instance is alive"
		return er
	}
	//ERROR 1045 (28000): Access denied for user 'root'@'localhost' (using password: NO)
	//如果是密码错误，是可以认为是好的
	if strings.Index(er.Message, "ERROR 1045") > -1 {
		er.Message = "mysql instance is alive"
		er.ExitCode = 0
		er.Changed = false
		er.Successful = true
		return er
	}
	return er
}

// InstanceAliveWithCode 判断MySQL实例是否可用，返回成功或者失败的数据
func InstanceAliveWithCode(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	er := MySQLInstanceReadiness(e, params)
	er.ResultData = make(map[string]string)
	if er.Successful {
		er.ResultData["alive"] = "true"
		er.Message = "mysql instance is alive"
		return er
	}
	//ERROR 1045 (28000): Access denied for user 'root'@'localhost' (using password: NO)
	//如果是密码错误，是可以认为是好的
	if strings.Index(er.Message, "ERROR 1045") > -1 {
		er.Message = "mysql instance is alive"
		er.ResultData["alive"] = "true"
		er.ExitCode = 0
		er.Changed = false
		er.Successful = true
		return er
	}
	er.Successful = true
	er.ExitCode = 0
	er.Changed = false
	er.Message = "mysql instance is not alive"
	er.ResultData["alive"] = "false"
	return er
}

// CheckMySQLInstanceAlive 判断MySQL实例是否可用
func CheckMySQLInstanceAlive(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var er executor.ExecuteResult
	timeout, err := executor.ExtractCmdFuncIntParam(params, CmdParamTimeout)
	if err != nil {
		timeout = MinDefaultTimeout
	}

	if timeout > MaxDefaultTimeout {
		timeout = MaxDefaultTimeout
	}

	start := time.Now().Unix()
	for {
		current := time.Now().Unix()
		pass := current - start
		if int(pass) > timeout {
			log.Debug("CheckMySQLInstanceAlive timeout  pass = %d", pass)
			break
		} else {
			er = MySQLInstanceAlive(e, params)
			if er.Successful {
				return er
			}
		}
		time.Sleep(time.Second * 5)
	}
	return er
}

// MySQLCmdSQL 执行sql语句
func MySQLCmdSQL(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdSQL, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlCmdSQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port := int(0)
	host := ""

	socketFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamSocketFile)

	if err != nil {
		if host, err = executor.ExtractCmdFuncStringParam(params, CmdParamHost); err != nil || host == "" {
			host = "localhost"
		}
		if port, err = executor.ExtractCmdFuncIntParam(params, CmdParamPort); err != nil || port == 0 {
			return executor.ErrorExecuteResult(err)
		}
		socketFile = ""
	}

	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		mysqlUser = "root"
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		mysqlPassword = ""
	}

	passwordStr := ""
	if mysqlPassword != "" {
		passwordStr = "-p" + strings.TrimSpace(mysqlPassword)
	}

	hostStr := ""
	if socketFile != "" {
		hostStr = "-S " + socketFile
	} else {
		hostStr = fmt.Sprintf("-h %s -P %d ", host, port)
	}

	userStr := ""
	if mysqlUser != "" {
		userStr = "-u" + mysqlUser
	}

	cmdStr := fmt.Sprintf("%s/bin/mysql --no-defaults %s %s %s  -E -e \"%s\" ", mysqlPath, hostStr, userStr, passwordStr, utils.EscapeShellCmd(cmdSQL))
	es, err := e.ExecShell(cmdStr)
	log.Debug("cmdStr string= %s", cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("result=%v", es.Stdout)

	err = executor.GetExecResult(es)
	if err != nil {
		log.Debug("executor sql failed=%v", err)
		return executor.ErrorExecuteResult(err)
	}

	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("mysql cmd sql statmtent execute successful"))
}

func MySQLUpgrade(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlUser)
	if err != nil {
		mysqlUser = "root"
	}

	mysqlPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncStringParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdStr := fmt.Sprintf("%s/bin/mysql_upgrade --protocol=tcp -P%s -u%s -p%s  --force", mysqlPath, port, mysqlUser, mysqlPassword)
	es, err := e.ExecShell(cmdStr)
	log.Debug("cmdStr string= %s", cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//log.Debug("result=%v", es.Stdout)

	err = executor.GetExecResult(es)
	if err != nil {
		log.Debug("executor sql failed=%v", err)
		return executor.ErrorExecuteResult(err)
	}

	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("mysql_upgrade execute successful"))
}

// InitializeMySQLInstanceForRecover 初始化异机器恢复 MySQL实例
func InitializeMySQLInstanceForRecover(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	dataBaseDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamDataBaseDir)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlDataDir, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlDataDir)
	if err != nil {
		mysqlDataDir = dataBaseDir + "/data"
	}
	ibdataFilename := mysqlDataDir + "/ibdata1"

	logFileSize, err := executor.ExtractCmdFuncIntParam(params, CmdParamLogFileSize)
	if err != nil {
		logFileSize = MinDefaultLogSize
	}

	if logFileSize > MaxDefaultLogSize {
		logFileSize = MaxDefaultLogSize
	}

	es, err := e.ExecShell("file " + ibdataFilename)

	if err != nil {
		return executor.ErrorExecuteResult(fmt.Errorf("can not detect mysql file '%s',error=%s", ibdataFilename, err.Error()))
	}
	msg := strings.Join(es.Stdout, "\n") + strings.Join(es.Stderr, "\n")

	//如果没有这个确认文件不存在的信息，就报错
	if (strings.Index(msg, "cannot open") == -1) || (strings.Index(msg, "No such file or directory") == -1) {
		return executor.ErrorExecuteResult(fmt.Errorf("mysql file '%s' exists or Permission denied", ibdataFilename))
	}
	//创建必要的目录
	execParams := executor.ExecutorCmdParams{}
	execParams["path"] = mysqlDataDir

	er := oscmd.MakeDir(e, &execParams)

	if !er.Successful {
		er.Message = "mkdir mysql data directory error:" + er.Message
		return er
	}
	execParams["path"] = dataBaseDir + "/logs"
	er = oscmd.MakeDir(e, &execParams)
	if !er.Successful {
		er.Message = "mkdir mysql logs directory error:" + er.Message
		return er
	}
	execParams["path"] = dataBaseDir + "/tmp"
	er = oscmd.MakeDir(e, &execParams)
	if !er.Successful {
		er.Message = "mkdir mysql tmp directory error:" + er.Message
		return er
	}
	execParams["path"] = dataBaseDir + "/run"
	er = oscmd.MakeDir(e, &execParams)
	if !er.Successful {
		er.Message = "mkdir mysql run directory error:" + er.Message
		return er
	}

	execParams = executor.ExecutorCmdParams{}
	execParams["filenamePattern"] = dataBaseDir + "/{data,logs,tmp,run}"
	execParams["own"] = "mysql"
	execParams["group"] = "mysql"

	er = oscmd.ChangeOwnAndGroup(e, &execParams)
	if !er.Successful {
		return er
	}
	return executor.SuccessulExecuteResult(es, true, "mysql instance "+dataBaseDir+" initialize successful")
}

//GetVersion 获取mysql版本
func GetVersion(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	var cmdStr = fmt.Sprintf("%s --version", path)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		er := executor.SuccessulExecuteResult(es, false, "get mysql version successful")
		if len(es.Stdout) == 0 {
			return executor.ErrorExecuteResult(errors.New("have no stdout"))
		}
		version := ""
		versionStr := es.Stdout[0]
		versionField := ""
		re, err := regexp.Compile(`\s\d+\.\d+\.\d+.*?\s`)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		regStr := re.FindString(versionStr)
		if regStr != "" {
			versionField = regStr
			versionField = strings.Trim(regStr, " ")
		}
		if strings.Contains(versionStr, "MySQL Community Server") {
			version = versionField + "-MySQL Community Server"
		} else if strings.Contains(versionStr, "MariaDB Server") {
			version = versionField + " Server"
		} else if strings.Contains(versionStr, "Percona Server") {
			version = versionField + "-Percona Server"
		}
		er.ResultData = make(map[string]string)
		er.ResultData["version"] = version
		return er
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//CheckMySQLCommunityClientExist 检查mySQLClient是否存在
func CheckMySQLCommunityClientExist(e executor.Executor, _ *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdStr := "mysql --version"
	es, err := e.ExecShell(cmdStr)
	log.Debug("cmdStr string= %s", cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode != 0 || len(es.Stderr) != 0 {
		es := executor.SuccessulExecuteResult(es, false, fmt.Sprintf("MySQL Community Client No Exist"))
		es.ResultData = make(map[string]string)
		es.ResultData["exist"] = fmt.Sprintf("%d", 1) //1是表示不存在
		return es
	}
	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("MySQL Community Client Exist"))
}

func init() {
	modules.AddModule(MysqlModuleName)
	executor.RegisterCmd(MysqlModuleName, CmdNameInitMysqlInstance, InitializeMySQLInstance)
	executor.RegisterCmd(MysqlModuleName, CmdNameShutDownMySQLInstance, ShutDownMySQLInstance)
	executor.RegisterCmd(MysqlModuleName, CmdNameStartupMysqlInstance, StartupMySQLInstance)
	executor.RegisterCmd(MysqlModuleName, CmdNameMysqlInstanceReadiness, MySQLInstanceReadiness)
	executor.RegisterCmd(MysqlModuleName, CmdNameMysqlInstanceAlive, MySQLInstanceAlive)
	executor.RegisterCmd(MysqlModuleName, CmdNameMysqlInstanceAliveByMyClient, MySQLInstanceAliveByMyClient)
	executor.RegisterCmd(MysqlModuleName, CmdNameInstanceAliveWithCode, InstanceAliveWithCode)
	executor.RegisterCmd(MysqlModuleName, CmdNameMysqlCmdSQL, MySQLCmdSQL)
	executor.RegisterCmd(MysqlModuleName, CmdNameMySQLCmdSQLEx, MySQLCmdSQLEx)
	executor.RegisterCmd(MysqlModuleName, CmdNameCmdSQLWithResult, CmdSQLWithResult)
	executor.RegisterCmd(MysqlModuleName, CmdNameStartMySQLService, StartMySQLService)
	executor.RegisterCmd(MysqlModuleName, CmdNameSetProxySQLConfig, SetProxySQLConfig)
	executor.RegisterCmd(MysqlModuleName, CmdNameSetProxySQLConfigEx, SetProxySQLConfigEx)
	executor.RegisterCmd(MysqlModuleName, CmdNameProxyMySQLCmdSQL, ProxyMySQLCmdSQL)
	executor.RegisterCmd(MysqlModuleName, CmdNameProxyMySQLCmdSQLEx, ProxyMySQLCmdSQLEx)
	executor.RegisterCmd(MysqlModuleName, CmdNameInstallPlugin, InstallPlugin)
	executor.RegisterCmd(MysqlModuleName, CmdNameCheckMySQLInstanceAlive, CheckMySQLInstanceAlive)
	executor.RegisterCmd(MysqlModuleName, CmdNameCheckMySQLUpgrade, MySQLUpgrade)
	executor.RegisterCmd(MysqlModuleName, CmdNameRegistOrchestrator, RegistOrchestrator)
	executor.RegisterCmd(MysqlModuleName, CmdNameSendBackupDataBySocat, SendBackupDataBySocat)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecvBackupDataBySocat, RecvBackupDataBySocat)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecoverBackupData, RecoverBackupData)
	executor.RegisterCmd(MysqlModuleName, CmdNameSetProxySQLMasterOffLineSoft, SetProxySQLMasterOffLineSoft)
	executor.RegisterCmd(MysqlModuleName, CmdNameSetProxySQLMemberOnLine, SetProxySQLMemberOnLine)
	executor.RegisterCmd(MysqlModuleName, CmdNameOrchestratorSwitchMaster, OrchestratorSwitchMaster)
	executor.RegisterCmd(MysqlModuleName, CmdNameOrchestratorStartReplicaption, OrchestratorStartReplicaption)
	executor.RegisterCmd(MysqlModuleName, CmdNameKillMasterInTransThread, KillMasterInTransThread)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecoverIncrementBackUpData, RecoverIncrementBackUpData)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecoverCompleteBackUpData, RecoverCompleteBackUpData)
	executor.RegisterCmd(MysqlModuleName, CmdNameInitializeMySQLInstanceForRecover, InitializeMySQLInstanceForRecover)
	executor.RegisterCmd(MysqlModuleName, CmdNameCopyFullBackUpDataToMySQLInstance, CopyFullBackUpDataToMySQLInstance)
	executor.RegisterCmd(MysqlModuleName, CmdNameSetSalveDelay, SetSalveDelay)
	executor.RegisterCmd(MysqlModuleName, CmdNameMGRAddMemberToProxySQL, MGRAddMemberToProxySQL)
	executor.RegisterCmd(MysqlModuleName, CmdNameAddMemberToProxySQL, AddMemberToProxySQL)
	executor.RegisterCmd(MysqlModuleName, CmdNameAddMembersToProxySQL, AddMembersToProxySQL)
	executor.RegisterCmd(MysqlModuleName, CmdNameRemoveMemberFromProxySQL, RemoveMemberFromProxySQL)
	executor.RegisterCmd(MysqlModuleName, CmdNameGetKeepAlivedRole, GetKeepalivedRole)
	executor.RegisterCmd(MysqlModuleName, CmdNameStopMysqlInstance, StopMysqlInstance)
	executor.RegisterCmd(MysqlModuleName, CmdNameCheckBackupData, CheckBackupData)
	executor.RegisterCmd(MysqlModuleName, CmdNameGetVersion, GetVersion)
	executor.RegisterCmd(MysqlModuleName, CmdNameMySQLClusterHasMaster, ClusterHasMaster)
	executor.RegisterCmd(MysqlModuleName, CmdNameSelectMySQLClusterMasterInfo, SelectMySQLClusterMasterInfo)
	executor.RegisterCmd(MysqlModuleName, CmdNameSetProxySQLMGRMasterOffLineSoft, SetProxySQLMGRMasterOffLineSoft)
	executor.RegisterCmd(MysqlModuleName, CmdNameResetMGRMasterStatus, ResetMGRMasterStatus)
	executor.RegisterCmd(MysqlModuleName, CmdNameSwitchMGRMaster, SwitchMGRMaster)
	executor.RegisterCmd(MysqlModuleName, CmdNameCheckMySQLCommunityClientExist, CheckMySQLCommunityClientExist)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecoverBackUpFiles, RecoverBackUpFiles)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecoverOneIncrementBackUpData, RecoverOneIncrementBackUpData)
	executor.RegisterCmd(MysqlModuleName, CmdNameRenameSchema, RenameSchema)
	executor.RegisterCmd(MysqlModuleName, CmdNameLogicalRecover, LogicalRecover)
	executor.RegisterCmd(MysqlModuleName, CmdNameLogicalBackup, LogicalBackup)
	executor.RegisterCmd(MysqlModuleName, CmdNameBinlogRecover, BinlogRecover)
	executor.RegisterCmd(MysqlModuleName, CmdNameInterruptDatabaseRemoteService, InterruptDatabaseRemoteService)
	executor.RegisterCmd(MysqlModuleName, CmdNameRecoverDatabaseRemoteService, RecoverDatabaseRemoteService)
}
