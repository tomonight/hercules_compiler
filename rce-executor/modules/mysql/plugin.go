//#管理mysql 插件的安装与卸载
package mysql

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"

	"strings"
)

// 函数名常量定义
const (
	CmdNameInstallPlugin = "InstallPlugin"
)

// 命令参数常量定义
const (
	CmdParamPluginName    = "pluginName"
	CmdParamPluginLibrary = "pluginLibrary"
)

const (
	ErrUDFExists = "ERROR 1125"
)

func InstallPlugin(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//CmdParamMysqlCmdSQL = "cmdSql"
	//plugin_name
	//plugin_library
	pluginName, err := executor.ExtractCmdFuncStringParam(params, CmdParamPluginName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	pluginLibrary, err := executor.ExtractCmdFuncStringParam(params, CmdParamPluginLibrary)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	sql := fmt.Sprintf("INSTALL PLUGIN %s SONAME '%s';", pluginName, pluginLibrary)
	lastParams := *params
	lastParams[CmdParamMysqlCmdSQL] = sql

	er := MySQLCmdSQL(e, &lastParams)
	if !er.Successful {
		if strings.Contains(strings.ToLower(er.Message), strings.ToLower(ErrUDFExists)) {
			er.Successful = true
			er.Message = ""
			return er
		}
	}
	return er
}
