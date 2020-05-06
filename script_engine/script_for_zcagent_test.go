package script_engine

import (
	"fmt"
	_ "hercules_compiler/rce-executor/modules/mysql"
	_ "hercules_compiler/rce-executor/modules/oscmd"
	_ "hercules_compiler/rce-executor/modules/osservice"
	"io/ioutil"
	"strings"
	"testing"
)

var script_outfile = `
#script begin
connect target target1 protocol=zcagent host=${host} port=${port} 
set target target1
execute oscmd.TextToFile filename="/tmp/222.txt" outText="${fileText}"  byBase64="true"
`

var script_install_mysql = `#script begin
#安装单实例mysql
#host -- 节点host
#port -- 节点port
#dataBaseDir -- 数据库基础目录: /mysql/testdb/testdb1
#mysqlParameterFile -- my.cnf文件路径: /mysql/testdb/testdb1/data/my.cnf
#mysqlPath -- mysql路径: /mysql/testdb/testdb1/mysql-5.7.15
#socketFile -- socket文件: /tmp/mysql3386.sock
#url -- 文件地址: http://192.168.0.185:7981/mysql-5.7.15-linux-glibc2.5-x86_64.tar.gz
#filename -- 文件名: mysql-5.7.15-linux-glibc2.5-x86_64.tar.gz
#unzippedFilename -- 解压后文件名: mysql-5.7.15-linux-glibc2.5-x86_64
#confFileText -- 配置文件内容
#mysqlDataDir -- 数据目录: /mysql/testdb/testdb1/data
#mysqlPort -- 端口: 3386
#mysqlRootPassword -- mysql root密码
#instanceName -- 实例名称


connect target target1 protocol=zcagent host=${host} port=${port}
set target target1
set exec failed continue
set exec successful stop
execute osservice.SysVServiceControl serviceName="zcloud_mysql_lilei_test1_lilei_test101" serviceAction="start"

`

