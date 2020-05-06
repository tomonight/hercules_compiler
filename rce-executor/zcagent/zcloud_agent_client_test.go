package zcagent

import (
	"fmt"
	"testing"
)

//@desc 定义相关变量
var (
	host        = "192.168.20.151"
	port        = 8100
	agentClient = NewClient(host, port)
)

func TestAgentConnect(t *testing.T) {
	err := agentClient.Open()
	if err != nil {
		t.Error("TestAgentConnect failed ", err)
	} else {
		t.Log("TestAgentConnect success!")
	}
}

func TestAgentGetKernel(t *testing.T) {
	res, err := agentClient.RunCmd("uname -r", 1000)
	if err != nil {
		t.Error("TestAgentGetKernel failed ", err)
	} else {
		t.Log("TestAgentGetKernel success res = ", res)
	}
}
func TestAgentGetFileSystem(t *testing.T) {
	res, err := agentClient.RunCmd("df -k", 1000)
	if err != nil {
		t.Error("TestAgentGetFileSystem failed ", err)
	} else {
		t.Log("TestAgentGetFileSystem success")
		fmt.Printf("\n")
		for _, s := range res.Stdout {
			fmt.Printf("%s\n", s)
		}
	}
}

func TestAgentDispayVar(t *testing.T) {
	res, err := agentClient.RunCmd("export ABC=12345;echo ABC=$ABC; cd /tmp;echo PWD=$PWD", 1000)
	if err != nil {
		t.Error("TestAgentDispayVar failed ", err)
	} else {
		t.Log("TestAgentDispayVar success")
		fmt.Printf("\n")
		for _, s := range res.Stdout {
			fmt.Printf("%s\n", s)
		}
	}
}

func TestAgentListFile(t *testing.T) {
	res, err := agentClient.RunCmd("ls -l /etc", 1000)
	if err != nil {
		t.Error("TestAgentListFile failed ", err)
	} else {
		t.Log("TestAgentListFile success")
		fmt.Printf("\n")
		for _, s := range res.Stdout {
			fmt.Printf("%s\n", s)
		}
	}
}
func TestAgentErrorCmd(t *testing.T) {
	res, err := agentClient.RunCmd("abc", 1000)
	if err != nil {
		t.Error("TestAgentErrorCmd failed ", err)
	} else {
		t.Log("TestAgentErrorCmd success")
		fmt.Printf("shell exit code=%d", res.ExitCode)
		fmt.Printf("\n")
		for _, s := range res.Stderr {
			fmt.Printf("%s\n", s)
		}
	}
}

func TestAgentExecTimeout(t *testing.T) {
	res, err := agentClient.RunCmd("sleep 3", 2000)
	if err != nil {
		t.Error("TestAgentExecTimeout failed ", err)
	} else {
		t.Log("TestAgentExecTimeout success")
		fmt.Printf("\n")
		fmt.Printf("shell exit code=%d", res.ExitCode)
		for _, s := range res.Stdout {
			fmt.Printf("%s\n", s)
		}
		fmt.Printf("shell stderr:\n")
		for _, s := range res.Stderr {
			fmt.Printf("%s\n", s)
		}
	}
}
func TestAgentClose(t *testing.T) {
	agentClient.Close()
}
