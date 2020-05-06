package modules

import (
	"errors"
	"flag"
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"strconv"
	"testing"
)

// executor类型常量
const (
	RCE = "rce"
	SSH = "ssh"
)

var (
	executorType = flag.String(
		"e", "rce",
		"Specify type of executor, rce or ssh",
	)
	host = flag.String(
		"h", "127.0.0.1",
		"Host for executor",
	)
	username = flag.String(
		"u", "",
		"Username for executor, only for SSH",
	)
	password = flag.String(
		"p", "",
		"Password for executor, only for SSH",
	)
	port = flag.String(
		"P", "5051",
		"Port for executor",
	)
)

// NewExecutor 根据参数生成ssh或rce executor
// 进行单元测试的命令（ref：https://stackoverflow.com/a/16353449/7452313）
// go test ./... -v -args -e=rce -h=127.0.0.1 -P=5051
// go test ./... -v -args -e=ssh -u=root -p=root123 -P=22 -h=192.168.0.161
func NewExecutor() (executor.Executor, error) {
	flag.Parse()

	if *executorType == RCE {
		return executor.NewRCEAgentExecutor(*host, *port)
	} else if *executorType == SSH {
		portInt, _ := strconv.Atoi(*port)
		return executor.NewSSHAgentExecutor(*host, *username, *password, portInt)
	} else {
		err := fmt.Sprintf("Unknow executor type: %s", *executorType)
		return nil, errors.New(err)
	}
}

// Execute 执行指定命令
func Execute(t *testing.T, cmdModule, cmdName string, cmdParams executor.ExecutorCmdParams) {
	cmdExecutor, err := NewExecutor()
	if err != nil {
		t.Error("Initialize executor error: ", err)
		return
	}

	defer cmdExecutor.Close()

	cmdFunc, ok := executor.GetCmdByModuleAndName(cmdModule, cmdName)
	if ok {
		er := cmdFunc(cmdExecutor, &cmdParams)
		if er.Successful {
			t.Logf("Command %s executed successful", cmdName)
			t.Logf("ResultData: %s", er.ResultData)
		} else {
			t.Errorf("Command %s executed failed", cmdName)
			t.Errorf("Error: %s", er.Message)
		}
	} else {
		t.Errorf("Can not find command %s", cmdName)
	}
}

var moduleList []string

func init() {
	moduleList = []string{"mysql", "oscmd", "cgroup", "osconfig", "osservice", "zdata"}
}

func AddModule(name string) {
	moduleList = append(moduleList, name)
}

func CheckModuleExsit(name string) bool {
	for _, v := range moduleList {
		if v == name {
			return true
		}
	}
	return false
}
