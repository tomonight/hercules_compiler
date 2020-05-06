package ssh

import (
	"fmt"
	"hercules_compiler/rce-executor/log"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	gossh "golang.org/x/crypto/ssh"

	"bufio"
	"io"
	"os"
	"path"

	"github.com/pkg/sftp"
)

//CmdStatus 命令的执行信息
type CmdStatus struct {
	StartTime int64
	StopTime  int64
	ExitCode  int64
	Stdout    []string
	Stderr    []string
	Error     string
}

//Client 定义ssh客户端
type Client struct {
	Host       string        //Host地址
	Username   string        //用户名
	Password   string        //密码
	Port       int           //端口号
	Client     *gossh.Client //ssh客户端
	LastResult string        //最近一次Run的结果
	KeyFile    string        // public key file路径
}

//默认超时时间
const (
	DefaultTimeOut = time.Second * 30
)

//Connect 连接目标主机
//
//return error 返回错误信息
func (c *Client) Connect() error {
	auth := []gossh.AuthMethod{gossh.Password(c.Password)}
	// 若密码为空，则使用证书登录
	if c.Password == "" {
		key, err := ioutil.ReadFile(c.KeyFile)
		if err != nil {
			return err
		}
		signer, err := gossh.ParsePrivateKey(key)
		if err != nil {
			return err
		}
		auth = []gossh.AuthMethod{gossh.PublicKeys(signer)}
	}
	config := &gossh.ClientConfig{
		User: c.Username,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key gossh.PublicKey) error {
			return nil
		},
		Timeout: DefaultTimeOut,
	}
	sshClient, err := gossh.Dial("tcp", c.Host+":"+strconv.Itoa(c.Port), config)
	if err != nil {
		c.Client = nil
		return err
	}
	c.Client = sshClient
	return nil
}

//Run 执行shell命令
//Run 函数通过shell 参数执行命令，执行结果通过error 判断是否成功
func (c *Client) Run(shell string) (string, error) {
	if c.Client == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)

	return c.LastResult, err
}

func (c *Client) RunCmdSetTimeout(shell string) (*CmdStatus, error) {
	cmdStatus := &CmdStatus{Stderr: []string{}, Stdout: []string{}, StartTime: time.Now().UnixNano(), StopTime: time.Now().UnixNano(), ExitCode: -1}
	var err error
	if ExecuteSetTimeout(shell) {
		log.Debug("this command %s need do timeout", shell)
		ch := make(chan int)
		go func(doChan chan int) {
			cmdStatus, err = c.RunCmd(shell)
			doChan <- 1
		}(ch)
		log.Debug("start do check timeout")
		select {
		case <-time.After(DefaultExecuteTimeout):
			log.Debug("command %s execute timeout", shell)
			return cmdStatus, fmt.Errorf("command %s execute timeout", shell)
		case <-ch:
			log.Debug("command %s execute finished", shell)
			break
		}
		log.Debug("execute command finished")
	} else {
		cmdStatus, err = c.RunCmd(shell)
	}
	return cmdStatus, err
}

//RunCmd 运行命令
// Add by xiongjun --20180513
func (c *Client) RunCmd(shell string) (*CmdStatus, error) {

	cmdStatus := &CmdStatus{Stderr: []string{}, Stdout: []string{}, StartTime: time.Now().UnixNano(), StopTime: time.Now().UnixNano(), ExitCode: -1}
	session, err := c.Client.NewSession()

	if err != nil {
		return cmdStatus, err
	}

	defer session.Close()

	//LANG=en_US.UTF8 设置语言环境
	session.Setenv("LANG", "en_US.UTF8")

	/*modes := gossh.TerminalModes{
		gossh.ECHO:          0, // disable echoing
		gossh.ECHOCTL:       0,
		gossh.TTY_OP_ISPEED: 115200, // input speed = 112kbaud
		gossh.TTY_OP_OSPEED: 115200, // output speed = 112kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		this.Close()
		return err
	}*/

	stdout, err := session.StdoutPipe()
	if err != nil {
		return cmdStatus, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return cmdStatus, err
	}

	outScanner := bufio.NewScanner(stdout)
	outScanner.Split(bufio.ScanLines)

	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(bufio.ScanLines)

	if err := session.Start(shell); err != nil {
		return cmdStatus, err
	}

	for outScanner.Scan() {
		cmdStatus.Stdout = append(cmdStatus.Stdout, outScanner.Text())
	}
	for errScanner.Scan() {
		cmdStatus.Stderr = append(cmdStatus.Stderr, errScanner.Text())
	}
	err = session.Wait()
	cmdStatus.StopTime = time.Now().UnixNano()
	//log.Debug("cmd string= %s ssh err %v", shell, err)
	if err != nil {
		switch v := err.(type) {
		case *gossh.ExitError:
			cmdStatus.ExitCode = int64(v.Waitmsg.ExitStatus())
			if cmdStatus.ExitCode == 0 {
				cmdStatus.ExitCode = -1
			}
			cmdStatus.Error = v.Waitmsg.Msg()
			return cmdStatus, nil
		}
		return cmdStatus, err
	}
	cmdStatus.ExitCode = 0
	return cmdStatus, nil
}

//Close 关闭客户端
func (c *Client) Close() {
	if c.Client != nil {
		c.Client.Close()
		c.Client = nil
	}
}

//SendFile 通过ssh传送本地文件到目标机器目录
//
//param localPath 	本地路径
//
//param remotePath 远程路径
//
func (c *Client) SendFile(localPath, remotePath string) error {
	var (
		err        error
		sftpClient *sftp.Client
	)
	err = c.Connect()
	if err != nil {
		return err
	}

	sftpClient, err = sftp.NewClient(c.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localPath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	return nil
}

//NewSSHClient 创建命令行对象
//
//param Host Host地址
//
//param username 用户名
//
//param password 密码
//
//param port 端口号
//
//param keyfile 证书文件路径，默认/zmysql/data/.ssh/id_rsa
//
func NewSSHClient(host, username, password, path string, port int, keyfile ...string) *Client {

	sshClient := new(Client)
	sshClient.Host = host
	sshClient.Username = username
	sshClient.Password = password
	sshClient.Port = port
	if len(keyfile) <= 0 {
		sshClient.KeyFile = os.ExpandEnv(path + "/data/.ssh/id_rsa")
	} else {
		sshClient.KeyFile = keyfile[0]
	}
	return sshClient
}
