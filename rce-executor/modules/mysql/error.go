//#定义需要忽略的错误信息
package mysql

import (
	"encoding/json"
	"fmt"
	"hercules_compiler/rce-executor/executor"
)

const (
	ErrMsg = "ERROR 3093"
)

type OrchestratorResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Details struct {
		HostName string `json:"Hostname"` //集群名字
		Port     uint   `json:"Port"`     //集群别名
	} `json:"Details"`
}

//GetOrchestratorAPIResult get orchestrator api result
func GetOrchestratorAPIResult(es executor.ExecuteResult) (err error) {
	if !es.Successful {
		err = fmt.Errorf("call orchestrator api failed, %s", es.Message)
		return
	}

	if es.Message == "" {
		return
	}

	orchestratorResponse := OrchestratorResponse{}
	err = json.Unmarshal([]byte(es.Message), &orchestratorResponse)
	if err != nil {
		return err
	}

	if orchestratorResponse.Code != "OK" {
		err = fmt.Errorf(orchestratorResponse.Message)
	}
	return
}
