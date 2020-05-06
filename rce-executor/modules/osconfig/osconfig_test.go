package osconfig

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var cmdExecutor executor.Executor

func init() {
	// initialize executor
	cmdExecutor, _ = modules.NewExecutor()
	if cmdExecutor == nil {
		fmt.Print("executor initialize error")
	}
}

func TestGetHostname(t *testing.T) {
	params := make(executor.ExecutorCmdParams)
	Convey("test get hostname", t, func() {
		er := GetHostname(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestSetHostname(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamHostname: "test",
	}
	Convey("test set hostname", t, func() {
		er := SetHostname(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestSetHostFile(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamIPAddr:   "192.168.0.149",
		CmdParamHostname: "elasticsearch",
	}
	Convey("test set host file", t, func() {
		er := SetHostsFile(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestSetTimezone(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamTimezone: "Asia/Shanghai",
	}
	Convey("test set timezone", t, func() {
		er := SetTimezone(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestGetDateTime(t *testing.T) {
	params := make(executor.ExecutorCmdParams)
	Convey("test get date time", t, func() {
		er := GetDateTime(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestSetDateTime(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamDateTime: "2006-01-02 15:04:05",
	}
	Convey("test set date time", t, func() {
		er := SetDateTime(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestSetConf(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamConfFile:    "/etc/sysctl.conf",
		CmdParamForceUpdate: true,
		CmdParamConfItem:    "net.ipv4.tcp_low_latency = 1",
	}
	Convey("test set conf", t, func() {
		er := SetSysConf(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}
