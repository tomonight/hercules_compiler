package executor

import (
	"fmt"
	"hercules_compiler/rce-executor/utils"
	"hercules_compiler/rce-executor/zcagent"
	"strings"
	"time"
)

//默认的超时时间为24小时
const defaultZCAgentExecuteTimeout = 3600 * 24 * 1000

//ZCAgent zCloud agent 执行器
type ZCAgentExecutor struct {
	timeout           int
	agentClient       zcagent.Client
	context           map[string]string
	remoteWorkingPath string
	remoteWorkingUser string
	sudoEnabled       bool
	remoteWorkingEnv  map[string]string
}

//NewZCAgentExecutor 生成新的zcloud agent执行器
func NewZCAgentExecutor(host string, port int) (Executor, error) {
	e := ZCAgentExecutor{
		agentClient:      zcagent.NewClient(host, port),
		context:          map[string]string{},
		remoteWorkingEnv: map[string]string{},
	}
	err := e.agentClient.Open()
	if err != nil {
		return nil, err
	}
	e.initContext()
	return &e, nil
}

//getHostName 获取主机名
func (e *ZCAgentExecutor) getHostName() error {
	es, err := e.ExecShell("hostname")
	if err != nil {
		return err
	}
	if es.ExitCode != 0 {
		return fmt.Errorf("Can not get hostname ,error code=%d,error message=%v", es.ExitCode, es.Stderr)
	}
	if len(es.Stdout) < 1 {
		return fmt.Errorf("Can not get hostname , no stdout")
	}
	hostnameText := strings.SplitN(strings.TrimLeft(es.Stdout[0], " "), " ", 2)
	if len(hostnameText) < 1 {
		return fmt.Errorf("Can not get hostname, stdout=%v", es.Stdout)
	}
	e.context[ContextNameHostname] = hostnameText[0]
	return nil
}

//initContext 初始化上下文
func (e *ZCAgentExecutor) initContext() error {

	es, err := e.ExecShell("uname -a")
	if err != nil {
		return err
	}
	if es.ExitCode != 0 {
		return fmt.Errorf("Can not get operation type,error code=%d,error message=%v", es.ExitCode, es.Stderr)
	}
	if len(es.Stdout) < 1 {
		return fmt.Errorf("Can not get operation type, no stdout")
	}
	ss := strings.SplitN(strings.TrimLeft(es.Stdout[0], " "), " ", 2)
	if len(ss) < 1 {
		return fmt.Errorf("Can not get operation type, stdout=%v", es.Stdout)
	}

	switch strings.ToLower(ss[0]) {
	case "linux":
		err := e.getLinuxDist()
		if err != nil {
			return err
		}
		e.context[ContextNameOSType] = Linux
		e.context[ContextNamePathSeparator] = "/"
		e.context[ContextNameTempDir] = "/tmp"
	case "darwin":
		e.context[ContextNameOSType] = MacOS
		e.context[ContextNamePathSeparator] = "/"
		e.context[ContextNameTempDir] = "/tmp"

	default:
		return fmt.Errorf("Can not get operation system type, stdout=%v", es.Stdout)
	}

	err = e.getHostName()
	if err == nil {
		return e.getLinuxDist()
	}
	return err
}

//Close 关闭连接
func (e *ZCAgentExecutor) Close() {
	e.agentClient.Close()
}

//GetExecutorContext 通过上下文名获取上下文值
func (e *ZCAgentExecutor) GetExecutorContext(contextName string) string {
	v, _ := e.context[contextName]
	return v
}

//SetTimeOut SetTimeOut
func (e *ZCAgentExecutor) SetTimeOut(timeout int) {
	e.timeout = timeout
}

//GetTimeOut GetTimeOut
func (e *ZCAgentExecutor) GetTimeOut() int {

	return e.timeout
}

//ExecShell 执行shell命令
func (e *ZCAgentExecutor) ExecShell(shellcmd string) (*ExecutedStatus, error) {
	timeout := e.GetTimeOut()
	if timeout > 0 {
		shellcmd = fmt.Sprintf("timeout %d %s", timeout, shellcmd)
	}

	if err := utils.CheckShellCmd(shellcmd); err != nil {
		return nil, err
	}

	if len(e.remoteWorkingPath) > 0 {
		shellcmd = "cd " + e.remoteWorkingPath + ";" + shellcmd
	}

	var envStr string
	envStr = ""

	for k, v := range e.remoteWorkingEnv {
		envStr = envStr + "export " + k + "=" + v + ";"
	}

	shellcmd = envStr + shellcmd

	if len(e.remoteWorkingUser) > 0 {
		shellcmd = "su " + e.remoteWorkingUser + " -c '" + utils.EscapeSingleQuoteShellCmd(shellcmd) + "'"
	}

	if e.sudoEnabled {
		shellcmd = "sudo -nE bash -c '" + utils.EscapeSingleQuoteShellCmd(shellcmd) + "'"
	}
	//fmt.Printf("cmd='%s'", shellcmd)

	startTime := time.Now().UnixNano()
	result, err := e.agentClient.RunCmd(shellcmd, defaultZCAgentExecuteTimeout)

	if err != nil {
		return nil, err
	}
	stopTime := time.Now().UnixNano()
	return &ExecutedStatus{RemoteStartTime: result.StartTime, RemoteStopTime: result.StopTime,
		ExitCode: result.ExitCode, Stderr: result.Stderr, Stdout: result.Stdout, ErrorMessage: result.Error,
		StartTime: startTime, StopTime: stopTime}, nil

}

