package oscmd

import (
	"hercules_compiler/rce-executor/executor"
	"testing"
)

func TestGrepLine(t *testing.T) {
	ex, err := executor.NewSSHAgentExecutor("192.168.11.167", "root", "123456", 22)
	if err != nil {
		t.Errorf("new executor failed :%s", err.Error())
	}
	params := executor.ExecutorCmdParams{
		CmdParamDirectory: "/root/test.cnf",
		CmdParamValue:     "loose_group_replication_ip_whitelist",
	}
	type args struct {
		e      executor.Executor
		params *executor.ExecutorCmdParams
	}
	tests := []struct {
		name string
		args args
		want executor.ExecuteResult
	}{
		{args: args{
			e:      ex,
			params: &params,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GrepLine(tt.args.e, tt.args.params)
			t.Log(got.ResultData["value"])
		})
	}
}
