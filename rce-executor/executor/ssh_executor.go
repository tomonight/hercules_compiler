package executor

import (
	"errors"
	"fmt"
	"hercules_compiler/rce-executor/ssh"
	"hercules_compiler/rce-executor/utils"
	"log"
	"strings"
	"time"

	"github.com/pkg/sftp"
)

const SIZE = 1 << 15

//SSHAgentExecutor ssh执行器
type SSHAgentExecutor struct {
	timeout           int
	SshClient         ssh.Client
	FTPClient         *sftp.Client
	context           map[string]string
	remoteWorkingPath string
	remoteWorkingUser string
	sudoEnabled       bool
	remoteWorkingEnv  map[string]string
}

//NewSSHAgentExecutorForSSHClient 通过sshclient生成新的ssh执行器
func NewSSHAgentExecutorForSSHClient(sshClient *ssh.Client) (Executor, error) {
	if sshClient == nil {
		return nil, errors.New("sshclient is nil")
	}
	e := SSHAgentExecutor{
		SshClient:        *sshClient,
		context:          map[string]string{},
		remoteWorkingEnv: map[string]string{},
	}
	err := e.SshClient.Connect()
	if err != nil {
		return nil, err
	}
	e.initContext()
	return &e, nil
}

//NewSSHAgentExecutor 生成新的ssh执行器
func NewSSHAgentExecutor(host, username, password string, port int, keyfile ...string) (Executor, error) {
	path, err := GetSShPath()
	if err != nil {
		return nil, err
	}
	e := SSHAgentExecutor{
		SshClient:        *ssh.NewSSHClient(host, username, password, path, port, keyfile...),
		context:          map[string]string{},
		remoteWorkingEnv: map[string]string{},
	}
	err = e.SshClient.Connect()
	if err != nil {
		return nil, err
	}
	e.initContext()
	return &e, nil
}

//NewSSHAgentExecutor 生成新的ssh执行器,返回为是
func NewSSHAgentExecutorReturnStruct(host, username, password string, port int, keyfile ...string) (*SSHAgentExecutor, error) {
	path, err := GetSShPath()
	if err != nil {
		return nil, err
	}
	e := SSHAgentExecutor{
		SshClient:        *ssh.NewSSHClient(host, username, password, path, port, keyfile...),
		context:          map[string]string{},
		remoteWorkingEnv: map[string]string{},
	}
	err = e.SshClient.Connect()
	if err != nil {
		return nil, err
	}
	e.initContext()
	c, err := sftp.NewClient(e.SshClient.Client, sftp.MaxPacket(SIZE))
	if err != nil {
		log.Printf("unable to start sftp subsytem: %v", err)
		return nil, err
	}
	e.FTPClient = c
	return &e, nil
}

//getHostInfor 获取主机信息
func (e *SSHAgentExecutor) getHostInfor() error {
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
	// 获取CPU信息，包括处理器数量、物理核心数及型号
	processorCountCmd := "cat /proc/cpuinfo| grep 'processor'| wc -l"
	physicalCoreCountCmd := "cat /proc/cpuinfo| grep 'physical id'| sort| uniq| wc -l"
	modelNameCmd := "cat /proc/cpuinfo | grep 'model name' | cut -f2 -d: | uniq -c"
	cmdstr := fmt.Sprintf("%s;%s;%s", processorCountCmd, physicalCoreCountCmd, modelNameCmd)
	es, err = e.ExecShell(cmdstr)
	if es.ExitCode != 0 {
		return fmt.Errorf("Can not get cpu information, error code=%d, error message=%v", es.ExitCode, es.Stderr)

	}
	e.context[ContextNameCPUProcessorCount] = es.Stdout[0]
	e.context[ContextNameCPUPhysicalCoreCount] = es.Stdout[1]
	e.context[ContextNameCPUModelName] = strings.TrimSpace(es.Stdout[2])

	// 获取内存信息
	cmdstr = "cat /proc/meminfo |grep MemTotal|cut -d: -f2"
	es, err = e.ExecShell(cmdstr)
	if es.ExitCode != 0 {
		return fmt.Errorf("Can not get memory size, error code=%d, error message=%v", es.ExitCode, es.Stderr)
	}
	memorySize := es.Stdout[0]
	memorySize = strings.Replace(memorySize, "kB", "", -1)
	memorySize = strings.TrimSpace(memorySize)
	e.context[ContextNameMemorySize] = memorySize
	return nil
}

//initContext 初始化上下文
func (e *SSHAgentExecutor) initContext() error {

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

	err = e.getHostInfor()
	if err == nil {
		return e.getLinuxDist()
	}
	// return e.setEnUSEnvLang()
	return nil
}

