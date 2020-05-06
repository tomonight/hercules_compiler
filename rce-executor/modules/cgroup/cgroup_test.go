package cgroup

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

func TestCreateNewCGroup(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamCGroupName:  "cg1",
		CmdParamControllers: "cpuset,memory",
	}
	Convey("test create new cgroup", t, func() {
		er := CreateNewCGroup(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})

}

func TestSetCGroupParams(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamCGroupName:       "cg1",
		CmdParamCGroupParamKey:   "memory.limit_in_bytes",
		CmdParamCGroupParamValue: "500m",
	}
	Convey("test set cgroup params", t, func() {
		er := SetCGroupParams(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})

	params = executor.ExecutorCmdParams{
		CmdParamCGroupName:       "cg1",
		CmdParamCGroupParamKey:   "cpuset.cpus",
		CmdParamCGroupParamValue: "0-1",
	}
	Convey("test set cgroup params", t, func() {
		er := SetCGroupParams(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}

func TestAddPIDIntoCGroup(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamCGroupName:  "cg1",
		CmdParamControllers: "cpuset,memory",
		CmdParamPID:         "47193",
	}
	Convey("test add PID into cgroup", t, func() {
		er := AddPIDIntoCGroup(cmdExecutor, &params)
		So(er.Successful, ShouldBeTrue)
	})
}
