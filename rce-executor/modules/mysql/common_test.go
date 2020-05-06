package mysql

import (
	"hercules_compiler/rce-executor/executor"
	"testing"
)

func Test_getBackupOrRecoverExecuteCommand(t *testing.T) {
	type args struct {
		backup bool
		params *executor.ExecutorCmdParams
	}
	tests := []struct {
		name           string
		args           args
		wantOutCommand string
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "when other command not set and back up is true then will get mysql default backup command",
			args: args{
				backup: true,
				params: &executor.ExecutorCmdParams{CmdParamXtrabackupPath: "/usr/local/bin"},
			},
			wantOutCommand: "/usr/local/bin/xtrabackup",
			wantErr:        false,
		},

		{
			name: "when other command not set and default command is nil will error result",
			args: args{
				backup: true,
				params: &executor.ExecutorCmdParams{},
			},
			wantOutCommand: "",
			wantErr:        true,
		},


		{
			name: "when other command not set and back up is false then will get mysql default recover command",
			args: args{
				backup: false,
				params: &executor.ExecutorCmdParams{CmdParamXtrabackupPath: "/usr/local/bin"},
			},
			wantOutCommand: "/usr/local/bin/xbstream",
			wantErr:        false,
		},

		{
			name: "when other command is set as mariabackup  then will get mariabackup backup command",
			args: args{
				backup: false,
				params: &executor.ExecutorCmdParams{"command": "mariabackup"},
			},
			wantOutCommand: "mariabackup",
			wantErr:        false,
		},

		{
			name: "when other command is set as mariabackup and options as tbs  then will get mariabackup backup command",
			args: args{
				backup: false,
				params: &executor.ExecutorCmdParams{"command": "mariabackup", "options":"tbs"},
			},
			wantOutCommand: "mariabackup tbs",
			wantErr:        false,
		},

		{
			name: "when other command is set as mbstream  then will get mbstream backup command",
			args: args{
				backup: false,
				params: &executor.ExecutorCmdParams{"command": "mbstream"},
			},
			wantOutCommand: "mbstream",
			wantErr:        false,
		},

		{
			name: "when other command is set as rm then will get error result",
			args: args{
				backup: false,
				params: &executor.ExecutorCmdParams{"command": "rm"},
			},
			wantOutCommand: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutCommand, err := getBackupOrRecoverExecuteCommand(tt.args.backup, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBackupOrRecoverExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutCommand != tt.wantOutCommand {
				t.Errorf("getBackupOrRecoverExecuteCommand() gotOutCommand = %v, want %v", gotOutCommand, tt.wantOutCommand)
			}
		})
	}
}
