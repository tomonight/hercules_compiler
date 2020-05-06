package http

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"testing"
)

var cmdExecutor executor.Executor

func init() {
	//initialize executor
	cmdExecutor, _ = executor.NewLocalExecutor()
	if cmdExecutor == nil {
		fmt.Print("executor initialize error")
	}
}

func TestHttpRequest(t *testing.T) {
	params := executor.ExecutorCmdParams{
		CmdParamURL:    "http://192.168.0.201:8080",
		CmdParamMethod: "get",
	}

	result := HttpRequest(cmdExecutor, &params)
	t.Log("result=", result)
}
