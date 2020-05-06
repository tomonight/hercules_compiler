package oscmd

import (
	"encoding/base64"
	"errors"
	"fmt"

	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// 模块名常量定义
const (
	OSCmdModuleName = "oscmd"
)

// 函数名常量定义
const (
	CmdNameAddGroup                = "AddGroup"
	CmdNameAddUser                 = "AddUser"
	CmdNameChangeOwnAndGroup       = "ChangeOwnAndGroup"
	CmdNameMakeDir                 = "MakeDir"
	CmdNameUnzipFile               = "UnzipFile"
	CmdNameDownloadFile            = "DownloadFile"
	CmdNameMD5Sum                  = "MD5Sum"
	CmdNameChangeMode              = "ChangeMode"
	CmdNameMove                    = "Move"
	CmdNameCopy                    = "Copy"
	CmdNameRemove                  = "Remove"
	CmdNameNohup                   = "Nohup"
	CmdNameKillProcessByPID        = "KillProcessByPID"
	CmdNameProcessStatus           = "ProcessStatus"
	CmdNameTextToFile              = "TextToFile"
	CmdNamePing                    = "Ping"
	CmdNameDisplayFileSystem       = "DisplayFileSystem"
	CmdNameDisplayFileSystemByPath = "DisplayFileSystemByPath"
	CmdNameFind                    = "Find"
	CmdNameSourceFile              = "SourceFile"
	CmdNameLink                    = "Link"
	CmdNameRpmInstall              = "RpmInstall"
	CmdNameAddSSHAuthorizedKeys    = "AddSSHAuthorizedKeys"
	CmdNameChangeDir               = "ChangeDir"
	CmdNameYumLocalPackagesInstall = "YumLocalPackagesInstall"
	CmdNameDirUseSpace             = "DirUseSpace"
	CmdNameSysCtl                  = "SysCtl"
	CmdNameGrepLine                = "GrepLine"
	CmdNameSetEnvironmentLang      = "SetEnvironmentLang"
	CmdNameGrepPidOf               = "GrepPidOf"
	CmdNameGetPortByPID            = "GetPortByPID"
	CmdNameGetMysqldPath           = "GetMysqldPath"
	CmdNameGetMysqlServerID        = "GetMysqlServerID"
	CmdAddParamsToConf             = "AddParamsToConf"
	CmdNameGetSlaveIP              = "GetSlaveIP"
	CmdNameLs                      = "Ls"
	CmdNameTest                    = "Test"
	CmdNameCheckMysqlPort          = "CheckMysqlPort"
	CmdNameCheckMysqlPortInConf    = "CheckMysqlPortInConf"
	CmdNameGetAddMGRIp             = "GetAddMGRIp"
	CmdNameGetMGRWhiteList         = "GetMGRWhiteList"
	CmdNameUnInstallSoftware       = "UnInstallSoftware"
	CmdNameSetScriptTaskFailed     = "SetScriptTaskFailed"
)

// 参数常量定义
const (
	CmdParamGroupName       = "groupName"
	CmdParamUserName        = "userName"
	CmdParamPassword        = "password"
	CmdParamFilename        = "filename"
	CmdParamPath            = "path"
	CmdParamDirectory       = "directory"
	CmdParamOwn             = "own"
	CmdParamGroup           = "group"
	CmdParamFilenamePattern = "filenamePattern"
	CmdParamRecursiveChange = "recursiveChange"
	CmdParamModeExp         = "modeExp"
	CmdParamOutputFilename  = "outputFilename"
	CmdParamSource          = "source"
	CmdParamTarget          = "target"
	CmdParamRecursiveCopy   = "recursiveCopy"
	CmdParamRecursiveRemove = "recursiveRemove"
	CmdParamURL             = "url"
	CmdParamMD5             = "md5"
	CmdParamExecutable      = "executable"
	CmdParamLogFile         = "logFile"
	CmdParamPIDFile         = "pidFile"
	CmdParamPID             = "pid"
	CmdParamForceKill       = "forceKill"
	CmdParamProcessName     = "processName"
	CmdParamOverwrite       = "overwrite"
	CmdParamOutText         = "outText"
	CmdParamByBase64        = "byBase64"
	CmdParamHost            = "host"
	CmdParamPrefix          = "prefix"
	CmdParamFileType        = "fileType"
	CmdParamsoftLink        = "softLink"
	CmdParamforceLink       = "forceLink"
	CmdParamGrubFile        = "grubFile"
	CmdParamInstallFlag     = "installFlag"
	CmdParamSSHRsaPubFile   = "SSHRsaPubFile"
	CmdParamCommand         = "command"
	CmdParamValue           = "value"
	CmdParamRemotePath      = "remotePath"
	CmdParamNeedTransfer    = "needTransfer"
	CmdLanguage             = "language"
	CmdParamSSHConnect      = "connect"
	CmdParamTestType        = "type"
	CmdParamAgentType       = "agentType"
	CmdParamRemoteRecvPort  = "remoteRecvPort"
	CmdParamTransferHost    = "transferHost"
	CmdParamMysqlPort       = "port"
)

// 结果集键定义
const (
	ResultDataKeyMD5            = "md5"
	ResultDataKeyUser           = "user"
	ResultDataKeyPID            = "pid"
	ResultDataKeyCPU            = "%cpu"
	ResultDataKeyMem            = "%mem"
	ResultDataKeyStat           = "stat"
	ResultDataKeyStart          = "start"
	ResultDataKeyTime           = "time"
	ResultDataKeyCommand        = "command"
	ResultDataKeyTempDir        = "TempDir"
	ResultDataKeyoutputFilename = "outputFilename"
)

// AddGroup 使用groupadd命令创建组
func AddGroup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	groupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamGroupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	es, err := e.ExecShell("groupadd " + groupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "user group "+groupName+" add successful")
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Index(errMsg, "groupadd:") >= 0 && strings.Index(errMsg, "already exists") >= 0 {
			return executor.SuccessulExecuteResult(es, false, "user group "+groupName+" exists")
		}
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)

	//getent group abcd | cut -d: -f3
}

