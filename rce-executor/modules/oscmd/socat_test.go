package oscmd

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/ssh"
	"testing"
)

//
//tar -C /home/MyData_BackUp/bms100/bms10001/binlog/20200319150045/mysql-bin.000002/  -cvf - mysql-bin.000002 | socat -b 10485760 -u stdio TCP:192.168.11.175:12808

func TestSocatTansferFiles(t *testing.T) {
	client := ssh.NewSSHClient("192.168.11.178", "root", "root123", "", 22)
	e, err := executor.NewSSHAgentExecutorForSSHClient(client)
	if err == nil {
		err = socatTransferFiles(e, "", "192.168.11.221", "192.168.11.221",
			"root", "/home/MyData_BackUp/bms100/bms10001/binlog/20200319150045/mysql-bin.000002",
			"/tmp/mytest", 22, 12808)
		t.Log("err = ", err)
	}

}
