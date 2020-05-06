package proxysql

import (
	"database/sql"
	//"fmt"
	"log"
)

type (
	//PsStatus 定义ps status
	PsStatus struct {
		ActiveTransactions         int64 `db:"Active_Transactions" json:"Active_Transactions"`
		BackendQueryTimeNsec       int64 `db:"Backend_query_time_nsec" json:"Backend_query_time_nsec"`
		ClientConnectionsAborted   int64 `db:"Client_Connections_aborted" json:"Client_Connections_aborted"`
		ClientConnectionsConnected int64 `db:"Client_Connections_connected" json:"Client_Connections_connected"`
		ClientConnectionsCreated   int64 `db:"Client_Connections_created" json:"Client_Connections_created"`
		ClientConnectionsNonIdle   int64 `db:"Client_Connections_non_idle" json:"Client_Connections_non_idle"`
		ComAutocommit              int64 `db:"Com_autocommit" json:"Com_autocommit"`
		ComAutocommitFiltered      int64 `db:"Com_autocommit_filtered" json:"Com_autocommit_filtered"`
		ComCommit                  int64 `db:"Com_commit" json:"Com_commit"`
		ComCommitFiltered          int64 `db:"Com_commit_filtered" json:"Com_commit_filtered"`
		ComRollback                int64 `db:"Com_rollback" json:"Com_rollback"`
		ComRollbackFiltered        int64 `db:"Com_rollback_filtered" json:"Com_rollback_filtered"`
		ComStmtClose               int64 `db:"Com_stmt_close" json:"Com_stmt_close"`
		ComStmtExecute             int64 `db:"Com_stmt_execute" json:"Com_stmt_execute"`
		ComStmtPrepare             int64 `db:"Com_stmt_prepare" json:"Com_stmt_prepare"`
		ConnPoolGetConnFailure     int64 `db:"ConnPool_get_conn_failure" json:"ConnPool_get_conn_failure"`
		ConnPoolGetConnImmediate   int64 `db:"ConnPool_get_conn_immediate" json:"ConnPool_get_conn_immediate"`
		ConnPoolGetconnSuccess     int64 `db:"ConnPool_get_conn_success" json:"ConnPool_get_conn_success"`
		ConnPoolMemoryBytes        int64 `db:"ConnPool_memory_bytes" json:"ConnPool_memory_bytes"`
		MySQLMonitorWorkers        int64 `db:"MySQL_Monitor_Workers" json:"MySQL_Monitor_Workers"`
		MySQLThreadWorkers         int64 `db:"MySQL_Thread_Workers" json:"MySQL_Thread_Workers"`
		ProxySQLUptime             int64 `db:"ProxySQL_Uptime" json:"ProxySQL_Uptime"`
		QueriesBackendsBytesRecv   int64 `db:"Queries_backends_bytes_recv" json:"Queries_backends_bytes_recv"`
		QueriesBackendsBytesSent   int64 `db:"Queries_backends_bytes_sent" json:"Queries_backends_bytes_sent"`
		QueryCacheEntries          int64 `db:"Query_Cache_Entries" json:"Query_Cache_Entries"`
		QueryCacheMemoryBytes      int64 `db:"Query_Cache_Memory_bytes" json:"Query_Cache_Memory_bytes"`
		QueryCachePurged           int64 `db:"Query_Cache_Purged" json:"Query_Cache_Purged"`
		QueryCacheBytesIN          int64 `db:"Query_Cache_bytes_IN" json:"Query_Cache_bytes_IN"`
		QueryCacheBytesOUT         int64 `db:"Query_Cache_bytes_OUT" json:"Query_Cache_bytes_OUT"`
		QueryCacheCountGET         int64 `db:"Query_Cache_count_GET" json:"Query_Cache_count_GET"`
		QueryCacheCountGETOK       int64 `db:"Query_Cache_count_GET_OK" json:"Query_Cache_count_GET_OK"`
		QueryCacheCountSET         int64 `db:"Query_Cache_count_SET" json:"Query_Cache_count_SET"`
		QueryProcessorTimeNsec     int64 `db:"Query_Processor_time_nsec" json:"Query_Processor_time_nsec"`
		Questions                  int64 `db:"Questions" json:"Questions"`
		SQLite3MemoryBytes         int64 `db:"SQLite3_memory_bytes" json:"SQLite3_memory_bytes"`
		ServerConnectionsAborted   int64 `db:"Server_Connections_aborted" json:"Server_Connections_aborted"`
		ServerConnectionsConnected int64 `db:"Server_Connections_connected" json:"Server_Connections_connected"`
		ServerConnectionsCreated   int64 `db:"Server_Connections_created" json:"Server_Connections_created"`
		ServersTableVersion        int64 `db:"Servers_table_version" json:"Servers_table_version"`
		SlowQueries                int64 `db:"Slow_queries" json:"Slow_queries"`
		StmtActiveTotal            int64 `db:"Stmt_Active_Total" json:"Stmt_Active_Total"`
		StmtActiveUnique           int64 `db:"Stmt_Active_Unique" json:"Stmt_Active_Unique"`
		StmtMaxStmtID              int64 `db:"Stmt_Max_Stmt_id" json:"Stmt_Max_Stmt_id"`
		MysqlBackendBuffersBytes   int64 `db:"mysql_backend_buffers_bytes" json:"mysql_backend_buffers_bytes"`
		MysqlFrontendBuffersBytes  int64 `db:"mysql_frontend_buffers_bytes" json:"mysql_frontend_buffers_bytes"`
		MysqlSessionInternalBytes  int64 `db:"mysql_session_internal_bytes" json:"mysql_session_internal_bytes"`
	}
	//Status define status info
	Status struct {
		VariablesName string `db:"Variable_name" json:"Variable_name"`
		Value         int64  `db:"Value" json:"Value"`
	}
)

