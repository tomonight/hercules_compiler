package ssh

import (
	"testing"
)

//@desc 定义相关变量
var (
	Host       = "192.168.0.181"
	Username   = "root"
	Password   = "root123"
	Port       = 22
	gSSHClient = NewSSHClient(Host, Username, Password, Port)
	sourceFile = "/Users/daiwei/tmp/1.txt"
	targetPath = "/home/daiwei/tmp"
)

//@desc 测试	ssh连接
func TestSSHConnect(t *testing.T) {
	err := gSSHClient.Connect()
	if err != nil {
		t.Error("TestSSHConnect failed ", err)
	} else {
		t.Log("TestSSHConnect success!")
	}
}

//@desc 测试ssh执行命令
func TestSSHRunShell(t *testing.T) {
	res, err := gSSHClient.Run("uname -r")
	if err != nil {
		t.Error("TestSSHRunShell failed ", err)
	} else {
		t.Log("TestSSHRunShell success res = ", res)
	}
}

//@desc 测试ssh执行命令
func TestSSHRunShellEx(t *testing.T) {
	res, err := gSSHClient.RunCmdSetTimeout("ifconfig")
	if err != nil {
		t.Error("TestSSHRunShell failed ", err)
	} else {
		t.Log("TestSSHRunShell success res = ", res)
	}
}

//@desc 测试ssh 发送文件
func TestSSHSendFile(t *testing.T) {
	return
	err := gSSHClient.SendFile(sourceFile, targetPath)
	if err != nil {
		t.Error("TestSSHSendFile failed ", err)
	} else {
		t.Log("TestSSHSendFile success !")
	}
}
