package oscmd

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/log"
	"hercules_compiler/rce-executor/ssh"
	"path"
	"strings"
	"time"
)

//define agent type
const (
	AgentTypeForZC = "ZCAgent"
)

//SocatSendFolder 发送
func SocatSendFolder(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	remoteHost, err := executor.ExtractCmdFuncStringParam(params, CmdParamHost)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	remotePort, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	//tar -cvf - 20190122105719 | socat -b 10485760 -u stdio TCP:192.168.11.202:9980
	cmdStr := fmt.Sprintf("tar -cvf - %s | socat -b 10485760 -u stdio TCP:%s:%d",
		path, remoteHost, remotePort)

	log.Debug("execute command %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.SuccessulExecuteResult(es, true, "socat send folder successful ")
}

//SocatRecvFolder 接收
func SocatRecvFolder(e executor.Executor, params *executor.ExecutorCmdParams) executor.ExecuteResult {

	remotePort, err := executor.ExtractCmdFuncIntParam(params, CmdParamPort)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	path, err := executor.ExtractCmdFuncStringParam(params, CmdParamPath)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}

	//socat -b 10485760 -u TCP-LISTEN:9980,reuseaddr stdio |tar -C /home/dw -xvf -
	cmdStr := fmt.Sprintf("socat -b 10485760 -u TCP-LISTEN:%d,reuseaddr stdio |tar -C %s -xvf -",
		remotePort, path)

	log.Debug("execute command %s", cmdStr)
	es, err := e.ExecShell(cmdStr)
	if err != nil {
		return executor.ErrorExecuteResult(err)
	}
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return executor.NotSuccessulExecuteResult(es, errMsg)
	}
	return executor.SuccessulExecuteResult(es, true, "socat recv folder successful ")
}

//getValidPort 获取适合的端口
func getValidPort(e executor.Executor, port uint) (uint, error) {
	for index := 0; index < 100; index++ {
		params := make(executor.ExecutorCmdParams)
		currentPort := port + uint(index)
		params[CmdParamPort] = fmt.Sprintf("%d", currentPort)
		result := PortValid(e, &params)
		if result.Successful {
			log.Debug("PortValid failed =%s", result.Message)
			return currentPort, nil
		}
	}
	return 0, fmt.Errorf("未能获取到适合的传输端口")
}

//exitSocatProcess 退出socat进程
func exitSocatProcess(e executor.Executor, port uint) error {
	netstatComand := fmt.Sprintf("netstat -tunlp|grep %d", port)
	log.Debug("netstatComand = %s", netstatComand)
	es, err := e.ExecShell(netstatComand)
	if err != nil {
		log.Debug("err = %v", err)
		return err
	}
	log.Debug("es.Stdout %v and len = %d", es.Stdout, len(es.Stdout))
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			log.Debug("no errmsg")
			return nil
		}
		errMsg = strings.Join(es.Stderr, "\n")
		log.Debug("err = %v", errMsg)
		return fmt.Errorf(errMsg)
	}

	if len(es.Stdout) > 0 {
		resText := es.Stdout[0]
		log.Debug("resText = %s", resText)
		resTextSplit := strings.Split(resText, " ")
		for _, text := range resTextSplit {
			if strings.Contains(text, "socat") {
				log.Debug("text = %s", text)
				processText := text
				processSlice := strings.Split(processText, "/")
				if len(processSlice) == 2 {
					processID := processSlice[0]
					log.Debug("processID = %s", processID)
					killCommand := fmt.Sprintf("kill -9 %s", processID)
					es, err := e.ExecShell(killCommand)
					if err != nil {
						return err
					}
					if es.ExitCode != 0 {
						var errMsg string
						if len(es.Stderr) == 0 {
							errMsg = executor.ErrMsgUnknow
						} else {
							errMsg = strings.Join(es.Stderr, "\n")
						}
						return fmt.Errorf(errMsg)
					}
					break
				}
			}
		}
	}
	return nil
}