// AddUser 使用useradd -g命令创建用户
func AddUser(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	userName, err := executor.ExtractCmdFuncStringParam(params, CmdParamUserName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	groupName, err := executor.ExtractCmdFuncStringParam(params, CmdParamGroupName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("useradd %s -M -g %s", userName, groupName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "user "+userName+" add successful")
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
		if strings.Index(errMsg, "useradd:") >= 0 && strings.Index(errMsg, "already exists") >= 0 {
			return executor.SuccessulExecuteResult(es, false, "user "+userName+" exists")
		}
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// MakeDir 使用mkdir命令创建目录
func MakeDir(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	cmdstr := fmt.Sprintf("mkdir -p %s", path)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("MkDir exitcode = %d, path %s, stdErr= %v stdOut= %v", es.ExitCode, path, es.Stderr, es.Stdout)
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "directory "+path+" create successful")
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// UnzipFile 根据目标文件类型使用unzip或tar命令解压文件
func UnzipFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filename, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilename)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	directory, err := executor.ExtractCmdFuncStringParam(params, CmdParamDirectory)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := ""

	if strings.HasSuffix(filename, ".zip") {
		cmdstr = fmt.Sprintf("unzip  %s -d %s", filename, directory)
	} else if strings.HasSuffix(filename, ".tar.gz") {
		cmdstr = fmt.Sprintf("tar -xzf %s -C %s", filename, directory)
	} else if strings.HasSuffix(filename, ".tgz") {
		cmdstr = fmt.Sprintf("tar -xzf %s -C %s", filename, directory)
	} else if strings.HasSuffix(filename, ".tar") {
		cmdstr = fmt.Sprintf("tar -xf %s -C %s", filename, directory)
	} else if strings.HasSuffix(filename, ".xz") {
		cmdstr = fmt.Sprintf("tar -xJf %s -C %s", filename, directory)
	} else {
		return executor.ErrorExecuteResult(fmt.Errorf("can not recognize compressed file"))
	}

	//log.Printf("UnzipFile cmd='%s\n", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "packaged file "+filename+" unflat to "+directory)
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// ChangeOwnAndGroup 使用chown命令修改文件的用户和组
func ChangeOwnAndGroup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	own, err := executor.ExtractCmdFuncStringParam(params, CmdParamOwn)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	group, err := executor.ExtractCmdFuncStringParam(params, CmdParamGroup)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	filenamePattern, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilenamePattern)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	recursiveChange, err := executor.ExtractCmdFuncBoolParam(params, CmdParamRecursiveChange)
	if err != nil {
		recursiveChange = false
	}

	var cmdstr string
	if recursiveChange {
		cmdstr = fmt.Sprintf("chown -R %s:%s %s", own, group, filenamePattern)
	} else {
		cmdstr = fmt.Sprintf("chown %s:%s %s", own, group, filenamePattern)
	}
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		return executor.SuccessulExecuteResult(es, true, "file or directory "+filenamePattern+" changed own:group to "+own+":"+group)
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Find 使用find命令查找文件或目录
func Find(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	directory, err := executor.ExtractCmdFuncStringParam(params, CmdParamDirectory)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	filenamePattern, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilenamePattern)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	fileType, err := executor.ExtractCmdFuncStringParam(params, CmdParamFileType)
	if err != nil {
		fileType = "file"
	}

	fileTypeFlag := ""

	var cmdstr string
	switch fileType {
	case "file":
		fileTypeFlag = "f"
	case "directory":
		fileTypeFlag = "d"
	default:
		return executor.ErrorExecuteResult(fmt.Errorf("Incorrect fileType parameter"))
	}

	cmdstr = fmt.Sprintf(`find  "%s" -type %s -name "%s"`, directory, fileTypeFlag, filenamePattern)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		er := executor.SuccessulExecuteResult(es, false, fmt.Sprintf("find file or directory '%s' in '%s' successful", filenamePattern, directory))
		er.ResultData = make(map[string]string)
		er.ResultData["result"] = strings.Join(es.Stdout, "\n")
		return er
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//DownloadFile 从指定的url下载到outputFilename指向的文件
func DownloadFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	url, err := executor.ExtractCmdFuncStringParam(params, CmdParamURL)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	outputFilename, err := executor.ExtractCmdFuncStringParam(params, CmdParamOutputFilename)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	md5, err := executor.ExtractCmdFuncStringParam(params, CmdParamMD5)
	if md5 != "" {
		log.Debug("input md5 value = %s", md5)
		if isTheSameFile(e, outputFilename, md5) {
			return executor.SuccessulExecuteResultNoData(outputFilename + " already download")
		}
	}

	//log.Printf("DownloadFile url=%s outputFile=%s", url, outputFilename)
	cmdstr := fmt.Sprintf("curl -f -s -o %s --create-dirs %s", outputFilename, url)
	log.Debug("command %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		log.Debug("ExecShell err %s", err)
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 {
		if md5 != "" {
			if !isTheSameFile(e, outputFilename, md5) {
				return executor.ErrorExecuteResult(fmt.Errorf("file %s md5 check failed", outputFilename))
			}
		}

		er := executor.SuccessulExecuteResult(es, true, "url "+url+" downloaded to "+outputFilename+" successful")
		er.ResultData = make(map[string]string)
		er.ResultData[ResultDataKeyoutputFilename] = outputFilename
		return er
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

func isTheSameFile(e executor.Executor, filePath, md5 string) bool {
	params := make(executor.ExecutorCmdParams)
	params[CmdParamFilename] = filePath
	er := MD5Sum(e, &params)
	log.Debug("MD5Sum req = %s", er)
	if er.Successful {
		outMd5 := er.ResultData[ResultDataKeyMD5]
		log.Debug("outMD5 value = %s", outMd5)
		if outMd5 == md5 {
			return true
		}
	}
	return false
}

// MD5Sum 根据不同平台使用md5或md5sum命令计算指定文件的md5值
func MD5Sum(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filename, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilename)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	osType := e.GetExecutorContext("OSType")
	cmdstr := fmt.Sprintf("md5sum %s | cut -d ' ' -f1", filename)
	// MacOS上使用md5命令
	if osType == executor.MacOS {
		cmdstr = fmt.Sprintf("md5 %s | cut -d ' ' -f4", filename)
	}
	log.Debug("MD5 cmdStr %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	log.Debug("MD5 Exec cmdStr %v %v", es, err)

	// if file not exist, es.ExitCode is also 0
	// but has es.Stderr
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "calculate md5 for file "+filename+" successful")
		resultData := map[string]string{ResultDataKeyMD5: es.Stdout[0]}
		er.ResultData = resultData
		return er
	}

	var errMsg string

	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// ChangeMode 使用`chmod`命令修改文件或目录的读写执行权限
func ChangeMode(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	modeExp, err := executor.ExtractCmdFuncStringParam(params, CmdParamModeExp)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	filenamePattern, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilenamePattern)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	recursiveChange, err := executor.ExtractCmdFuncBoolParam(params, CmdParamRecursiveChange)
	if err != nil {
		recursiveChange = false
	}

	var cmdstr string
	if recursiveChange {
		cmdstr = fmt.Sprintf("chmod -R %s %s", modeExp, filenamePattern)
	} else {
		cmdstr = fmt.Sprintf("chmod %s %s", modeExp, filenamePattern)
	}

	log.Debug("cmdStr = %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "file or directory "+filenamePattern+" changed mode to "+modeExp)
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Move 使用`mv`命令移动或重命名文件或目录
func Move(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	source, err := executor.ExtractCmdFuncStringParam(params, CmdParamSource)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	target, err := executor.ExtractCmdFuncStringParam(params, CmdParamTarget)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("mv %s %s", source, target)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "file or directory "+source+" moved to "+target)
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// ChangeDir 使用`cd`命令切换目录
func ChangeDir(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	targetPath, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("cd %s", targetPath)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, " change path to "+targetPath)
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Copy 使用`cp`命令复制文件或目录
func Copy(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	source, err := executor.ExtractCmdFuncStringParam(params, CmdParamSource)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	target, err := executor.ExtractCmdFuncStringParam(params, CmdParamTarget)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	recursiveCopy, err := executor.ExtractCmdFuncBoolParam(params, CmdParamRecursiveCopy)
	if err != nil {
		recursiveCopy = false
	}
	overwrite, err := executor.ExtractCmdFuncBoolParam(params, CmdParamOverwrite)
	if err != nil {
		overwrite = true
	}

	cpcmd := "cp"
	if overwrite {
		cpcmd = "\\cp -f"
	}

	var cmdstr string
	if recursiveCopy {
		cmdstr = fmt.Sprintf("%s -R %s %s", cpcmd, source, target)
	} else {
		cmdstr = fmt.Sprintf("%s %s %s", cpcmd, source, target)
	}
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "file or directory "+source+" copied to "+target)
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Remove 使用`rm`命令移除文件或目录
func Remove(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filenamePattern, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilenamePattern)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if filenamePattern == "" {
		return executor.SuccessulExecuteResultNoData("file or directory " + filenamePattern + " deleted")
	}
	recursiveRemove, err := executor.ExtractCmdFuncBoolParam(params, CmdParamRecursiveRemove)
	if err != nil {
		recursiveRemove = false
	}

	var cmdstr string
	if recursiveRemove {
		cmdstr = fmt.Sprintf("rm -rf %s", filenamePattern)
	} else {
		cmdstr = fmt.Sprintf("rm -f %s", filenamePattern)
	}
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "file or directory "+filenamePattern+" deleted")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Nohup 使用`nohup`命令在后台启动服务
func Nohup(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	executable, err := executor.ExtractCmdFuncStringParam(params, CmdParamExecutable)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	logFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamLogFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	pidFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamPIDFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("nohup %s > %s 2>&1 & echo $! > %s && echo $!", executable, logFile, pidFile)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, executable+" running in backgroud successful")
		er.ResultData = map[string]string{
			ResultDataKeyPID: es.Stdout[0],
		}
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// KillProcessByPID 使用`kill`命令停止进程
func KillProcessByPID(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	pid, err := executor.ExtractCmdFuncStringParam(params, CmdParamPID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	forceKill, err := executor.ExtractCmdFuncBoolParam(params, CmdParamForceKill)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	var cmdstr string
	if forceKill {
		cmdstr = fmt.Sprintf("kill -9 %s", pid)
	} else {
		cmdstr = fmt.Sprintf("kill %s", pid)
	}
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "process id "+pid+" killed")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.ToLower(es.Stderr[0])
		if strings.Index(errMsg, "kill:") >= 0 && strings.Index(errMsg, "no such process") >= 0 {
			return executor.SuccessulExecuteResult(es, false, "process id "+pid+"not exists")
		}
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// ProcessStatus 使用`ps -ef`命令查看进程状态
func ProcessStatus(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	processName, err := executor.ExtractCmdFuncStringParam(params, CmdParamProcessName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf(`ps aux | egrep "%s" | grep -v grep`, processName)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		// ref: https://stackoverflow.com/a/13737890/7452313
		stdout := strings.Fields(es.Stdout[0])
		er := executor.SuccessulExecuteResult(es, true, "process status for "+processName+" get successful")
		er.ResultData = map[string]string{
			ResultDataKeyUser:    stdout[0],
			ResultDataKeyPID:     stdout[1],
			ResultDataKeyCPU:     stdout[2],
			ResultDataKeyMem:     stdout[3],
			ResultDataKeyStat:    stdout[7],
			ResultDataKeyStart:   stdout[8],
			ResultDataKeyTime:    stdout[9],
			ResultDataKeyCommand: strings.Join(stdout[10:], " "),
		}
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		if es.ExitCode == 1 {
			errMsg = executor.ErrProcessNotFound
		} else {
			errMsg = executor.ErrMsgUnknow
		}
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// TextToFile 将一个文本输出到文件中
// overwrite 可设置覆盖文件内容或追加 Todo追加需要判断追加字符串是否已经存在
func TextToFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	filename, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilename)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	outText, err := executor.ExtractCmdFuncStringParam(params, CmdParamOutText)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	overwrite, err := executor.ExtractCmdFuncBoolParam(params, CmdParamOverwrite)
	if err != nil {
		// 默认为覆盖模式
		overwrite = true
	}

	byBase64, err := executor.ExtractCmdFuncBoolParam(params, CmdParamByBase64)
	if err != nil {
		// 默认Base64为false
		byBase64 = false
	}

	var writeMode string
	if overwrite {
		writeMode = ">"
	} else { //追加
		//追加字符串判断是否已经存在
		cmdStr := fmt.Sprintf("cat %s", filename)
		es, err := e.ExecShell(cmdStr)
		if err != nil {
			log.Warn(cmdStr, "executor failed ", err)
			return executor.NotSuccessulExecuteResult(es, err.Error())
		}
		err = executor.GetExecResult(es)
		if err != nil {
			log.Warn(cmdStr, "executor failed ", err)
			return executor.NotSuccessulExecuteResult(es, err.Error())
		}

		for _, text := range es.Stdout {
			//			if strings.Contains(text, outText) {
			//				if !strings.Contains(text, "#") && !strings.Contains(outText, "#") {
			//					return executor.SuccessulExecuteResult(es, true, "text saved to file "+filename)
			//				}
			//			}
			if text == outText {
				return executor.SuccessulExecuteResult(es, true, "text saved to file "+filename)
			}
		}

		writeMode = ">>"
	}

	cmdstr := ""
	if byBase64 {
		encoded := base64.StdEncoding.EncodeToString([]byte(outText))
		cmdstr = fmt.Sprintf(`echo -e "%s" | base64 -d - %s %s`, utils.EscapeShellCmd(encoded), writeMode, filename)

	} else {
		cmdstr = fmt.Sprintf(`echo -e "%s" %s %s`, utils.EscapeShellCmd(outText), writeMode, filename)

	}
	log.Debug("cmdstr= %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	fmt.Printf(strings.Join(es.Stderr, "\n"))
	fmt.Printf(strings.Join(es.Stdout, "\n"))

	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "text saved to file "+filename)
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Ping 使用ping命令检查主机连接性，发包5次，超时时间为5秒
// success表示可以ping通
func Ping(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	host, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	osType := e.GetExecutorContext("OSType")
	cmdstr := fmt.Sprintf("ping -w 5 -c 5 %s", host)
	if osType == executor.MacOS {
		cmdstr = fmt.Sprintf("ping -t 5 -c 5 %s", host)
	}
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "ping successful")
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// DisplayFileSystem 使用df -klP命令查看文件系统使用情况
// 结果使用|进行分割，字段顺序为：
// Filesystem|1K-blocks|Used|Available|Capacity|Mounted on
// 数据单位均为KB
func DisplayFileSystem(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmdstr := fmt.Sprintf("df -hlP")
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "display file system successful")
		er.ResultData = make(map[string]string)
		for index, data := range es.Stdout[1:] {
			dataList := strings.Fields(data)
			dataString := strings.Join(dataList, "|")
			er.ResultData[strconv.Itoa(index)] = dataString
		}
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// DisplayFileSystemByPath 使用df -klP命令查看文件系统使用情况
// 结果使用|进行分割，字段顺序为：
// Filesystem|1K-blocks|Used|Available|Capacity|Mounted on
// 数据单位均为KB
func DisplayFileSystemByPath(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		path = e.GetExecutorContext(executor.ContextNameTempDir)
	}
	cmdstr := fmt.Sprintf("df -kP %s", path)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "display file system successful")
		er.ResultData = make(map[string]string)
		dataList := strings.Fields(es.Stdout[1])
		er.ResultData["total"] = dataList[1]
		er.ResultData["used"] = dataList[2]
		er.ResultData["available"] = dataList[3]
		er.ResultData["mountOn"] = dataList[5]
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

func nextRandom() string {
	r := uint32(time.Now().UnixNano())
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// CreateTempDir 在指定的目录下创建一个临时目录，目录名是随机的
//@param path 在path指定路径下创建临时目录
//@out param 输出TempDir 获得创建的路径名
func CreateTempDir(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		path = e.GetExecutorContext(executor.ContextNameTempDir)
	}
	prefix, _ := executor.ExtractCmdFuncStringParam(params, CmdParamPrefix)
	tempName := path + e.GetExecutorContext(executor.ContextNamePathSeparator) + prefix + nextRandom()
	cmdParams := executor.ExecutorCmdParams{CmdParamPath: tempName}
	er := MakeDir(e, &cmdParams)
	if er.Successful {
		er.ResultData = make(map[string]string)
		er.ResultData[ResultDataKeyTempDir] = tempName
		er.Message = "temporary directory " + tempName + " created"
	}
	return er
}

// SourceFile 封装source命令
func SourceFile(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filename, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilename)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("source %s", filename)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "source file successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// SysCtl 封装sysctl命令
func SysCtl(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	command, err := executor.ExtractCmdFuncStringParam(params, CmdParamCommand)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	value, err := executor.ExtractCmdFuncStringParam(params, CmdParamValue)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf("sysctl -w %s=%s", command, value)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "SysCtl excute successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// Link 使用ln命令生成连接
func Link(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	source, err := executor.ExtractCmdFuncStringParam(params, CmdParamSource)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	target, err := executor.ExtractCmdFuncStringParam(params, CmdParamTarget)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	forceLink, err := executor.ExtractCmdFuncBoolParam(params, CmdParamforceLink)
	if err != nil {
		// default to true
		forceLink = true
	}

	softLink, err := executor.ExtractCmdFuncBoolParam(params, CmdParamsoftLink)
	if err != nil {
		// default to true
		softLink = true
	}

	var forceLinkFlag, softLinkFlag string = "", ""
	if forceLink {
		forceLinkFlag = " -f"
	}
	if softLink {
		softLinkFlag = " -s"
	}
	cmdstr := fmt.Sprintf("ln%s%s %s %s", forceLinkFlag, softLinkFlag, source, target)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "link creation successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

func getRpmInstallError(errMsg string) error {
	flag := []string{"already installed]"}
	for _, value := range flag {
		if strings.Contains(errMsg, value) {
			return nil
		}
	}
	return errors.New(errMsg)
}

// RpmInstall 使用rpm命令安装软件包
func RpmInstall(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	filename, err := executor.ExtractCmdFuncStringParam(params, CmdParamFilename)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	installFlags, err := executor.ExtractCmdFuncStringParam(params, CmdParamInstallFlag)
	if err != nil {
		installFlags = " -ivh"
	}

	cmdstr := fmt.Sprintf("rpm%s %s", installFlags, filename)
	log.Debug("cmdStr %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	log.Debug("res:", es)
	log.Debug("exitCode:%d err %v", es.ExitCode, err)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 || es.ExitCode == 1 /*&& len(es.Stderr) == 0*/ {

		return executor.SuccessulExecuteResult(es, true, "rpm installed successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	err = getRpmInstallError(strings.Join(es.Stderr, ","))
	if err != nil {
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}

	//return executor.NotSuccessulExecuteResult(es, errMsg)
	return executor.SuccessulExecuteResult(es, false, "rpm installed successful")
}

// AddSSHAuthorizedKeys 添加 authorized_keys 文件，配置ssh互信
func AddSSHAuthorizedKeys(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.GetSShPath()
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	connect, err := executor.ExtractCmdFuncBoolParam(params, CmdParamSSHConnect)
	if err != nil {
		connect = true
	}
	rsaPubFile, err := executor.ExtractCmdFuncStringParam(params, CmdParamSSHRsaPubFile)
	if err != nil {
		rsaPubFile = os.ExpandEnv(path + "/data/.ssh/id_rsa.pub")
		rsaprikeyFile := os.ExpandEnv(path + "/data/.ssh/id_rsa")
		if !utils.PathExist(rsaPubFile) || !utils.PathExist(rsaprikeyFile) || connect == false {
			rsaPrivDir := os.ExpandEnv(path + "/data/.ssh")
			sshKeyGenCmd := fmt.Sprintf("rm -f %s/id_rsa && mkdir -p %s && ssh-keygen -f %s/id_rsa -P ''", rsaPrivDir, rsaPrivDir, rsaPrivDir)
			localExecutor, err := executor.NewLocalExecutor()
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
			_, err = localExecutor.ExecShell(sshKeyGenCmd)
			if err != nil {
				return executor.ErrorExecuteResult(err)
			}
		}
	}
	file, err := os.Open(rsaPubFile)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	bRes, err := ioutil.ReadAll(file)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	body := string(bRes)
	body = strings.TrimSuffix(body, "\n")

	cmdstr := fmt.Sprintf(
		`mkdir -p ~/.ssh; grep -q "%s" ~/.ssh/authorized_keys 2>/dev/null|| echo "%s" >> ~/.ssh/authorized_keys; chmod 644 ~/.ssh/authorized_keys; chmod 700 ~/.ssh/;restorecon -R -v /root/.ssh`,
		body, body,
	)
	log.Debug("cmdStr %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, "SSH authorized_keys added successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//SetEnvironmentLang set env LANG=en_US 设置环境语言
func SetEnvironmentLang(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	//command = "set env LANG=en_US"
	lang, err := executor.ExtractCmdFuncStringParam(params, CmdLanguage)
	if err != nil {
		lang = "en_US"
	}

	cmdstr := fmt.Sprintf("set env LANG=%s", lang)
	log.Debug("cmdStr %s", cmdstr)
	es, err := e.ExecShell(cmdstr)
	log.Debug("res:", es)
	log.Debug("exitCode:%d err %v", es.ExitCode, err)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		return executor.SuccessulExecuteResult(es, true, cmdstr+" successful")
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

////AddReadOnlyToConf 添加readonly到conf文件
//func AddReadOnlyToConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
//	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
//	if err != nil {
//		return executor.ErrorExecuteResult(err)
//	}
//	var cmd = `sed -i -r '/^read_only|[[:space:]]+read_only|super_read_only/d' %s && sed -i '/\[mysqld\]/a read_only=1\nsuper_read_only=1' %s`
//	cmdstr := fmt.Sprintf(cmd, path, path)
//	log.Debug("cmdStr %s", cmdstr)
//	es, err := e.ExecShell(cmdstr)
//	log.Debug("res:", es)
//	log.Debug("exitCode:%d err %v", es.ExitCode, err)
//	if err != nil {
//		return executor.ErrorExecuteResult(err)
//	}
//
//	if es.ExitCode == 0 && len(es.Stderr) == 0 {
//		return executor.SuccessulExecuteResult(es, true, cmdstr+" successful")
//	}
//
//	var errMsg string
//	if len(es.Stderr) == 0 {
//		errMsg = executor.ErrMsgUnknow
//	} else {
//		errMsg = es.Stderr[0]
//	}
//	return executor.NotSuccessulExecuteResult(es, errMsg)
//}

//AddParamsToConf 添加params到conf文件
func AddParamsToConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	cmds, err := executor.ExtractCmdFuncStringListParam(params, CmdParamCommand, "*#*")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	var errMsg string
	for _, cmd := range cmds {
		log.Debug("cmdStr %s", cmd)
		es, err := e.ExecShell(cmd)
		if err != nil {
			return executor.ErrorExecuteResult(err)
		}
		log.Debug("res:%v", es)
		log.Debug("res std out %v ", es.Stdout)
		log.Debug("res std error %v ", es.Stderr)
		log.Debug("exitCode:%d err %v", es.ExitCode, err)
		if es.ExitCode == 0 {
			er := executor.SuccessulExecuteResult(es, true, fmt.Sprintf("command【%s】 execute successful", cmd))
			er.ResultData = make(map[string]string)
			er.ResultData["test"] = "successsssss"
			return er
			// continue
		}

		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}

	return executor.SuccessulExecuteResultNoData(errMsg)
}

//查看是否存在安装目录
func Ls(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	var cmdStr = "ls " + path
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "dir exists")
		er.ResultData = make(map[string]string)
		files := []string{}
		for _, v := range es.Stdout {
			files = append(files, v)
		}
		er.ResultData["files"] = strings.Join(files, ",")
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = es.Stderr[0]
	}
	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//Test 判断文件是否存在
func Test(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	mtype, err := executor.ExtractCmdFuncStringParam(params, CmdParamTestType)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	var cmdStr = fmt.Sprintf("test -%s %s", mtype, path)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	er := executor.SuccessulExecuteResult(es, true, "test file successful")
	er.ResultData = make(map[string]string)
	er.ResultData["code"] = fmt.Sprintf("%d", es.ExitCode)
	return er
}

// CheckMysqlPort 检查端口是否是当前mysql端口
func CheckMysqlPort(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	pid, err := executor.ExtractCmdFuncIntParam(params, CmdParamPID)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamMysqlPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf(`ps -ef | grep mysqld | grep %d | grep %d`, pid, port)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "CheckMysqlPort successful ,result:"+es.Stdout[0])
		for _, resultStr := range es.Stdout {
			if strings.Contains(resultStr, cmdstr) {
				continue
			}
			er.ResultData = map[string]string{
				"value": "true",
			}
			return er
		}
		er.ResultData = map[string]string{
			"value": "false",
		}
		return er
	}

	errMsg := strings.Join(es.Stderr, "\n")

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// CheckMysqlPortInConf 检查端口是否是当前mysql端口，从配置文件查找
func CheckMysqlPortInConf(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	port, err := executor.ExtractCmdFuncIntParam(params, CmdParamMysqlPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	cmdstr := fmt.Sprintf(`grep -i 'port' %s|grep %d`, path, port)
	es, err := e.ExecShell(cmdstr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode == 0 && len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "CheckMysqlPort successful ,result:"+es.Stdout[0])
		if len(es.Stdout) > 0 {
			er.ResultData = map[string]string{
				"value": "true",
			}
		} else {
			er.ResultData = map[string]string{
				"value": "false",
			}
		}
		return er
	}

	var errMsg string
	if len(es.Stderr) == 0 {
		errMsg = executor.ErrMsgUnknow
	} else {
		errMsg = strings.Join(es.Stderr, "\n")
	}

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

// GetAddMGRIp GetAddMGRIp
func GetAddMGRIp(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	sourceStr, err := executor.ExtractCmdFuncStringParam(params, "sourceStr")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	addStr, err := executor.ExtractCmdFuncStringParam(params, "addStr")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	returnStr := sourceStr

	if !strings.Contains(sourceStr, addStr) {
		returnStr = returnStr + "," + addStr
	}

	es, err := e.ExecShell("hostname")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "GetAddMGRIp successful ,result:"+es.Stdout[0])
		er.ResultData = map[string]string{
			"value": returnStr,
		}
		return er
	}

	errMsg := strings.Join(es.Stderr, "\n")

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

//GetMGRWhiteList GetMGRWhiteList
func GetMGRWhiteList(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	whiteListStr, err := executor.ExtractCmdFuncStringParam(params, "whiteListStr")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	returnStr := getMGRWhiteListFunc(whiteListStr)
	es, err := e.ExecShell("hostname")
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if len(es.Stderr) == 0 {
		er := executor.SuccessulExecuteResult(es, true, "GetAddMGRIp successful ,result:"+es.Stdout[0])
		er.ResultData = map[string]string{
			"value": returnStr,
		}
		return er
	}

	errMsg := strings.Join(es.Stderr, "\n")

	return executor.NotSuccessulExecuteResult(es, errMsg)
}

func getMGRWhiteListFunc(whiteListStr string) string {

	whiteList := strings.Split(whiteListStr, ",")
	newWhiteList := make([]string, 0)
	aSlice := make([][]string, 0)
	bSlice := make([][]string, 0)
	cSlice := make([][]string, 0)

	for _, while := range whiteList {
		if while == "" {
			continue
		}
		split := strings.Split(while, ".")
		num, err := strconv.Atoi(split[0])
		if err != nil {
			log.Debug("Atoi err", err)
			continue
		}
		if num > 0 && num < 128 {
			split[1] = "0"
			split[2] = "0"
			split[3] = "0/8"
			aSlice = append(aSlice, split)
		}

		if num > 127 && num < 192 {
			split[2] = "0"
			split[3] = "0/16"
			bSlice = append(bSlice, split)
		}

		if num > 191 && num < 233 {
			split[3] = "0/24"
			cSlice = append(cSlice, split)
		}
	}
	if len(aSlice) != 0 {
		newWhiteList = append(newWhiteList, getIPListByType(aSlice, 0)...)
	}
	if len(bSlice) != 0 {
		newWhiteList = append(newWhiteList, getIPListByType(bSlice, 1)...)
	}
	if len(cSlice) != 0 {
		newWhiteList = append(newWhiteList, getIPListByType(cSlice, 2)...)
	}

	return strings.Join(newWhiteList, ",")
}

func getIPListByType(aSlice [][]string, index int) []string {
	var newWhiteList []string
	if len(aSlice) == 1 {
		ip := strings.Join(aSlice[0], ".")
		newWhiteList = append(newWhiteList, ip)
		return newWhiteList
	}

	aMap := make(map[string]int)
	for i := 0; i < len(aSlice); i++ {
		var targetStr string
		for i, v := range aSlice[i] {
			if i == (index + 1) {
				break
			}
			targetStr += v
		}
		aMap[targetStr]++

	}
	if len(aMap) == len(aSlice) {
		for i := 0; i < len(aSlice); i++ {
			ip := strings.Join(aSlice[i], ".")
			newWhiteList = append(newWhiteList, ip)
		}
	} else if len(aSlice) > len(aMap) {
		if len(aMap) == 1 {
			ip := strings.Join(aSlice[0], ".")
			newWhiteList = append(newWhiteList, ip)
		} else {
			for e := range aMap {
				if aMap[e] > 1 {
					for i := 0; i < len(aSlice); i++ {
						var targetSir string
						for i, v := range aSlice[i] {
							if i == (index + 1) {
								break
							}
							targetSir += v
						}
						if targetSir == e {
							ip := strings.Join(aSlice[i], ".")
							newWhiteList = append(newWhiteList, ip)
							break
						}
					}
					for i := 0; i < len(aSlice); i++ {
						var targetSir string
						for i, v := range aSlice[i] {
							if i == (index + 1) {
								break
							}
							targetSir += v
						}
						if targetSir != e {
							ip := strings.Join(aSlice[i], ".")
							newWhiteList = append(newWhiteList, ip)
						}
					}
				}
			}
		}

	}
	return newWhiteList
}

//UnInstallSoftware 卸载平台安装的软件并删除配置文件
func UnInstallSoftware(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	processName, err := executor.ExtractCmdFuncStringParam(params, CmdParamProcessName)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	version, err := executor.ExtractCmdFuncStringParam(params, CmdParamOSVersion)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	es, err := e.ExecShell("rm -rf " + path)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	var errMsg string
	if es.ExitCode != 0 || len(es.Stderr) != 0 {
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}

	rmCmd := "rm -rf " + "/etc/systemd/system/" + strings.ToLower(processName) + ".service"
	if version == "6" {
		rmCmd = "rm -rf" + "/etc/init.d/" + strings.ToLower(processName)
	}
	es, err = e.ExecShell(rmCmd)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	if es.ExitCode != 0 || len(es.Stderr) != 0 {
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = es.Stderr[0]
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}

	return executor.SuccessulExecuteResultNoData(errMsg)
}

//SetScriptFailed 设置脚本任务执行失败
func SetScriptTaskFailed(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {
	return executor.ErrorExecuteResult(fmt.Errorf("脚本任务执行失败"))
}
