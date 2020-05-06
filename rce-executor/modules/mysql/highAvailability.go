package mysql

import (
	"encoding/json"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/http"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/modules/osservice"
	"hercules_compiler/rce-executor/utils"
	"strconv"
	"strings"
	"time"
)

//定义高可用中间件类型
const (
	HighAvailabilityForProxySQL = 1
	HighAvailabilityForHAProxy  = 2
)

// 函数名常量定义
const (
	CmdNameSetProxySQLConfig                 = "SetProxySQLConfig"
	CmdNameSetProxySQLConfigEx               = "SetProxySQLConfigEx"
	CmdNameProxyMySQLCmdSQL                  = "ProxyMySQLCmdSQL"
	CmdNameProxyMySQLCmdSQLEx                = "ProxyMySQLCmdSQLEx"
	CmdNameRegistOrchestrator                = "RegistOrchestrator"
	CmdNameOrchestratorSwitchMaster          = "OrchestratorSwitchMaster"
	CmdNameSetProxySQLMasterOffLineSoft      = "SetProxySQLMasterOffLineSoft"
	CmdNameSetProxySQLMGRMasterOffLineSoft   = "SetProxySQLMGRMasterOffLineSoft"
	CmdNameMySQLCmdSQLEx                     = "MySQLCmdSQLEx"
	CmdNameCmdSQLWithResult                  = "SQLWithResult"
	CmdNameSetProxySQLMemberOnLine           = "SetProxySQLMemberOnLine"
	CmdNameOrchestratorStartReplicaption     = "OrchestratorStartReplicaption"
	CmdNameKillMasterInTransThread           = "KillMasterInTransThread"
	CmdNameRecoverCompleteBackUpData         = "RecoverCompleteBackUpData"
	CmdNameRecoverIncrementBackUpData        = "RecoverIncrementBackUpData"
	CmdNameInitializeMySQLInstanceForRecover = "InitializeMySQLInstanceForRecover"
	CmdNameGetKeepAlivedRole                 = "GetKeepalivedRole"
	CmdNameAddMembersToProxySQL              = "AddMembersToProxySQL"
	CmdNameRemoveMemberFromProxySQL          = "RemoveMemberFromProxySQL"
	CmdNameSwitchMGRMaster                   = "SwitchMGRMaster"
	CmdNameResetMGRMasterStatus              = "ResetMGRMasterStatus"
)

// 命令参数常量定义
const (
	CmdParamProxySQLPath           = "path"
	CmdParamProxySQLConfPath       = "confPath"
	CmdParamProxySQLBinPath        = "binPath"
	CmdParamMasterHost             = "masterHost"
	CmdParamMasterPort             = "masterPort"
	CmdParamSlaveHost              = "slaveHost"
	CmdParamSlavePort              = "slavePort"
	CmdParamMemberHost             = "memberHost"
	CmdParamMemberRole             = "memberRole"
	CmdParamMemberPort             = "memberPort"
	CmdParamMemberUser             = "memberUser"
	CmdParamMemberPassword         = "memberPassword"
	CmdParamHAProxyPort            = "haProxyPort"
	CmdParamHAProxyConfText        = "haProxyConfText"
	CmdParamHAProxyNewConfText     = "haProxyNewConfText"
	CmdParamSlaveReadAllowMaxDelay = "slaveReadAllowMaxDelay"
	CmdParamsMasterReadWeight      = "masterReadWeight"
	CmdParamSlaveReadWeight        = "slaveReadWeight"
)