//SetWorkingPath 设置工作目录
func (e *ZCAgentExecutor) SetWorkingPath(workingPath string) {
	e.remoteWorkingPath = workingPath
}

//SetEnv 设置环境变量
func (e *ZCAgentExecutor) SetEnv(envName, envValue string) {
	e.remoteWorkingEnv[envName] = envValue
}

//SetExecuteUser 设置用户
func (e *ZCAgentExecutor) SetExecuteUser(username string) {
	e.remoteWorkingUser = username
}

func (e *ZCAgentExecutor) SetSudoEnabled(sudoEnabled bool) {
	e.sudoEnabled = sudoEnabled
}

//getLinuxDist 系统版本信息
func (e *ZCAgentExecutor) getLinuxDist() error {
	var (
		redhat   = "/etc/redhat-release"
		suse     = "/etc/SUSE-release"
		mandrake = "/etc/mandrake-release"
		debian   = "/etc/debian_version"
	)

	var (
		kernel       string
		architecture string
		dist         string
		version      string
		psuedoname   string
	)

	//@remark do redhat check
	//PSUEDONAME=`cat /etc/redhat-release | sed s/.*\(// | sed s/\)//`
	sedParamHead := "sed s/.*\\(//" //替换 ( 的前半部分为空字符串
	sedParamBack := "sed s/\\)//"   //替换后）的半部分为空
	redhatCmd := fmt.Sprintf("cat %s | %s | %s", redhat, sedParamHead, sedParamBack)
	//fmt.Println("cmd =", redhatCmd)
	res, err := e.ExecShell(redhatCmd)
	//fmt.Println("res=", res, " err = ", err)
	if err == nil {
		dist = "RedHat"
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			//fmt.Println("stdout=", res.Stdout)
			psuedoname = res.Stdout[0]
		}
		//REV=`cat /etc/redhat-release | sed s/.*release\ // | sed s/\ .*//`
		sedParamHead = "sed s/.*release\\ //" //替换relase 前半部分为空字符串
		sedParamBack = "sed s/\\ .*//"
		redhatCmd = fmt.Sprintf("cat %s | %s | %s", redhat, sedParamHead, sedParamBack)
		res, err = e.ExecShell(redhatCmd)
		if err == nil {
			if stdoutLen := len(res.Stdout); stdoutLen > 0 {
				version = res.Stdout[0]
			}
		} else {
			return err
		}

	} else {
		//@remark do suse check
		//DIST=`cat /etc/SUSE-release | tr "\n" ' '| sed s/VERSION.*//`
		trHead := "tr \"\n\" ' '"         //换行符替换为空格
		sedParamBack = "sed s/VESION.*//" //VERSION后面的全部替换为空字符
		suseCmd := fmt.Sprintf("cat %s | %s | %s", suse, trHead, sedParamBack)
		res, err = e.ExecShell(suseCmd)
		if err == nil {
			if stdoutLen := len(res.Stdout); stdoutLen > 0 {
				dist = res.Stdout[0]
			}

			//REV=`cat /etc/SUSE-release | tr "\n" ' ' | sed s/.*=\ //`
			sedParamBack = "sed s/.*=\\ //"
			suseCmd = fmt.Sprintf("cat %s | %s | %s", suse, trHead, sedParamBack)
			res, err = e.ExecShell(suseCmd)
			if err == nil {

				if stdoutLen := len(res.Stdout); stdoutLen > 0 {
					version = res.Stdout[0]
				}
			} else {
				return err
			}

		} else {
			//@remark do  mandrake check
			//PSUEDONAME=`cat /etc/mandrake-release | sed s/.*\(// | sed s/\)//`
			mandrakeCmd := fmt.Sprintf("cat %s | %s | %s", mandrake, sedParamHead, sedParamBack)
			res, err = e.ExecShell(mandrakeCmd)
			if err == nil {
				dist = "Mandrake"
				if stdoutLen := len(res.Stdout); stdoutLen > 0 {
					psuedoname = res.Stdout[0]
				}

				//REV=`cat /etc/mandrake-release | sed s/.*release\ // | sed s/\ .*//`
				sedParamHead = "s/.*release\\ //" //替换relase 前半部分为空字符串
				sedParamBack = "s/\\ .*//"
				mandrakeCmd = fmt.Sprintf("cat %s | %s | %s", mandrake, sedParamHead, sedParamBack)
				res, err = e.ExecShell(mandrakeCmd)
				if err == nil {
					if stdoutLen := len(res.Stdout); stdoutLen > 0 {
						version = res.Stdout[0]
					}
				} else {
					return err
				}

			} else {
				//@remark do debian check
				//DIST = "Debian `cat /etc/debian_version`"
				debianCmd := fmt.Sprintf("cat %s", debian)
				res, err = e.ExecShell(debianCmd)
				if err == nil {
					if stdoutLen := len(res.Stdout); stdoutLen > 0 {
						dist = "Debian " + res.Stdout[0]
					}
					version = ""
				} else {
					return err
				}

			}
		}
	}

	machCmd := "uname -m"
	res, err = e.ExecShell(machCmd)
	if err == nil {
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			architecture = res.Stdout[0]
		}
	} else {
		return err
	}

	kernelCmd := "uname -r"
	res, err = e.ExecShell(kernelCmd)
	if err == nil {
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			kernel = res.Stdout[0]
		}
	} else {
		return err
	}

	e.context[ContextNameDist] = dist
	e.context[ContextNameVersion] = version
	e.context[ContextNamePsuedoname] = psuedoname
	e.context[ContextNameArchitecture] = architecture
	e.context[ContextNameKernel] = kernel
	return nil
}