//export LANG=en_US.UTF8
func (e *SSHAgentExecutor) setEnUSEnvLang() error {
	command := "export LANG=en_US.UTF8"
	es, err := e.ExecShell(command)
	if err != nil {
		return err
	}
	if es.ExitCode != 0 {
		return fmt.Errorf("%s,error code=%d,error message=%v", command, es.ExitCode, es.Stderr)
	}
	return err
}

//Close 关闭连接
func (e *SSHAgentExecutor) Close() {
	e.SshClient.Close()
}

//GetExecutorContext 通过上下文名获取上下文值
func (e *SSHAgentExecutor) GetExecutorContext(contextName string) string {
	v, _ := e.context[contextName]
	return v
}

//SetTimeOut SetTimeOut
func (e *SSHAgentExecutor) SetTimeOut(timeout int) {
	e.timeout = timeout
}

//GetTimeOut GetTimeOut
func (e *SSHAgentExecutor) GetTimeOut() int {
	return e.timeout
}

//ExecShell 执行shell命令
func (e *SSHAgentExecutor) ExecShell(shellcmd string) (*ExecutedStatus, error) {
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

	result, err := e.SshClient.RunCmd(shellcmd)

	if err != nil {
		return nil, err
	}
	stopTime := time.Now().UnixNano()
	return &ExecutedStatus{RemoteStartTime: result.StartTime, RemoteStopTime: result.StopTime,
		ExitCode: result.ExitCode, Stderr: result.Stderr, Stdout: result.Stdout, ErrorMessage: result.Error,
		StartTime: startTime, StopTime: stopTime}, nil

}

//SetWorkingPath 设置工作目录
func (e *SSHAgentExecutor) SetWorkingPath(workingPath string) {
	e.remoteWorkingPath = workingPath
}

//SetEnv 设置环境变量
func (e *SSHAgentExecutor) SetEnv(envName, envValue string) {
	e.remoteWorkingEnv[envName] = envValue
}

//SetExecuteUser 设置用户
func (e *SSHAgentExecutor) SetExecuteUser(username string) {
	e.remoteWorkingUser = username
}

func (e *SSHAgentExecutor) SetSudoEnabled(sudoEnabled bool) {
	e.sudoEnabled = sudoEnabled
}

//getLinuxDist 系统版本信息
func (e *SSHAgentExecutor) getLinuxDist() error {
	var (
		centos   = "/etc/centos-release"
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

	//@remark do centos check
	//PSUEDONAME=`cat /etc/centos-release | sed s/.*\(// | sed s/\)//`
	sedParamHead := "sed s/.*\\(//" //替换 ( 的前半部分为空字符串
	sedParamBack := "sed s/\\)//"   //替换后）的半部分为空
	centosCmd := fmt.Sprintf("cat %s | %s | %s", centos, sedParamHead, sedParamBack)
	res, err := e.ExecShell(centosCmd)
	if err == nil && len(res.Stderr) == 0 {
		dist = "CentOS"
		if stdoutLen := len(res.Stdout); stdoutLen > 0 {
			//fmt.Println("stdout=", res.Stdout)
			psuedoname = res.Stdout[0]
		}
		//REV=`cat /etc/centos-release | sed s/.*release\ // | sed s/\ .*//`
		sedParamHead = "sed s/.*release\\ //" //替换relase 前半部分为空字符串
		sedParamBack = "sed s/\\ .*//"
		centosCmd = fmt.Sprintf("cat %s | %s | %s", centos, sedParamHead, sedParamBack)
		res, err = e.ExecShell(centosCmd)
		if err == nil {
			if stdoutLen := len(res.Stdout); stdoutLen > 0 {
				version = res.Stdout[0]
			}
		} else {
			return err
		}
	} else {
		//@remark do redhat check
		//PSUEDONAME=`cat /etc/redhat-release | sed s/.*\(// | sed s/\)//`
		sedParamHead = "sed s/.*\\(//" //替换 ( 的前半部分为空字符串
		sedParamBack = "sed s/\\)//"   //替换后）的半部分为空
		redhatCmd := fmt.Sprintf("cat %s | %s | %s", redhat, sedParamHead, sedParamBack)
		//fmt.Println("cmd =", redhatCmd)
		res, err = e.ExecShell(redhatCmd)
		//fmt.Println("res=", res, " err = ", err)
		if err == nil && len(res.Stderr) == 0 {
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
			if err == nil && len(res.Stderr) == 0 {
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
				if err == nil && len(res.Stderr) == 0 {
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
					if err == nil && len(res.Stderr) == 0 {
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
func GetSShPath() (path string, err error) {
	localExecutor, err := NewLocalExecutor()
	out, err := localExecutor.ExecShell("grep '^WorkingDirectory' /etc/systemd/system/mydata_service.service|awk -F'WorkingDirectory=' '{print $2}'")
	if err != nil {
		return
	}
	if out.ErrorMessage == "" && len(out.Stdout) > 0 {
		path = out.Stdout[0]
	}
	localExecutor.Close()
	return
}
