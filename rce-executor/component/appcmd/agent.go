package appcmd

import (
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/modules/oscmd"
	"hercules_compiler/rce-executor/ssh"
	"path"
)

//定义
const (
	CmdNameInstallAgent = "installAgent"
)

//InstallAgentInTarget  向目标机器安装agent文件
//
//param sourceFilePath 源文件地址 现在支持本地资源，不支持网络资源
//
//param targetPath 目标主机安装目录
//
//param host       目标机器ip地址
//
//param username   目标机器用户名
//
//param password   目标机器密码
//
//return error     错误信息
func InstallAgentInTarget(sourceFilePath, targetPath, host, username, password string, port ...int) error {
	sshPath, err := executor.GetSShPath()
	if err != nil {
		return err
	}
	sshClient := ssh.NewSSHClient(host, username, password,sshPath, port...)
	defer sshClient.Close()

	//@desc 通过sshClient 生成新的sshExecutor
	sshExecutor, err := executor.NewSSHAgentExecutorForSSHClient(sshClient)
	if err != nil {
		return err
	}

	//@desc 获取文件名
	remoteFileName := path.Base(sourceFilePath)
	remoteFullPath := targetPath + "/" + remoteFileName

	//@desc 判断目标机器文件是否存在
	fileExist := false
	log.Debug("check agent file exist or not")
	lsParams := executor.ExecutorCmdParams{oscmd.CmdParamPath: remoteFullPath}
	res, err := executor.ExecCmd(sshExecutor, oscmd.OSCmdModuleName, oscmd.CmdNameList, lsParams)
	if err == nil {
		lsPath, ok := res[oscmd.ResultDataKeyPaths]
		if ok {
			if lsPath == remoteFullPath {
				fileExist = true
			}
		}
	}
	//@desc 目标文件已经存在
	if fileExist {
		log.Debug("agent file already exist")
	} else {
		//@desc 传送文件到
		log.Debug("start upload file to target........")
		err = sshClient.SendFile(sourceFilePath, targetPath)
		if err != nil {
			return err
		}

		//@desc 给文件添加可执行权限
		log.Debug("start chmod for %s", remoteFullPath)
		changeModeParams := executor.ExecutorCmdParams{"modeExp": "u=rwx,g=rwx,o=rwx", "recursiveChange": true}
		changeModeParams["filenamePattern"] = remoteFullPath
		res, err := executor.ExecCmd(sshExecutor, oscmd.OSCmdModuleName, oscmd.CmdNameChangeMode, changeModeParams)
		if err != nil {
			log.Debug("%s execute failed: %v", oscmd.CmdNameChangeMode, err)
			return err
		}
		log.Debug("%s execute success res:: %v", oscmd.CmdNameChangeMode, res)
	}

	//@desc 启动服务
	log.Debug("start agent service ")
	runInBgParams := executor.ExecutorCmdParams{oscmd.CmdParamExecutable: remoteFullPath}
	_, err = executor.ExecCmd(sshExecutor, AppModuleName, CmdNameServiceStart, runInBgParams)
	if err != nil {
		log.Debug("%s execute failed: %v", AppModuleName, err)
		return err
	}
	log.Debug("%s execute success res: %v", AppModuleName, res)
	return nil
}
