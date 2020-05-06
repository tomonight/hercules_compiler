package http

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"strings"
)

// 模块名常量定义
const (
	HttpModuleName = "http"
)

// 函数名常量定义
const (
	CmdNameHttpRequest = "HttpRequest"
)

// 命令参数常量定义
const (
	CmdParamMethod = "method"
	CmdParamBody   = "body"
	CmdParamURL    = "url"
	CmdParamHead   = "head"
)

//HttpRequest http请求
func HttpRequest(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	method, err := executor.ExtractCmdFuncStringParam(params, CmdParamMethod)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	body, _ := executor.ExtractCmdFuncStringParam(params, CmdParamBody)
	requestURL, err := executor.ExtractCmdFuncStringParam(params, CmdParamURL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	httpHead := NewHttpHead()
	head, _ := executor.ExtractCmdFuncStringParam(params, CmdParamHead)
	if head != "" {
		//头的传送方式为 key:value;key1:value1
		headList := strings.Split(head, ";")
		for _, headValue := range headList {
			currentHead := strings.Split(headValue, ":")
			if len(currentHead) == 2 {
				if currentHead[0] != "" && currentHead[1] != "" {
					httpHead[currentHead[0]] = currentHead[1]
				}
			}
		}
	}
	fmt.Println("httpHead = ", httpHead)
	res, err := httpDo(requestURL, method, body, httpHead)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	return executor.SuccessulExecuteResultNoData(res)
}

func init() {
	executor.RegisterCmd(HttpModuleName, CmdNameHttpRequest, HttpRequest)
}