var proxySQLConfigText = `
#!/bin/bash 
#
# chkconfig: 345 99 01
# description: High Performance and Advanced Proxy for MySQL and forks. 
# It provides advanced features like connection pool, query routing and rewrite, 
# firewalling, throttling, real time analysis, error-free failover
### BEGIN INIT INFO
# Provides:          %s
# Required-Start:    $local_fs
# Required-Stop:     $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: High Performance Advanced Proxy for MySQL
# Description :      High Performance and Advanced Proxy for MySQL and forks.
#                    It provides advanced features like connection pool, query routing and rewrite,
#                    firewalling, throttling, real time analysis, error-free failover
### END INIT INFO

OLDDATADIR="%s/run"
DATADIR="%s/data"
OPTS="-c %s/etc/proxysql.cnf -D $DATADIR"
PIDFILE="%s/data/proxysql.pid"

ulimit -n 102400
ulimit -c 1073741824

getpid() {
  if [ -f $PIDFILE ]
  then
	if [ -r $PIDFILE ]
	then
	  pid=$(cat $PIDFILE)
	  if [ "X$pid" != "X" ]
	  then
		# Verify that a process with this pid is still running.
		pid=$(ps -p $pid | grep $pid | grep -v grep | awk '{print $1}' | tail -1)
		if [ "X$pid" = "X" ]
		then
		  # This is a stale pid file.
			rm -f $PIDFILE
		  echo "Removed stale pid file: $PIDFILE"
		fi
	  fi
	else
	  echo "Cannot read $PIDFILE."
	  exit 1
	fi
  fi
}


testpid() {
	pid=$(ps -p $pid | grep $pid | grep -v grep | awk '{print $1}' | tail -1)
	if [ "X$pid" = "X" ]
	then
	# Process is gone so remove the pid file.
		rm -f $PIDFILE
	fi
}

initial() {
	OPTS="--initial $OPTS"
	start
}

reload() {
	OPTS="--reload $OPTS"
	start
}

start() {
  echo -n "Starting ProxySQL: "
	mkdir $DATADIR 2>/dev/null
  getpid
  if [ "X$pid" = "X" ]
   then
		if [ -f $OLDDATADIR/proxysql.db ]
		then
			if [ ! -f $DATADIR/proxysql.db ]
			then
				mv -iv $OLDDATADIR/proxysql.db $DATADIR/proxysql.db
			fi
		fi
%s $OPTS
		if [ "$?" = "0" ]; then
			echo "DONE!"
			return 0
		else
			echo "FAILED!"
			return 1
		fi
	else
		echo "ProxySQL is already running."
		exit 0
	fi
}

stop() {
  echo -n "Shutting down ProxySQL: "
  getpid
  if [ "X$pid" = "X" ]
	then
		echo "ProxySQL was not running."
		exit 0
	else
		# Note: we send a kill to all the processes, not just to the child
		for i in $(echo $pid; pgrep -x proxysql -P $pid) ; do
			if [ "$i" != "$$" ]; then
				kill $i
			fi
		done
	#  Loop until it does.
		savepid=$pid
		CNT=0
		TOTCNT=0
		while [ "X$pid" != "X" ]
		do
			# Loop for up to 20 second
			if [ "$TOTCNT" -lt "200" ]
			then
				if [ "$CNT" -lt "10" ]
				then
					CNT=$(expr $CNT + 1)
				else
					echo -n "."
					CNT=0
				fi
				TOTCNT=$(expr $TOTCNT + 1)

				sleep 0.1

				testpid
			else
				pid=
			fi
		done
		pid=$savepid
		testpid
		if [ "X$pid" != "X" ]
		then
			echo
			echo "Timed out waiting for ProxySQL to exit."
			echo "  Attempting a forced exit..."
			for i in $(echo $pid; pgrep proxysql -P $pid) ; do
				if [ "$i" != "$$" ]; then
					kill -9 $i
				fi
			done
		fi

		pid=$savepid
		testpid
		if [ "X$pid" != "X" ]
		then
			echo "Failed to stop ProxySQL"
			exit 1
		else
			echo "DONE!"
		fi
	fi
}


status() {
  getpid
  if [ "X$pid" = "X" ]
  then
		echo "ProxySQL is not running."
		exit 1
  else
		echo "ProxySQL is running ($pid)."
	exit 0
  fi
}

case "$1" in
	start)
		start
	;;
	initial)
		initial
	;;
	reload)
		reload
	;;
	stop)
		stop
	;;
	status)
		status
	;;
	restart)
		stop
		start
	;;
	*)
		echo "Usage: ProxySQL {start|stop|status|reload|restart|initial}"
		exit 1
	;;
esac
exit $?
`

