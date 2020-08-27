package main

import (
	"fmt"
	"hercules_compiler/engine/compile"
	_ "hercules_compiler/rce-executor/modules/cgroup"
	_ "hercules_compiler/rce-executor/modules/etcd"
	_ "hercules_compiler/rce-executor/modules/http"
	_ "hercules_compiler/rce-executor/modules/mysql"
	_ "hercules_compiler/rce-executor/modules/oscmd"
	_ "hercules_compiler/rce-executor/modules/osservice"
	_ "hercules_compiler/rce-executor/modules/zdata"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "success", args: args{path: "test_script"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.path)
		})
		noders := compile.Noders

		for _, node := range noders {
			for _, err := range node.Errors() {
				fmt.Println(err.Error())
			}
		}
		fmt.Println(noders)
	}
}

func TestRunScript(t *testing.T) {

	kk := &KK{Name: "kk", ID: 123}
	c1 := &Class{Name: "c1"}
	c2 := &Class{Name: "c2"}
	c3 := &Class{Name: "c3"}
	kk.Class = []*Class{}
	kk.Class = append(kk.Class, c1)
	kk.Class = append(kk.Class, c2)
	kk.Class = append(kk.Class, c3)
	// aa, _ := json.Marshal(kk)
	mmm := `
	{"installClusterID":24,"creator":"admin","callBackAddr":"http://192.168.11.181:9999/api/v1/task/callback/mysql/install","appID":200000,"name":"kkki","description":"安装mysql单实例: kkki","version":"8.0","clusterName":"kkki","backupPackageURL":"","mysqlPackageURL":"http://192.168.11.181:8080/download/mysql-8.0.14-linux-glibc2.12-x86_64.tar.xz","mysqlMd5":"b9a04efa353f4b16d5c034a4e6e73cc8","backupPort":0,"packageHFSAddr":"http://192.168.11.181:8080/sysfile","orchestratorHost":"","orchestratorPort":"","orchestratorHostname":"","orchestratorSSHUser":"","orchestratorSSHPassword":"","orchestratorSSHPort":"","orchestratorDatabase":"","orchestratorUser":"","orchestratorPassword":"","mytype":1,"instanceList":[{"role":1,"target":{"id":5,"host":"192.168.11.221","hostName":"zmysql221","port":22,"username":"root","password":"root123","operateSystem":"linux","osVersion":"RedHat 7.2","cpuCount":0,"cpuInfo":"","memorySize":5947148,"memoryInfo":"","kernelInfo":""},"instance_name":"kkki01","conf_file_text":"[mysql]\nprompt = [\\\\u@\\\\p][\\\\d]\u003e\\\\_\nno_auto_rehash = 1\n[mysqldump]\nsingle_transaction = 1\nquick = 1\nmax_allowed_packet = 1G\n[mysqld]\nuser = mysql\nautocommit = 1\nserver_id = 200825073\ncharacter_set_server = utf8mb4\ndatadir = /home/zmysql/db/kkki/kkki01/data\ntmpdir = /home/zmysql/db/kkki/kkki01/tmp\nsocket = /home/zmysql/db/kkki/kkki01/run/mysql.sock\ntransaction_isolation = READ-COMMITTED\nexplicit_defaults_for_timestamp = 1\nmax_allowed_packet = 1G\nevent_scheduler = ON\nopen_files_limit = 65535\ninteractive_timeout = 1800\nwait_timeout = 1800\nlock_wait_timeout = 1800\nskip_name_resolve = 1\nmax_connections = 1024\nmax_connect_errors = 1000000\ntable_open_cache = 4096\ntable_definition_cache = 4096\ntable_open_cache_instances = 64\nread_buffer_size = 8M\nread_rnd_buffer_size = 4M\nsort_buffer_size = 4M\ntmp_table_size = 32M\njoin_buffer_size = 4M\nthread_cache_size = 64\nlog_error = /home/zmysql/db/kkki/kkki01/logs/mysql-error-log.err\nlog_error_verbosity = 2\ngeneral_log_file = /home/zmysql/db/kkki/kkki01/logs/general.log\nslow_query_log = 1\nslow_query_log_file = /home/zmysql/db/kkki/kkki01/logs/slow.log\nlog_queries_not_using_indexes = 1\nlog_slow_admin_statements = 1\nlog_slow_slave_statements = 1\nlog_throttle_queries_not_using_indexes = 10\nexpire_logs_days = 10\nlong_query_time = 1\nmin_examined_row_limit = 100\nlog_bin_trust_function_creators = 1\nlog_timestamps = SYSTEM\ninnodb_buffer_pool_size = 523M\ninnodb_buffer_pool_instances = 16\ninnodb_buffer_pool_load_at_startup = 1\ninnodb_buffer_pool_dump_at_shutdown = 1\ninnodb_lru_scan_depth = 4096\ninnodb_lock_wait_timeout = 5\ninnodb_io_capacity = 2000\ninnodb_io_capacity_max = 4000\ninnodb_flush_method = O_DIRECT\ninnodb_undo_tablespaces = 3\ninnodb_flush_neighbors = 0\ninnodb_log_file_size = 1G\ninnodb_log_files_in_group = 2\ninnodb_log_buffer_size = 64M\ninnodb_purge_threads = 4\ninnodb_thread_concurrency = 64\ninnodb_print_all_deadlocks = 1\ninnodb_strict_mode = 1\ninnodb_sort_buffer_size = 128M\ninnodb_write_io_threads = 16\ninnodb_read_io_threads = 16\ninnodb_file_per_table = 1\ninnodb_stats_persistent_sample_pages = 64\ninnodb_autoinc_lock_mode = 2\ninnodb_online_alter_log_max_size = 1G\ninnodb_open_files = 65535\nloose-innodb_numa_interleave = 1\ninnodb_buffer_pool_dump_pct = 40\ninnodb_page_cleaners = 8\ninnodb_undo_log_truncate = 1\ninnodb_max_undo_log_size = 2G\ninnodb_purge_rseg_truncate_frequency = 128\ninnodb_status_output = 0\ninnodb_status_output_locks = 0\nmaster-info-repository = TABLE\nrelay_log_info_repository = TABLE\nsync_binlog = 1\ngtid-mode = ON\nenforce_gtid_consistency = 1\nbinlog_format = ROW\nbinlog_rows_query_log_events = 1\nrelay-log = /home/zmysql/db/kkki/kkki01/logs/mysql-relay\nrelay_log_recovery = 1\nslave_rows_search_algorithms = INDEX_SCAN,HASH_SCAN\nslave_parallel_type = LOGICAL_CLOCK\nslave_parallel_workers = 8\nslave_preserve_commit_order = 1\nslave_transaction_retries = 128\nbinlog_gtid_simple_recovery = 1\nlog_slave_updates = 1\nlog-bin = /home/zmysql/db/kkki/kkki01/logs/mysql-bin\ndefault_authentication_plugin = mysql_native_password\ninnodb_monitor_enable = all\nloose-mysqlx = 0\nloose_rpl_semi_sync_master_enabled = 1\nloose_rpl_semi_sync_slave_enabled = 1\nloose_rpl_semi_sync_master_timeout = 5000\nloose_rpl_semi_sync_master_wait_point = AFTER_SYNC\nloose_rpl_semi_sync_master_wait_for_slave_count = 1\nplugin_load = rpl_semi_sync_master=semisync_master.so;rpl_semi_sync_slave=semisync_slave.so\nslave_net_timeout = 4\nlower_case_table_names = 1\nbinlog_transaction_dependency_tracking = WRITESET\ntransaction_write_set_extraction = XXHASH64\nport = 19978\nreport_host = 192.168.11.221\n","install_base_dir":"/home","softInsallPath":"/home/zmysql/product","installPath":"","mySQLUser":"root","socketFilePath":"","mysql_port":19978,"mysql_root_password":"root123","install_params":{"backupPassword":"5h1VVBqP5896","backupUser":"mydata_bk","basic_dir":"/home","conf_dir":"/home/zmysql/db/kkki/kkki01/conf","datadir":"/home/zmysql/db/kkki/kkki01/data","first_conf_text":"[mysql]\nprompt = [\\\\u@\\\\p][\\\\d]\u003e\\\\_\nno_auto_rehash = 1\n[mysqldump]\nsingle_transaction = 1\nquick = 1\nmax_allowed_packet = 1G\n[mysqld]\nuser = mysql\nautocommit = 1\nserver_id = 200825073\ncharacter_set_server = utf8mb4\ndatadir = /home/zmysql/db/kkki/kkki01/data\ntmpdir = /home/zmysql/db/kkki/kkki01/tmp\nsocket = /home/zmysql/db/kkki/kkki01/run/mysql.sock\ntransaction_isolation = READ-COMMITTED\nexplicit_defaults_for_timestamp = 1\nmax_allowed_packet = 1G\nevent_scheduler = ON\nopen_files_limit = 65535\ninteractive_timeout = 1800\nwait_timeout = 1800\nlock_wait_timeout = 1800\nskip_name_resolve = 1\nmax_connections = 1024\nmax_connect_errors = 1000000\ntable_open_cache = 4096\ntable_definition_cache = 4096\ntable_open_cache_instances = 64\nread_buffer_size = 8M\nread_rnd_buffer_size = 4M\nsort_buffer_size = 4M\ntmp_table_size = 32M\njoin_buffer_size = 4M\nthread_cache_size = 64\nlog_error = /home/zmysql/db/kkki/kkki01/logs/mysql-error-log.err\nlog_error_verbosity = 2\ngeneral_log_file = /home/zmysql/db/kkki/kkki01/logs/general.log\nslow_query_log = 1\nslow_query_log_file = /home/zmysql/db/kkki/kkki01/logs/slow.log\nlog_queries_not_using_indexes = 1\nlog_slow_admin_statements = 1\nlog_slow_slave_statements = 1\nlog_throttle_queries_not_using_indexes = 10\nexpire_logs_days = 10\nlong_query_time = 1\nmin_examined_row_limit = 100\nlog_bin_trust_function_creators = 1\nlog_timestamps = SYSTEM\ninnodb_buffer_pool_size = 523M\ninnodb_buffer_pool_instances = 16\ninnodb_buffer_pool_load_at_startup = 1\ninnodb_buffer_pool_dump_at_shutdown = 1\ninnodb_lru_scan_depth = 4096\ninnodb_lock_wait_timeout = 5\ninnodb_io_capacity = 2000\ninnodb_io_capacity_max = 4000\ninnodb_flush_method = O_DIRECT\ninnodb_undo_tablespaces = 3\ninnodb_flush_neighbors = 0\ninnodb_log_file_size = 1G\ninnodb_log_files_in_group = 2\ninnodb_log_buffer_size = 64M\ninnodb_purge_threads = 4\ninnodb_thread_concurrency = 64\ninnodb_print_all_deadlocks = 1\ninnodb_strict_mode = 1\ninnodb_sort_buffer_size = 128M\ninnodb_write_io_threads = 16\ninnodb_read_io_threads = 16\ninnodb_file_per_table = 1\ninnodb_stats_persistent_sample_pages = 64\ninnodb_autoinc_lock_mode = 2\ninnodb_online_alter_log_max_size = 1G\ninnodb_open_files = 65535\nloose-innodb_numa_interleave = 1\ninnodb_buffer_pool_dump_pct = 40\ninnodb_page_cleaners = 8\ninnodb_undo_log_truncate = 1\ninnodb_max_undo_log_size = 2G\ninnodb_purge_rseg_truncate_frequency = 128\ninnodb_status_output = 0\ninnodb_status_output_locks = 0\nmaster-info-repository = TABLE\nrelay_log_info_repository = TABLE\nsync_binlog = 1\ngtid-mode = ON\nenforce_gtid_consistency = 1\nbinlog_format = ROW\nbinlog_rows_query_log_events = 1\nrelay-log = /home/zmysql/db/kkki/kkki01/logs/mysql-relay\nrelay_log_recovery = 1\nslave_rows_search_algorithms = INDEX_SCAN,HASH_SCAN\nslave_parallel_type = LOGICAL_CLOCK\nslave_parallel_workers = 8\nslave_preserve_commit_order = 1\nslave_transaction_retries = 128\nbinlog_gtid_simple_recovery = 1\nlog_slave_updates = 1\nlog-bin = /home/zmysql/db/kkki/kkki01/logs/mysql-bin\ndefault_authentication_plugin = mysql_native_password\ninnodb_monitor_enable = all\nloose-mysqlx = 0\nloose_rpl_semi_sync_master_enabled = 1\nloose_rpl_semi_sync_slave_enabled = 1\nloose_rpl_semi_sync_master_timeout = 5000\nloose_rpl_semi_sync_master_wait_point = AFTER_SYNC\nloose_rpl_semi_sync_master_wait_for_slave_count = 1\nplugin_load = rpl_semi_sync_master=semisync_master.so;rpl_semi_sync_slave=semisync_slave.so\nslave_net_timeout = 4\nlower_case_table_names = 1\nbinlog_transaction_dependency_tracking = WRITESET\ntransaction_write_set_extraction = XXHASH64\nport = 19978\nreport_host = 192.168.11.221\n","instance_dir":"/home/zmysql/db/kkki/kkki01","monitorPassword":"41Rvq3fc65e6","monitorUser":"mydata_monitor","port":"19978","replicationPassword":"3a802d24fc4f","replicationUser":"mydata_repl","server_id":"200825073","socket":"/home/zmysql/db/kkki/kkki01/run/mysql.sock"},"cpu_limit":-1,"memory_limit":-1,"version":"","available":false,"BackupSourceFlag":false,"sourcePath":"","backupPath":"","dataDir":"","backDir":"","logErrorDir":"","binlogDir":"","relayLogDir":""}],"install":true,"force":false,"recover":false,"recoverCluster":false,"backUpInfo":null,"sourceClusterArch":"","newNodeIP":"","proxysql":null,"doProxySQL":false,"delayTime":0,"extraParams":{"readFileDownloadURL":"http://192.168.11.181:8080/sysfile/readfile.zip","readFileName":"readfile","readFilePackageName":"readfile.zip","readFileVersion":"1.0","softwarePath":"/home/zmysql/product"}}
	`
	type args struct {
		name   string
		params interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{name: "backup", params: mmm}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init("test_script")
			if err := RunScript(tt.args.name, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("RunScript() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// go get -u -v github.com/nsf/gocode
// go get -u -v github.com/rogpeppe/godef
// go get -u -v github.com/golang/lint/golint
// go get -u -v github.com/lukehoban/go-find-references
// go get -u -v github.com/lukehoban/go-outline
// go get -u -v sourcegraph.com/sqs/goreturns
// go get -u -v golang.org/x/tools/cmd/gorename
// go get -u -v github.com/tpng/gopkgs
// go get -u -v github.com/newhook/go-symbols
