package osservice

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/ssh"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrite(t *testing.T) {
	client := ssh.NewSSHClient("192.168.11.151", "root", "root123", "", 22, "")
	e, err := executor.NewSSHAgentExecutorForSSHClient(client)
	Convey("测试写入文件", t, func() {
		So(e, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
	if err == nil {
		filename := "/etc/init.d/zmysql_yyyy_yyyy01"
		textToFileParams := executor.ExecutorCmdParams{
			oscmd.CmdParamFilename: filename,
			oscmd.CmdParamOutText:  "mysqlServiceTemperlate",
		}
		er := oscmd.TextToFile(e, &textToFileParams)
		Convey("测试写入文件", t, func() {
			So(er, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	}
}

func TestSystemdServiceStatus(t *testing.T) {
	Convey("成功连接目标主机", t, func() {
		client := ssh.NewSSHClient("192.168.11.167", "root", "123456", "", 22, "")
		e, err := executor.NewSSHAgentExecutorForSSHClient(client)
		So(e, ShouldNotBeNil)
		So(err, ShouldBeNil)
		if err == nil && e != nil {
			Convey("获取mydata_service状态", func() {
				params := executor.ExecutorCmdParams{
					CmdParamServiceName: "zmysql_167_16701",
				}
				er := SystemdServiceStatus(e, &params)
				t.Log(er.ResultData, er.Message)

			})
		}
	})
}

func TestSysServiceStatus(t *testing.T) {
	Convey("成功连接目标主机", t, func() {
		client := ssh.NewSSHClient("192.168.11.151", "root", "root123", "", 22, "")
		e, err := executor.NewSSHAgentExecutorForSSHClient(client)
		So(e, ShouldNotBeNil)
		So(err, ShouldBeNil)
		if err == nil && e != nil {
			Convey("获取mydata_service状态", func() {
				params := executor.ExecutorCmdParams{
					CmdParamServiceName: "zmysql_zc66_zc6601",
				}
				er := SysServiceStatus(e, &params)
				t.Log(er, er.Message)

			})
		}
	})
}