//SetProxySQLConfigEx 设置porxysql配置
func SetProxySQLConfigEx(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	serviceName, err := executor.ExtractCmdFuncStringParam(params, osservice.CmdParamServiceName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	proxyConfPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamProxySQLConfPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	basicPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamProxySQLPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	binPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamProxySQLBinPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	newConfText := fmt.Sprintf(proxySQLConfigText, serviceName, basicPath, basicPath, basicPath, basicPath, binPath)
	newParams := executor.ExecutorCmdParams{}
	newParams[oscmd.CmdParamFilename] = proxyConfPath
	newParams[oscmd.CmdParamOutText] = newConfText
	newParams[oscmd.CmdParamOverwrite] = true
	return oscmd.TextToFile(e, &newParams)
}

// SetProxySQLConfig 配置ProxySQL
func SetProxySQLConfig(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	proxyConfPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamProxySQLConfPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	basicPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamProxySQLPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	binPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamProxySQLBinPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	runPath := basicPath + "/run"
	text := "OLDDATADIR=\"" + runPath + "\""
	commond := fmt.Sprintf("sed -i 's,^OLDDATADIR.*,%s,' %s", text, proxyConfPath)
	es, err := e.ExecShell(commond)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	dataPath := basicPath + "/data"
	text = "DATADIR=\"" + dataPath + "\""
	commond = fmt.Sprintf("sed -i 's,^DATADIR.*,%s,' %s", text, proxyConfPath)
	es, err = e.ExecShell(commond)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	text = "OPTS=\"-c " + basicPath + "/etc/proxysql.cnf -D $DATADIR\""
	commond = fmt.Sprintf("sed -i 's,^OPTS.*,%s,' %s", text, proxyConfPath)
	es, err = e.ExecShell(commond)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	text = binPath + " $OPTS"
	commond = fmt.Sprintf("sed -i '84s,^.*$,%s,' %s", text, proxyConfPath)
	es, err = e.ExecShell(commond)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("")
}

// ProxyMySQLCmdSQL 执行sql语句
func ProxyMySQLCmdSQL(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdSQL, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlCmdSQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port := int(0)
	host := "localhost"

	if port, err = executor.ExtractCmdFuncIntParam(params, CmdParamPort); err != nil || port == 0 {
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

	passwordStr := "-p" + strings.TrimSpace(mysqlPassword)
	hostStr := fmt.Sprintf("-h%s -P%d ", host, port)
	userStr := "-u" + mysqlUser

	cmdStr := fmt.Sprintf("%s/bin/mysql --protocol=tcp --default-auth=mysql_native_password %s %s %s --prompt \"%s> \" -e \"%s\" ", mysqlPath, hostStr, userStr, passwordStr, mysqlUser, utils.EscapeShellCmd(cmdSQL))
	es, err := e.ExecShell(cmdStr)
	log.Debug("cmdStr string=%s", cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	err = executor.GetExecResult(es)
	if err != nil {
		log.Debug("executor sql failed=%v", err)
		return executor.ErrorExecuteResult(err)
	}

	return executor.SuccessulExecuteResult(es, false, fmt.Sprintf("proxsql cmd sql statmtent execute successful"))
}

// ProxyMySQLCmdSQLEx 执行proxysql sql语句扩展函数
func ProxyMySQLCmdSQLEx(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdSQL, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlCmdSQL)
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

	err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, cmdSQL, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("proxsql cmd sql statmtent execute successful")
}

// MySQLCmdSQLEx 执行mysql sql语句扩展函数
func MySQLCmdSQLEx(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdSQL, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlCmdSQL)
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

	err = ExecuteMySQLCommand(host, mysqlUser, mysqlPassword, cmdSQL, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("mysql cmd sql statmtent execute successful")
}

// CmdSQLWithResult 执行mysql sql语句传回查询值
func CmdSQLWithResult(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdSQL, err := executor.ExtractCmdFuncStringParam(params, CmdParamMysqlCmdSQL)
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

	es, err := e.ExecShell("hostname")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	result, err := ExecuteMySQLCommandWithResult(host, mysqlUser, mysqlPassword, cmdSQL, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("result=====", result)
	er := executor.SuccessulExecuteResult(es, true, fmt.Sprintf(" excute sql %s successful", cmdSQL))
	er.ResultData = map[string]string{
		"1": result,
	}
	return er
}

// RegistOrchestrator 主从注册高可用
func RegistOrchestrator(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncStringParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	hostname, err := executor.ExtractCmdFuncStringParam(params, CmdParamHostName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPort, err := executor.ExtractCmdFuncStringParam(params, CmdParamMySQLPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	username, _ := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamUserName)
	password, _ := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamPassword)
	var basicAuthHead string
	if username != "" && password != "" {
		basicAuthHead = http.GetBasicAuthHeadInfo(username, password)
	}

	//注册api地址格式 http://192.168.0.119:3000/api/discover
	requestURL := fmt.Sprintf("http://%s:%s/api/discover/%s/%s", host, port, hostname, mysqlPort)
	method := "get"
	log.Debug("requestURL = %s and head info = %s ", requestURL, basicAuthHead)
	newParams := executor.ExecutorCmdParams{
		http.CmdParamURL:    requestURL,
		http.CmdParamMethod: method,
		http.CmdParamHead:   basicAuthHead,
	}

	es := http.HttpRequest(e, &newParams)
	err = GetOrchestratorAPIResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return es

}

//OrchestratorSwitchMaster 调用Orchestrator 接口切换主从
func OrchestratorSwitchMaster(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncStringParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	hostname, err := executor.ExtractCmdFuncStringParam(params, CmdParamHostName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPort, err := executor.ExtractCmdFuncStringParam(params, CmdParamMySQLPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	clusterName, err := executor.ExtractCmdFuncStringParam(params, CmdParamClusterName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	username, _ := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamUserName)
	password, _ := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamPassword)
	var basicAuthHead string
	if username != "" && password != "" {
		basicAuthHead = http.GetBasicAuthHeadInfo(username, password)
	}

	//切换主从http://192.168.11.200:3000/api/graceful-master-takeover/cluster01/node-1/18923
	requestURL := fmt.Sprintf("http://%s:%s/api/graceful-master-takeover/%s/%s/%s", host, port, clusterName, hostname, mysqlPort)
	method := "get"
	log.Debug("requestURL = %s and head info = %s", requestURL, basicAuthHead)
	newParams := executor.ExecutorCmdParams{
		http.CmdParamURL:    requestURL,
		http.CmdParamMethod: method,
		http.CmdParamHead:   basicAuthHead,
	}
	es := http.HttpRequest(e, &newParams)
	err = GetOrchestratorAPIResult(es)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return es
}

//SetProxySQLMasterOffLineSoft 设置porxysql配置
func SetProxySQLMasterOffLineSoft(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to set ProxySQL Master OffLineSoft")
	}

	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	middlewareType, err := executor.ExtractCmdFuncIntParam(params, CmdParamMiddlewareType)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//高可用类型为 HAProxy
	if middlewareType == HighAvailabilityForHAProxy {
		servicePort, err := executor.ExtractCmdFuncIntParam(params, CmdParamHAProxyPort)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		confText, err := executor.ExtractCmdFuncStringParam(params, CmdParamHAProxyConfText)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = stopHAProxyMySQLProxy(host, "", uint(servicePort))
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = reloadHAProxy(host, "", confText, "/etc/haproxy/haproxy.cfg")
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		return executor.SuccessulExecuteResultNoData("SetProxySQLMasterOffLineSoft execute successful")
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

	masterHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMasterHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	masterPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMasterPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	updateSQL := "update mysql_servers set status='OFFLINE_SOFT' where hostname='%s' and port=%d;"
	updateSQL = fmt.Sprintf(updateSQL, masterHost, masterPort)
	loadSQL := "load mysql servers to runtime;"
	//saveSQL := "save mysql servers to disk"
	// sql := fmt.Sprintf("%s%s%s", updateSQL, loadSQL, saveSQL)
	sql := fmt.Sprintf("%s%s", updateSQL, loadSQL)
	log.Debug("update sql = %s", sql)
	err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	doOffLineHard := true
	for index := 0; index <= 50; index++ {
		querySQL := "SELECT IFNULL(SUM(ConnUsed),0) FROM stats_mysql_connection_pool WHERE status='OFFLINE_SOFT' AND srv_host='%s' and srv_port=%d;"
		querySQL = fmt.Sprintf(querySQL, masterHost, masterPort)
		connPool := -1
		log.Debug("querySQL sql = %s", querySQL)
		err = ExecutePorxySQLCommandQuery(host, mysqlUser, mysqlPassword, querySQL, port, &connPool)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		log.Debug("current connPool = %d", connPool)
		if connPool == 0 {
			doOffLineHard = false
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	log.Debug("doOffLineHard = %v", doOffLineHard)
	if doOffLineHard {
		updateSQL = fmt.Sprintf("update mysql_servers set status='OFFLINE_HARD' where hostname='%s' and port=%d;", masterHost, masterPort)
		// sql = fmt.Sprintf("%s%s%s", updateSQL, loadSQL, saveSQL)
		sql = fmt.Sprintf("%s%s", updateSQL, loadSQL)
		log.Debug("updateSQL sql = %s", sql)
		err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
	}
	return executor.SuccessulExecuteResultNoData("SetProxySQLMasterOffLineSoft execute successful")
}

//SetProxySQLMasterOffLineSoft 设置MGR porxysql配置
func SetProxySQLMGRMasterOffLineSoft(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to set ProxySQL Master OffLineSoft")
	}

	hostList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamHostList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	portList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamPortList, ",")
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

	masterHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMasterHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	masterPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMasterPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//update mysql_servers set hostgroup_id=3 where hostgroup_id=2;
	//update mysql_servers set hostgroup_id=2 where hostname='192.168.11.177' and port=20900;
	//update mysql_servers set status='OFFLINE_SOFT' where hostname='192.168.11.177' and port=20900;
	//load mysql servers to runtime;
	for k, host := range hostList {
		port, _ := strconv.Atoi(portList[k])

		updateSQL := "update mysql_servers set hostgroup_id=3 where hostgroup_id=2; update mysql_servers set hostgroup_id=2 where hostname='%s' and port=%d; update mysql_servers set status='OFFLINE_SOFT' where hostname='%s' and port=%d;"
		updateSQL = fmt.Sprintf(updateSQL, masterHost, masterPort, masterHost, masterPort)
		loadSQL := "load mysql servers to runtime;"
		//saveSQL := "save mysql servers to disk"
		// sql := fmt.Sprintf("%s%s%s", updateSQL, loadSQL, saveSQL)
		sql := fmt.Sprintf("%s%s", updateSQL, loadSQL)
		err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		doOffLineHard := true
		for index := 0; index <= 5; index++ {
			querySQL := "SELECT IFNULL(SUM(ConnUsed),0) FROM stats_mysql_connection_pool WHERE status='OFFLINE_SOFT' AND srv_host='%s' and srv_port=%d;"
			querySQL = fmt.Sprintf(querySQL, masterHost, masterPort)
			connPool := -1
			err = ExecutePorxySQLCommandQuery(host, mysqlUser, mysqlPassword, querySQL, port, &connPool)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			log.Debug("current connPool = %d", connPool)
			if connPool == 0 {
				doOffLineHard = false
				break
			}
			time.Sleep(time.Second * 1)
		}
		log.Debug("doOffLineHard = %v", doOffLineHard)
		if doOffLineHard {
			updateSQL = fmt.Sprintf("update mysql_servers set status='OFFLINE_HARD' where hostname='%s' and port=%d;", masterHost, masterPort)
			// sql = fmt.Sprintf("%s%s%s", updateSQL, loadSQL, saveSQL)
			sql = fmt.Sprintf("%s%s", updateSQL, loadSQL)
			err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
		}
	}

	return executor.SuccessulExecuteResultNoData("SetProxySQLMasterOffLineSoft execute successful")
}

//SetProxySQLMemberOnLine 设置porxysql member online
func SetProxySQLMemberOnLine(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to set ProxySQL Master OffLineSoft")
	}

	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	middlewareType, err := executor.ExtractCmdFuncIntParam(params, CmdParamMiddlewareType)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//高可用类型为 HAProxy
	if middlewareType == HighAvailabilityForHAProxy {
		newConfText, err := executor.ExtractCmdFuncStringParam(params, CmdParamHAProxyNewConfText)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		err = reloadHAProxy(host, "", newConfText, "/etc/haproxy/haproxy.cfg")
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		return executor.SuccessulExecuteResultNoData("SetProxySQLMasterOffLineSoft execute successful")
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

	memberHost, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberHost, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(memberHost) != 2 {
		return executor.ErrorExecuteResult(fmt.Errorf("必须同时输入新旧两个节点实例信息"))
	}

	memberPort, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberPort, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(memberHost) != 2 {
		return executor.ErrorExecuteResult(fmt.Errorf("必须同时输入新旧两个节点实例信息"))
	}

	slaveReadAllowMaxDelay, err := executor.ExtractCmdFuncStringParam(params, CmdParamSlaveReadAllowMaxDelay)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	masterReadWeight, err := executor.ExtractCmdFuncStringParam(params, CmdParamsMasterReadWeight)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	slaveReadWeight, err := executor.ExtractCmdFuncStringParam(params, CmdParamSlaveReadWeight)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	var sql string
	//默认第一个为旧主host
	for index := range memberHost {
		oldMaster := false
		if index == 0 {
			oldMaster = true
		}

		//针对新主节点
		//replace into mysql_servers(hostgroup_id,hostname,port,status) values (2,'192.168.11.221',1541,'ONLINE');
		//针对所有节点
		//replace into mysql_servers(hostgroup_id,hostname,port,status) values (3,'192.168.11.221',1541,'ONLINE');
		//replace into mysql_servers(hostgroup_id,hostname,port,status) values (3,'192.168.11.222',1541,'ONLINE');
		//针对原主节点
		//delete from mysql_servers where hostgroup_id=2 and hostname='192.168.11.222' and port=1541;
		if oldMaster {
			sql += fmt.Sprintf("delete from mysql_servers where hostgroup_id=2 and hostname='%s' and port=%s;", memberHost[index], memberPort[index])
			sql += fmt.Sprintf("replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);", memberHost[index], memberPort[index], slaveReadWeight, slaveReadAllowMaxDelay)

		} else {
			sql += fmt.Sprintf("replace into mysql_servers(hostgroup_id,hostname,port,status) values (2,'%s',%s,'ONLINE');", memberHost[index], memberPort[index])
			//replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);
			sql += fmt.Sprintf("replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);", memberHost[index], memberPort[index], masterReadWeight, "0")
		}
	}
	sql += "load mysql servers to runtime;"
	sql += "save mysql servers to disk;"

	err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("SetProxySQLMemberOnLine execute successful")
}

//OrchestratorStartReplicaption 调用Orchestrator 启动主从复制
func OrchestratorStartReplicaption(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncStringParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	hostname, err := executor.ExtractCmdFuncStringParam(params, CmdParamHostName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	mysqlPort, err := executor.ExtractCmdFuncStringParam(params, CmdParamMySQLPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	username, _ := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamUserName)
	password, _ := executor.ExtractCmdFuncStringParam(params, oscmd.CmdParamPassword)
	var basicAuthHead string
	if username != "" && password != "" {
		basicAuthHead = http.GetBasicAuthHeadInfo(username, password)
	}

	var er executor.ExecuteResult
	for index := 0; index < 5; index++ {
		//启动主从复制http://127.0.0.1:3000/api/start-replica/192.168.0.200/20213
		requestURL := fmt.Sprintf("http://%s:%s/api/start-replica/%s/%s", host, port, hostname, mysqlPort)
		method := "get"
		log.Debug("requestURL = %s and head info = %s", requestURL, basicAuthHead)
		newParams := executor.ExecutorCmdParams{
			http.CmdParamURL:    requestURL,
			http.CmdParamMethod: method,
			http.CmdParamHead:   basicAuthHead,
		}
		reqRest := http.HttpRequest(e, &newParams)
		err = GetOrchestratorAPIResult(reqRest)
		if err == nil {
			return reqRest
		}
		//if reqRest.Successful {
		//	return reqRest
		//}
		er = reqRest
		time.Sleep(time.Second)
	}
	return er
}

//KillMasterInTransThread kill主节点所有的事务
func KillMasterInTransThread(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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

	err = KillThreads(host, mysqlUser, mysqlPassword, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("KillMasterInTransThread execute successful")
}

//GetKeepalivedRole 获取keepalived角色
func GetKeepalivedRole(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := fmt.Sprintf("busctl tree org.keepalived.Vrrp1")
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		log.Debug("busctl tree org.keepalived.Vrrp1 failed ")
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		role := 0
		keepalivedENOText := ""
		if len(es.Stdout) == 8 {
			textList := strings.Split(strings.TrimSpace(es.Stdout[6]), "org")
			if len(textList) == 2 {
				lastValue := textList[1]
				if lastValue != "" {
					keepalivedENOText = "/org" + lastValue
				}
				log.Debug("keepalivedENOText = %s", keepalivedENOText)
			}
		}
		log.Debug("keepalivedENOText = ", keepalivedENOText)
		if keepalivedENOText != "" {
			cmdstr = fmt.Sprintf("busctl introspect org.keepalived.Vrrp1 %s org.keepalived.Vrrp1.Instance", keepalivedENOText)
			log.Debug("commond = %s", cmdstr)
			es, err = e.ExecShell(cmdstr)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			if es.ExitCode == 0 && len(es.Stderr) == 0 {
				if len(es.Stdout) == 5 {
					if strings.Contains(es.Stdout[3], "Master") {
						log.Debug("current keeaplived is master")
						role = 1
					} else if strings.Contains(es.Stdout[3], "Backup") {
						log.Debug("current keepalived is backup")
						role = 2
					} else if strings.Contains(es.Stdout[3], "Fault") {
						log.Debug("current keeaplived is fault")
						role = 3
					} else {
						log.Debug("current keeaplived unknown role")
					}
				}
			}
		}

		er := executor.SuccessulExecuteResult(es, false, "get keepalived role successul")
		er.ResultData = make(map[string]string)
		er.ResultData["role"] = fmt.Sprintf("%d", role)
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Contains(errMsg, "Failed to introspect object") {
			er := executor.SuccessulExecuteResult(es, false, "get keepalived role successul")
			er.ResultData = make(map[string]string)
			er.ResultData["role"] = "0"
			log.Debug("current keeaplived unknown role")
			return er
		}
		log.Debug("errMsg = ", errMsg)
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//SetSalveDelay 设置从库延迟
func SetSalveDelay(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	doSlaveDelay, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoSlaveDelay)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doSlaveDelay {
		return executor.SuccessulExecuteResultNoData("No need to set salve delay")
	}

	delayTime, err := executor.ExtractCmdFuncIntParam(params, CmdParamDelayTime)
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

	// stop slave sql_thread;
	// change master to master_delay=3600;
	// start slave sql_thread;
	slaveDelaySQL := "stop slave sql_thread;change master to master_delay=%d;start slave sql_thread;"

	sql := fmt.Sprintf(slaveDelaySQL, delayTime)

	err = ExecuteMySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("SetSalveDelay execute successful")
}

//MGRAddMemberToProxySQL mgr添加porxysql member
func MGRAddMemberToProxySQL(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	log.Debug("AddMemberToProxySQL command come in")
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to add member to proxysql")
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

	memberHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMemberHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMemberPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// insert into mysql_servers (hostgroup_id,hostname,port) values (3,'192.168.11.103',3318);
	// delete from mysql_servers where hostname="192.168.11.223";
	deleteSQL := "delete from mysql_servers where hostname='%s';"
	deleteSQL = fmt.Sprintf(deleteSQL, memberHost)
	insertSQL := "insert into mysql_servers(hostgroup_id, hostname, port, status) values(3, '%s', %d, 'ONLINE');"
	insertSQL = fmt.Sprintf(insertSQL, memberHost, memberPort)
	loadSQL := "load mysql servers to runtime;"
	saveSQL := "save mysql servers to disk"
	sql := fmt.Sprintf("%s%s%s%s", deleteSQL, insertSQL, loadSQL, saveSQL)
	log.Debug("MGRAddMemberToProxySQL sql list = %v", sql)
	err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("AddMemberToProxySQL execute successful")
}

//AddMemberToProxySQL 添加porxysql member
func AddMemberToProxySQL(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	log.Debug("AddMemberToProxySQL command come in")
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to add member to proxysql")
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

	memberHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMemberHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMemberPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	slaveReadAllowMaxDelay, err := executor.ExtractCmdFuncStringParam(params, CmdParamSlaveReadAllowMaxDelay)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	slaveReadWeight, err := executor.ExtractCmdFuncStringParam(params, CmdParamSlaveReadWeight)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// insert into mysql_servers (hostgroup_id,hostname,port) values (3,'192.168.11.103',3318);
	// delete from mysql_servers where hostname="192.168.11.223";
	deleteSQL := "delete from mysql_servers where hostname='%s';"
	deleteSQL = fmt.Sprintf(deleteSQL, memberHost)
	insertSQL := "replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%d,'ONLINE',%s,%s);"
	insertSQL = fmt.Sprintf(insertSQL, memberHost, memberPort, slaveReadWeight, slaveReadAllowMaxDelay)
	loadSQL := "load mysql servers to runtime;"
	saveSQL := "save mysql servers to disk"
	sql := fmt.Sprintf("%s%s%s%s", deleteSQL, insertSQL, loadSQL, saveSQL)
	log.Debug("AddMemberToProxySQL sql list = %v", sql)
	err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("AddMemberToProxySQL execute successful")
}

//AddMembersToProxySQL 批量添加proxysql member
func AddMembersToProxySQL(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

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

	memberHost, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberHost, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPort, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberPort, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberRole, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberRole, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	count := len(memberHost)
	if len(memberPort) != count || len(memberRole) != count {
		return executor.ErrorExecuteResult(fmt.Errorf("输入参数不合法"))
	}

	slaveReadAllowMaxDelay, err := executor.ExtractCmdFuncStringParam(params, CmdParamSlaveReadAllowMaxDelay)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	masterReadWeight, err := executor.ExtractCmdFuncStringParam(params, CmdParamsMasterReadWeight)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	slaveReadWeight, err := executor.ExtractCmdFuncStringParam(params, CmdParamSlaveReadWeight)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//2--write || 3--read
	for index := 0; index < count; index++ {
		sql := ""
		writeSetSQL := ""
		masterSetReadSQL := ""
		slaveSetReadSQL := ""
		if memberRole[index] == MasterText {
			//replace into mysql_servers(hostgroup_id,hostname,port,status) values (2,'192.168.11.222',1541,'ONLINE');
			writeSetSQL = "replace into mysql_servers(hostgroup_id, hostname, port, status) values(2, '%s', %s, 'ONLINE');"
			writeSetSQL = fmt.Sprintf(writeSetSQL, memberHost[index], memberPort[index])
			sql += writeSetSQL

			//replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);
			//主节点读延迟设置为0, 延迟是从节点相对主节点的延迟
			masterSetReadSQL = "replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);"
			masterSetReadSQL = fmt.Sprintf(masterSetReadSQL, memberHost[index], memberPort[index], masterReadWeight, "0")
			sql += masterSetReadSQL

		} else {
			//replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);
			slaveSetReadSQL = "replace into mysql_servers(hostgroup_id,hostname,port,status,weight,max_replication_lag) values (3,'%s',%s,'ONLINE',%s,%s);"
			slaveSetReadSQL = fmt.Sprintf(slaveSetReadSQL, memberHost[index], memberPort[index], slaveReadWeight, slaveReadAllowMaxDelay)
			sql += slaveSetReadSQL
		}

		loadSQL := "load mysql servers to runtime;"
		sql += loadSQL
		saveSQL := "save mysql servers to disk;"
		sql += saveSQL

		err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
	}
	return executor.SuccessulExecuteResultNoData("AddMembersToProxySQL execute successful")
}

//RemoveMemberFromProxySQL 移除 porxysql member
func RemoveMemberFromProxySQL(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	log.Debug("RemoveMemberFromProxySQL command come in")
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to add member to proxysql")
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

	memberHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMemberHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	// memberPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMemberPort)
	// if err != nil {
	// 	return executor.ErrorExecuteResult(err)
	// }

	// insert into mysql_servers (hostgroup_id,hostname,port) values (3,'192.168.11.103',3318);
	// delete from mysql_servers where hostname="192.168.11.223";
	deleteSQL := "delete from mysql_servers where hostname='%s';"
	deleteSQL = fmt.Sprintf(deleteSQL, memberHost)
	// insertSQL := "insert into mysql_servers(hostgroup_id, hostname, port, status) values(3, '%s', %d, 'ONLINE');"
	// insertSQL = fmt.Sprintf(insertSQL, memberHost, memberPort)
	loadSQL := "load mysql servers to runtime;"
	saveSQL := "save mysql servers to disk"
	// sql := fmt.Sprintf("%s%s%s%s", deleteSQL, insertSQL, loadSQL, saveSQL)
	sql := fmt.Sprintf("%s%s%s", deleteSQL, loadSQL, saveSQL)
	log.Debug("RemoveMemberFromProxySQL sql list = %v", sql)
	err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, sql, port)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData("RemoveMemberFromProxySQL execute successful")
}

//ClusterHasMaster mysql cluster has master or not
func ClusterHasMaster(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	memberHost, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberHost, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPort, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberPort, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberUser, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberUser, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPassword, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberPassword, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	count := len(memberHost)
	if len(memberPort) != count || len(memberHost) != count ||
		len(memberUser) != count || len(memberPassword) == 0 {
		return executor.ErrorExecuteResult(fmt.Errorf("输入参数不合法"))
	}

	masterExist := false
	for index := 0; index < count; index++ {
		port, err := strconv.Atoi(memberPort[index])
		if err != nil {
			return executor.ErrorExecuteResult(fmt.Errorf("输入参数不合法"))
		}

		instance := NewInstanceInfo(memberHost[index], memberUser[index], memberPassword[index], port)
		var version string
		version, err = getMySQLVersion(instance)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		masterExist, err = existMasterInfo(instance, version)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}

		if masterExist {
			break
		}
	}

	er := executor.SuccessulExecuteResultNoData("MySQLClusterHasMaster execute successful")
	er.ResultData = make(map[string]string)
	er.ResultData["hasMaster"] = fmt.Sprintf("%t", masterExist)
	return er
}

//SelectMySQLClusterMasterInfo select mysql cluster master info
func SelectMySQLClusterMasterInfo(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	memberHost, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberHost, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPort, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberPort, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberUser, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberUser, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPassword, err := executor.ExtractCmdFuncStringListParam(params, CmdParamMemberPassword, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	count := len(memberHost)
	var instanceInfoList []*InstanceInfo
	var master *InstanceInfo
	var masterExist bool
	for index := 0; index < count; index++ {
		port, err := strconv.Atoi(memberPort[index])
		if err != nil {
			return executor.ErrorExecuteResult(fmt.Errorf("输入参数不合法"))
		}
		instanceInfo := NewInstanceInfo(memberHost[index], memberUser[index], memberPassword[index], port)
		instanceInfoList = append(instanceInfoList, instanceInfo)

	}

	master, err = getMaster(instanceInfoList)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if master != nil {
		masterExist = true
	}

	er := executor.SuccessulExecuteResultNoData("MySQLClusterHasMaster execute successful")
	er.ResultData = make(map[string]string)
	er.ResultData["hasMaster"] = fmt.Sprintf("%t", masterExist)
	if masterExist {
		er.ResultData["masterHost"] = fmt.Sprintf("%s", master.Host)
		er.ResultData["masterPort"] = fmt.Sprintf("%d", master.Port)
		er.ResultData["masterUsername"] = fmt.Sprintf("%s", master.Username)
		er.ResultData["masterPassword"] = fmt.Sprintf("%s", master.Password)

		var slaveInfoList []map[string]string
		for index := 0; index < count; index++ {
			if memberHost[index] == master.Host {
				continue
			} else {
				slaveInfo := make(map[string]string)
				hostKey := "slaveHost"
				slaveInfo[hostKey] = memberHost[index]

				portKey := "slavePort"
				slaveInfo[portKey] = memberPort[index]

				userKey := "slaveUsername"
				slaveInfo[userKey] = memberUser[index]

				passwordKey := "slavePassword"
				slaveInfo[passwordKey] = memberPassword[index]
				slaveInfoList = append(slaveInfoList, slaveInfo)
			}
		}
		slaveInfoListText := MapSliceToJSONString(slaveInfoList)
		er.ResultData["slaveInfo"] = slaveInfoListText
	}
	return er
}

func SwitchMGRMaster(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
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
	querySQL := "select MEMBER_ID from performance_schema.replication_group_members where MEMBER_HOST='%s' and MEMBER_PORT='%d';"
	sqlStr := fmt.Sprintf(querySQL, host, port)
	var MEMBER_ID string
	err = ExecuteMySQLCommandQuery(host, mysqlUser, mysqlPassword, sqlStr, port, false, &MEMBER_ID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if MEMBER_ID != "" {
		switchSQL := fmt.Sprintf("select group_replication_set_as_primary('%s');", MEMBER_ID)
		err = ExecuteMySQLCommand(host, mysqlUser, mysqlPassword, switchSQL, port)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
	}
	return executor.SuccessulExecuteResultNoData("SwitchMGRMaster execute successful")
}

//SetProxySQLMemberOnLine 设置porxysql member online
func ResetMGRMasterStatus(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	doProxySQL, err := executor.ExtractCmdFuncBoolParam(params, CmdParamDoProxySQL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if !doProxySQL {
		return executor.SuccessulExecuteResultNoData("No need to set ProxySQL Master OffLineSoft")
	}

	hostList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamHostList, ",")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	portList, err := executor.ExtractCmdFuncStringListParam(params, CmdParamPortList, ",")
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

	memberHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMemberHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	masterHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamMasterHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	memberPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMemberPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	masterPort, err := executor.ExtractCmdFuncIntParam(params, CmdParamMasterPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberUser, err := executor.ExtractCmdFuncStringParam(params, CmdParamMemberUser)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	memberPassword, err := executor.ExtractCmdFuncStringParam(params, CmdParamMemberPassword)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	getMasterSQL := "select MEMBER_HOST,MEMBER_PORT from performance_schema.replication_group_members where MEMBER_ROLE='PRIMARY';"
	masterInfo := struct {
		MemberHost string `orm:"column(MEMBER_HOST)"`
		MemberPort int    `orm:"column(MEMBER_PORT)"`
	}{}
	err = ExecuteMySQLCommandQuery(memberHost, memberUser, memberPassword, getMasterSQL, memberPort, false, &masterInfo)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//success//
	//update mysql_servers set hostgroup_id=3,status='ONLINE' where hostname='192.168.11.177' and port=20900;
	//update mysql_servers set hostgroup_id=2,status='ONLINE' where hostname='192.168.11.179' and port=20900;
	//load mysql servers to runtime;
	//save mysql servers to disk;
	//failure //
	//update mysql_servers set status='ONLINE';
	//load mysql servers to runtime;
	//save mysql servers to disk;
	for k, host := range hostList {
		port, _ := strconv.Atoi(portList[k])
		var CmdSQL = "update mysql_servers set status='ONLINE';load mysql servers to runtime;save mysql servers to disk;"

		if masterInfo.MemberHost == memberHost && masterInfo.MemberPort == memberPort {
			CmdSQL = fmt.Sprintf("update mysql_servers set hostgroup_id=3,status='ONLINE' where hostname='%s' and port=%d;update mysql_servers set hostgroup_id=2,status='ONLINE' where hostname='%s' and port=%d;load mysql servers to runtime;save mysql servers to disk;", masterHost, masterPort, memberHost, memberPort)
		}

		err = ExecuteProxySQLCommand(host, mysqlUser, mysqlPassword, CmdSQL, port)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
	}

	return executor.SuccessulExecuteResultNoData("ResetMGRMasterStatus execute successful")

}

//MapSliceToJSONString  map slice to json string
func MapSliceToJSONString(param []map[string]string) (jsonSTR string) {
	bRes, _ := json.Marshal(param)
	jsonSTR = string(bRes)
	return
}