//socatTransferFiles
func socatTransferFiles(e executor.Executor, agentType, destHost, transferHost, destUserName, sourcePath, destPath string, destPort, socatPort uint) error {
	var destExecutor executor.Executor
	var err error

	sourceDir, lastPath := path.Split(sourcePath)
	if sourceDir == "" || lastPath == ""{
		return fmt.Errorf("传输目录:%s 格式不合法", sourcePath)
	}

	if destPath == ""{
		return fmt.Errorf("目标文件路径:%s 不合法", sourcePath)
	}

	if agentType == AgentTypeForZC {
		//create dest executor
		destExecutor, err = executor.NewZCAgentExecutor(destHost, int(destPort))
		if err != nil {
			return err
		}
	} else {
		//create dest executor
		path, err := executor.GetSShPath()
		if err != nil {
			return err
		}
		destClient := ssh.NewSSHClient(destHost, destUserName, "", path, int(destPort))
		destExecutor, err = executor.NewSSHAgentExecutorForSSHClient(destClient)
		if err != nil {
			return err
		}
		defer destClient.Close()
	}

	defer exitSocatProcess(destExecutor, socatPort)

	log.Debug("destHost=%s", destHost)
	if agentType != AgentTypeForZC {
		socatPort, err = getValidPort(destExecutor, socatPort)
		if err != nil {
			return err
		}
	}

	//create target folder
	mkdirCommand := fmt.Sprintf("mkdir -p %s", destPath)
	es, err := destExecutor.ExecShell(mkdirCommand)
	if err != nil {
		return err
	}
	if es.ExitCode != 0 {
		var errMsg string
		if len(es.Stderr) == 0 {
			errMsg = executor.ErrMsgUnknow
		} else {
			errMsg = strings.Join(es.Stderr, "\n")
		}
		return fmt.Errorf(errMsg)
	}

	doNext := make(chan int)
	transferResult := make(chan error, 2)

	go func() {
		var resError error
		destCommand := fmt.Sprintf("socat -b 10485760 -u TCP-LISTEN:%d,reuseaddr stdio |tar -C %s -xvf -",
			socatPort, destPath)
		log.Debug("execute command %s", destCommand)
		doNext <- 1
		es, err := destExecutor.ExecShell(destCommand)
		if err != nil {
			resError = err
		}
		if es.ExitCode != 0 {
			var errMsg string
			if len(es.Stderr) == 0 {
				errMsg = executor.ErrMsgUnknow
			} else {
				errMsg = strings.Join(es.Stderr, "\n")
			}
			resError = fmt.Errorf(errMsg)
		}
		transferResult <- resError
	}()



	go func() {
		<-doNext
		time.Sleep(time.Second * 10)
		var resError error
		sourceCommand := fmt.Sprintf("tar -C %s  -cvf - %s | socat -b 10485760 -u stdio TCP:%s:%d",
			sourceDir, lastPath, transferHost, socatPort)
		log.Debug("execute command %s", sourceCommand)
		es, err := e.ExecShell(sourceCommand)
		if err != nil {
			resError = err
		}
		if es.ExitCode != 0 {
			var errMsg string
			if len(es.Stderr) == 0 {
				errMsg = executor.ErrMsgUnknow
			} else {
				errMsg = strings.Join(es.Stderr, "\n")
			}
			resError = fmt.Errorf(errMsg)
		}
		transferResult <- resError
	}()

	doneResult := make(chan error)
	go func() {
		var socatError error
		defer func() {
			doneResult <- socatError
		}()

		for index := 0; index < 2; index++ {
			err := <-transferResult
			if err != nil {
				socatError = err
				return
			}
		}
	}()

	var lastErr error
	select {
	case err := <-doneResult:
		lastErr = err
	case <-time.After(time.Hour * 2):
		lastErr = fmt.Errorf("socatTansferFiles timeout")
	}

	if lastErr != nil {
		return lastErr
	}

	////mv temp folder to dest folder
	//sourcePath = filepath.Join(destPath, sourcePath)
	//log.Debug("target sourcePath=%s", sourcePath)
	//fileName := utils.GetFileName(sourcePath)
	//destPath = filepath.Join(destPath, fileName)
	//log.Debug("target dest path=%s", destPath)
	//mvCommand := fmt.Sprintf("mv -fT %s %s", sourcePath, destPath)
	//log.Debug("mv command = %s", mvCommand)
	//es, err = destExecutor.ExecShell(mvCommand)
	if err != nil {
		lastErr = err
	} else {
		if es.ExitCode != 0 {
			var errMsg string
			if len(es.Stderr) == 0 {
				errMsg = executor.ErrMsgUnknow
			} else {
				errMsg = strings.Join(es.Stderr, "\n")
			}
			lastErr = fmt.Errorf(errMsg)
		}
	}
	log.Debug("lastError= %v", lastErr)
	if lastErr != nil {
		if strings.Contains(lastErr.Error(), "File exists") ||
			strings.Contains(lastErr.Error(),"文件已存在"){
			lastErr = nil
		}
	}
	return lastErr
}