// execute oscmd.AddGroup groupName="mysql"
// execute oscmd.AddUser userName="mysql" groupName="mysql"
// execute oscmd.MakeDir path="/soft"
// execute oscmd.DownloadFile url="${url}" outputFilename="/soft/${filename}"
// execute oscmd.MakeDir path="${dataBaseDir}"
// execute oscmd.UnzipFile directory="${dataBaseDir}" filename="/soft/${filename}"
// execute oscmd.Move source="${dataBaseDir}/${unzippedFilename}" target="${mysqlPath}"
// print "starting initialize mysql instance"
// execute mysql.InitializeMySQLInstance port=${mysqlPort}  mysqlPath=${mysqlPath} dataBaseDir=${dataBaseDir} serverId=${serverId} socketFile=${socketFile} mysqlDataDir=${mysqlDataDir} user="mysql"
// execute oscmd.TextToFile filename="${mysqlParameterFile}" outText="${confFileText}"
// execute oscmd.ChangeOwnAndGroup own="mysql" group="mysql" filenamePattern="${mysqlParameterFile}"
// execute mysql.StartupMySQLInstance port=${mysqlPort} mysqlPath=${mysqlPath} mysqlDataDir=${mysqlDataDir} socketFile=${socketFile} user="mysql" mysqlParameterFile="${mysqlParameterFile}" dbName="${dbName}" instanceName="${instanceName}"
// execute mysql.StartMySQLService dbName="${dbName}" instanceName="${instanceName}"
// sleep 120000
// execute mysql.MySQLInstanceAlive mysqlPath=${mysqlPath} socketFile="${socketFile}" user="root"
// set var sqls="SET @@SESSION.SQL_LOG_BIN=0;DELETE FROM mysql.user ;CREATE USER 'root'@'%' IDENTIFIED BY '${mysqlRootPassword}' ;GRANT ALL ON *.* TO 'root'@'%' WITH GRANT OPTION ;DROP DATABASE IF EXISTS test ;FLUSH PRIVILEGES ;"
// execute mysql.MySQLCmdSQL cmdSql="${sqls}" mysqlPath=${mysqlPath} socketFile="${socketFile}" user="root"
// set exec failed continue
// set exec successful stop
func TestInstallMySQLByZCAgent(t *testing.T) {
	//script_engine.ExecuteScript(strings.Split(script_install_zminor, "\n"), map[string]string{"password": "root123", "host": "enmo.wicp.net", "port": "8188", "username": "root",
	//script_engine.ExecuteScript(strings.Split(installEtcd, "\n"), map[string]string{}, "install etcd", nil)
	ExecuteScript(strings.Split(script_install_mysql, "\n"), map[string]string{
		"host":               "192.168.88.201",
		"port":               "8100",
		"mysqlPath":          "/mysql/testdb/testdb1/mysql",
		"socketFile":         "/tmp/mysql3386.sock",
		"url":                "http://192.168.11.181:9999/userfiles/mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz",
		"filename":           "mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz",
		"unzippedFilename":   "mysql-5.7.24-linux-glibc2.12-x86_64",
		"dataBaseDir":        "/mysql/testdb/testdb1",
		"mysqlDataDir":       "/mysql/testdb/testdb1/data",
		"serverId":           "190127001",
		"mysqlPort":          "3386",
		"mysqlRootPassword":  "root123",
		"mysqlParameterFile": "/mysql/testdb/testdb1/data/my.cnf",
		"dbName":             "testdb",
		"instanceName":       "testdb1",
		"confFileText": `
		[mysqldump]
		single_transaction = 1
		quick = 1
		max_allowed_packet = 1G
		[mysqld]
		user = mysql
		sql_mode = STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO
		autocommit = 1
		server_id = 190127001
		character_set_server = utf8mb4
		datadir = /mysql/testdb/testdb1/data
		tmpdir = /mysql/testdb/testdb1/tmp
		socket = /tmp/mysql3386.sock
		transaction_isolation = READ-COMMITTED
		explicit_defaults_for_timestamp = 1
		max_allowed_packet = 1G
		event_scheduler = ON
		open_files_limit = 65535
		interactive_timeout = 1800
		wait_timeout = 1800
		lock_wait_timeout = 1800
		skip_name_resolve = 1
		max_connections = 1024
		max_connect_errors = 1000000
		table_open_cache = 4096
		table_definition_cache = 4096
		table_open_cache_instances = 64
		read_buffer_size = 8M
		read_rnd_buffer_size = 4M
		sort_buffer_size = 4M
		tmp_table_size = 32M
		join_buffer_size = 4M
		thread_cache_size = 64
		log_error = /mysql/testdb/testdb1/logs/mysql-error-log.err
		log_error_verbosity = 2
		general_log_file = /mysql/testdb/testdb1/logs/general.log
		slow_query_log = 1
		slow_query_log_file = /mysql/testdb/testdb1/logs/slow.log
		log_queries_not_using_indexes = 1
		log_slow_admin_statements = 1
		log_slow_slave_statements = 1
		log_throttle_queries_not_using_indexes = 10
		expire_logs_days = 10
		long_query_time = 1
		min_examined_row_limit = 100
		log_bin_trust_function_creators = 1
		log_timestamps = SYSTEM
		innodb_buffer_pool_size = 4G
		innodb_buffer_pool_instances = 16
		innodb_buffer_pool_load_at_startup = 1
		innodb_buffer_pool_dump_at_shutdown = 1
		innodb_lru_scan_depth = 4096
		innodb_lock_wait_timeout = 5
		innodb_io_capacity = 2000
		innodb_io_capacity_max = 4000
		innodb_flush_method = O_DIRECT
		innodb_flush_neighbors = 0
		innodb_log_file_size = 128M
		innodb_log_files_in_group = 2
		innodb_log_buffer_size = 64M
		innodb_purge_threads = 4
		innodb_thread_concurrency = 64
		innodb_print_all_deadlocks = 1
		innodb_strict_mode = 1
		innodb_sort_buffer_size = 128M
		innodb_write_io_threads = 16
		innodb_read_io_threads = 16
		innodb_file_per_table = 1
		innodb_stats_persistent_sample_pages = 64
		innodb_autoinc_lock_mode = 2
		innodb_online_alter_log_max_size = 1G
		innodb_open_files = 65535
		loose-innodb_numa_interleave = 1
		innodb_buffer_pool_dump_pct = 40
		innodb_page_cleaners = 8
		innodb_undo_log_truncate = 1
		innodb_max_undo_log_size = 2G
		innodb_purge_rseg_truncate_frequency = 128
		innodb_status_file = 1
		innodb_status_output = 0
		innodb_status_output_locks = 0
		master-info-repository = TABLE
		relay_log_info_repository = TABLE
		sync_binlog = 1
		gtid-mode = ON
		enforce_gtid_consistency = 1
		binlog_format = ROW
		binlog_rows_query_log_events = 1
		relay-log = /mysql/testdb/testdb1/logs/mysql-relay
		relay_log_recovery = 1
		slave_rows_search_algorithms = INDEX_SCAN,HASH_SCAN
		slave_parallel_type = LOGICAL_CLOCK
		slave_parallel_workers = 8
		slave_preserve_commit_order = 1
		slave_transaction_retries = 128
		binlog_gtid_simple_recovery = 1
		log_slave_updates = 1
		log-bin = /mysql/testdb/testdb1/logs/mysql-bin
		default_authentication_plugin = mysql_native_password
		innodb_monitor_enable = all
		loose-mysqlx = 0
		loose_rpl_semi_sync_master_enabled = 1
		loose_rpl_semi_sync_slave_enabled = 1
		loose_rpl_semi_sync_master_timeout = 3600000
		loose_rpl_semi_sync_master_wait_point = AFTER_SYNC
		loose_rpl_semi_sync_master_wait_for_slave_count = 1
		plugin_load = rpl_semi_sync_master=semisync_master.so;rpl_semi_sync_slave=semisync_slave.so
		slave_net_timeout = 4
		report_host = 192.168.11.171
		slave_net_timeout = 8
		read_only = 1
		super_read_only = 1
		binlog_transaction_dependency_tracking=WRITESET
		transaction_write_set_extraction=XXHASH64
		[mysql]
		prompt = [\u@\p][\d]>\_
		no_auto_rehash = 1
		`,
	}, "mysql", nil)
}

func TestOutputFileByZCAgent(t *testing.T) {
	fileText, err := ioutil.ReadFile("/Users/cm/action_group.sql")
	if err != nil {
		fmt.Errorf("error:%v\n", err)
		return
	}
	ExecuteScript(strings.Split(script_outfile, "\n"), map[string]string{
		"host":     "192.168.88.201",
		"port":     "8100",
		"fileText": string(fileText),
	}, "outputfile", nil)
}
