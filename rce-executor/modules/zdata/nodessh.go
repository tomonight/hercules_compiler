package zdata

import (
	"encoding/json"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"io/ioutil"
	"os"
)

//Node 定义节点数据格式
type Node struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

//GenerateNodeRsa generate rsa, node could ssh to another node without entering password
func GenerateNodeRsa(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	var (
		cmdStr string
		es     *executor.ExecutedStatus
		err    error
	)
	//@desc get SSHRsaPubFile
	rsaPubFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamaSSHRsaPubFile)
	if err != nil {
		// default to ~/.ssh/id_rsa
		path, err := executor.GetSShPath()
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		rsaPubFile =path+ "/data/.ssh/id_rsa"
	}
	//@desc get NodeList
	nodeListStr, err := executor.ExtractCmdFuncStringParam(params, CmdParamaNodeList)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//@desc 解析nodeListstr JsonString
	var nodeList []*Node
	err = json.Unmarshal([]byte(nodeListStr), &nodeList)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//@desc get ssh_ras_pub file content
	file, err := os.Open(rsaPubFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	bRes, err := ioutil.ReadAll(file)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	body := string(bRes)

	for _, node := range nodeList {
		e, err := executor.NewSSHAgentExecutor(node.Host, node.UserName, node.Password, node.Port)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		defer e.Close()

		cmdStr = "mkdir -p ~/.ssh/"
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}

		cmdStr = fmt.Sprintf("echo \"%s\" > ~/.ssh/authorized_keys", body)
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}

		cmdStr = "chmod 644 ~/.ssh/authorized_keys"
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}

		cmdStr = "chmod 700 ~/.ssh/"
		es, err = e.ExecShell(cmdStr)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}
		err = executor.GetExecResult(es)
		if err != nil {
			log.Warn("cmd %s failed %v", cmdStr, err)
			return executor.ErrorExecuteResult(err)
		}
	}
	return executor.SuccessulExecuteResult(es, true, "RSA generated successful")
}
