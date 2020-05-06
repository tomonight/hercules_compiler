package executor

import (
	"fmt"
	"testing"
)

//@desc 定义相关变量
var (
	host        = "139.9.235.10"
	port        = 22
	executor, _ = NewSSHAgentExecutor(host, "root", "Root_123", port)
)

func TestZCAgentExecutorListFile(t *testing.T) {
	executor.SetSudoEnabled(true)
	executor.SetExecuteUser("mysql")
	es, err := executor.ExecShell("ls -l ~;ls -l /; whoami")
	//es, err := executor.ExecShell(`/zmysql/db/mgr01/mgr0102/mysql/bin/mysql -h127.0.0.1 -P20000 -uroot -proot123 -e"SET @@SESSION.SQL_LOG_BIN=0;delete from mysql.user where user not in ('mysql.infoschema','mysql.session','mysql.sys','root') or host !='localhost';FLUSH PRIVILEGES;CREATE USER 'root'@'%' IDENTIFIED BY 'root123';ALTER USER 'root'@'localhost' IDENTIFIED BY 'root123';GRANT ALL ON *.* TO 'root'@'%' WITH GRANT OPTION ;DROP DATABASE IF EXISTS test ;FLUSH PRIVILEGES ;"`)
	if err != nil {
		t.Error("TestZCAgentExecutorListFile failed ", err)
	} else {
		t.Log("TestZCAgentExecutorListFile success")
		fmt.Printf("\n")
		fmt.Printf("ExitCode:%d\n", es.ExitCode)
		fmt.Printf("ErrorMessage:%s\n", es.ErrorMessage)

		fmt.Printf("Stdout:\n")
		for _, s := range es.Stdout {
			fmt.Printf("%s\n", s)
		}
		fmt.Printf("Stderr:\n")
		for _, s := range es.Stderr {
			fmt.Printf("%s\n", s)
		}
	}
}

func TestZCAgentClose(t *testing.T) {
	executor.Close()
}