//定义mysql status
const (
	StmtMySQLStatus = `SHOW MYSQL STATUS`
)

//GetProxySQLStatus 获取proxysql status
func (ps *PsStatus) GetProxySQLStatus(db *sql.DB) PsStatus {

	var tmp Status
	rows, err := db.Query(StmtMySQLStatus)
	if err != nil {
		log.Print("db.Query", StmtMySQLStatus)
	}
	for rows.Next() {
		tmp = Status{}
		err = rows.Scan(&tmp.VariablesName, &tmp.Value)
		if err != nil {
			log.Print("err = ", err)
			return PsStatus{}
		}

		switch tmp.VariablesName {
		case "Active_Transactions":
			ps.ActiveTransactions = tmp.Value
		case "Backend_query_time_nsec":
			ps.BackendQueryTimeNsec = tmp.Value
		case "Client_Connections_aborted":
			ps.ClientConnectionsAborted = tmp.Value
		case "Client_Connections_connected":
			ps.ClientConnectionsConnected = tmp.Value
		case "Client_Connections_created":
			ps.ClientConnectionsCreated = tmp.Value
		case "Client_Connections_non_idle":
			ps.ClientConnectionsNonIdle = tmp.Value
		case "Com_autocommit":
			ps.ComAutocommit = tmp.Value
		case "Com_autocommit_filtered":
			ps.ComAutocommitFiltered = tmp.Value
		case "Com_commit":
			ps.ComCommit = tmp.Value
		case "Com_commit_filtered":
			ps.ComCommitFiltered = tmp.Value
		case "Com_rollback":
			ps.ComRollback = tmp.Value
		case "Com_rollback_filtered":
			ps.ComRollbackFiltered = tmp.Value
		case "Com_stmt_close":
			ps.ComStmtClose = tmp.Value
		case "Com_stmt_execute":
			ps.ComStmtExecute = tmp.Value
		case "Com_stmt_prepare":
			ps.ComStmtPrepare = tmp.Value
		case "ConnPool_get_conn_failure":
			ps.ConnPoolGetConnFailure = tmp.Value
		case "ConnPool_get_conn_immediate":
			ps.ConnPoolGetConnImmediate = tmp.Value
		case "ConnPool_get_conn_success":
			ps.ConnPoolGetconnSuccess = tmp.Value
		case "ConnPool_memory_bytes":
			ps.ConnPoolMemoryBytes = tmp.Value
		case "MySQL_Monitor_Workers":
			ps.MySQLMonitorWorkers = tmp.Value
		case "MySQL_Thread_Workers":
			ps.MySQLThreadWorkers = tmp.Value
		case "ProxySQL_Uptime":
			ps.ProxySQLUptime = tmp.Value
		case "Queries_backends_bytes_recv":
			ps.QueriesBackendsBytesRecv = tmp.Value
		case "Queries_backends_bytes_sent":
			ps.QueriesBackendsBytesSent = tmp.Value
		case "Query_Cache_Entries":
			ps.QueryCacheEntries = tmp.Value
		case "Query_Cache_Memory_bytes":
			ps.QueryCacheMemoryBytes = tmp.Value
		case "Query_Cache_Purged":
			ps.QueryCachePurged = tmp.Value
		case "Query_Cache_bytes_IN":
			ps.QueryCacheBytesIN = tmp.Value
		case "Query_Cache_bytes_OUT":
			ps.QueryCacheBytesOUT = tmp.Value
		case "Query_Cache_count_GET":
			ps.QueryCacheCountGET = tmp.Value
		case "Query_Cache_count_GET_OK":
			ps.QueryCacheCountGETOK = tmp.Value
		case "Query_Cache_count_SET":
			ps.QueryCacheCountSET = tmp.Value
		case "Query_Processor_time_nsec":
			ps.QueryProcessorTimeNsec = tmp.Value
		case "Questions":
			ps.Questions = tmp.Value
		case "SQLite3_memory_bytes":
			ps.SQLite3MemoryBytes = tmp.Value
		case "Server_Connections_aborted":
			ps.ServerConnectionsAborted = tmp.Value
		case "Server_Connections_connected":
			ps.ServerConnectionsConnected = tmp.Value
		case "Server_Connections_created":
			ps.ServerConnectionsCreated = tmp.Value
		case "Servers_table_version":
			ps.ServersTableVersion = tmp.Value
		case "Slow_queries":
			ps.SlowQueries = tmp.Value
		case "Stmt_Active_Total":
			ps.StmtActiveTotal = tmp.Value
		case "Stmt_Active_Unique":
			ps.StmtActiveUnique = tmp.Value
		case "Stmt_Max_Stmt_id":
			ps.StmtMaxStmtID = tmp.Value
		case "mysql_backend_buffers_bytes":
			ps.MysqlBackendBuffersBytes = tmp.Value
		case "mysql_frontend_buffers_bytes":
			ps.MysqlFrontendBuffersBytes = tmp.Value
		case "mysql_session_internal_bytes":
			ps.MysqlSessionInternalBytes = tmp.Value
		default:
			log.Print("GetProxySqlStatus()", tmp.VariablesName)
		}
	}
	log.Printf("GetProxySqlStatus = %#v", *ps)
	return *ps
}
