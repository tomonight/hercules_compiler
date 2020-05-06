package mysql

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"path/filepath"
)

var canDoSourceCommandList = []string{"mariabackup", "mbstream"}

func ErrCommandNotSupport(command string) error {
	return fmt.Errorf("not support this source command %s", command)
}

func commandCanExecute(command string) (canDo bool) {
	for _, canDoCommand := range canDoSourceCommandList {
		if canDoCommand == command {
			canDo = true
			return
		}
	}
	return
}

//getCommandAndOptionsParams  get command and options from params
func getCommandAndOptionsParams(params *executor.ExecutorCmdParams) (command, options string) {
	command, _ = executor.ExtractCmdFuncStringParam(params, "command")
	options, _ = executor.ExtractCmdFuncStringParam(params, "options")
	log.Info("getCommandAndOptionsParams output command = %s options = %s", command, options)
	return
}

//getBackupOrRecoverExecuteCommand get execute command
func getBackupOrRecoverExecuteCommand(backup bool, params *executor.ExecutorCmdParams) (outCommand string, err error) {
	command, options := getCommandAndOptionsParams(params)
	if command != "" {
		if !commandCanExecute(command) {
			err = ErrCommandNotSupport(command)
			return
		}
		outCommand = command
		if options != "" {
			outCommand += " " + options
		}
		return
	}
	return getXtrabackupCommand(backup, params)
}

func getXtrabackupCommand(backup bool, params *executor.ExecutorCmdParams) (xtrabackupCmd string, err error) {
	var xtrabackupPath string
	xtrabackupPath, err = executor.ExtractCmdFuncStringParam(params, CmdParamXtrabackupPath)
	if err != nil {
		return
	}

	if backup {
		xtrabackupCmd = filepath.Join(xtrabackupPath, "xtrabackup")
	} else {
		xtrabackupCmd = filepath.Join(xtrabackupPath, "xbstream")
	}
	return
}
