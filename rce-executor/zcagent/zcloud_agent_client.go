package zcagent

import (
	"strconv"
	"strings"
	"time"

	"hercules_compiler/rce-executor/zcagent/command"

	"fmt"

	"github.com/nu7hatch/gouuid"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
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

// A Client calls a remote agent (server) to execute commands.
type Client interface {
	// Connect to a remote agent.
	Open() error

	// Close connection to a remote agent.
	Close() error

	// Return hostname and port of remote agent, if connected.
	AgentAddr() (string, int)

	//ExecuteCommand
	RunCmd(shell string, timeout int64) (*CmdStatus, error)
}

type client struct {
	host  string
	port  int
	conn  *grpc.ClientConn
	agent command.CommandV3ServerClient
}

// NewClient makes a new Client.
func NewClient(host string, port int) Client {
	return &client{host: host, port: port}
}

func (c *client) Open() error {
	var opt grpc.DialOption
	opt = grpc.WithInsecure()
	conn, err := grpc.Dial(
		c.host+":"+fmt.Sprintf("%d", c.port),
		opt, // insecure or with TLS

		// Block = actually connect. Timeout = max time to retry on failure
		// (no option to set retry count). Backoff delay = time between retries,
		// up to Timeout.
		grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(10)*time.Second),
		grpc.WithBackoffMaxDelay(time.Duration(2)*time.Second),
	)
	if err != nil {
		return err
	}
	c.conn = conn
	c.agent = command.NewCommandV3ServerClient(conn)
	return nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *client) AgentAddr() (string, int) {
	return c.host, c.port
}

func id() string {
	uuid, _ := uuid.NewV4()
	return strings.Replace(uuid.String(), "-", "", -1)
}

//RunCmd 执行远程shell命令，timeout以毫秒为单位
func (c *client) RunCmd(shell string, timeout int64) (*CmdStatus, error) {

	cmdStatus := &CmdStatus{Stderr: []string{}, Stdout: []string{}, StartTime: time.Now().UnixNano(), StopTime: time.Now().UnixNano(), ExitCode: -1}

	cmd := &command.CommandV3{
		CommandContent: shell,
		RequestId:      id(),
		TimeOut:        timeout,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout+1000)*time.Millisecond)
	defer cancel()

	result, err := c.agent.ExecuteCommand(ctx, cmd)
	cmdStatus.StopTime = time.Now().UnixNano()

	if err != nil {
		return cmdStatus, err
	}
	if result == nil {
		return cmdStatus, fmt.Errorf("Remote agent return nil message")
	}
	msg, err := result.Recv()
	if err != nil {
		return cmdStatus, err
	}
	if msg == nil {
		return cmdStatus, fmt.Errorf("Get empty message from remote agent")
	}
	if len(msg.ResultMsg) == 0 {
		cmdStatus.Stdout = []string{}
	} else {
		cmdStatus.Stdout = strings.Split(strings.Replace(msg.ResultMsg, "\r\n", "\n", -1), "\n")
	}

	if len(msg.ErrorMsg) == 0 {
		cmdStatus.Stderr = []string{}
	} else {
		cmdStatus.Stderr = strings.Split(strings.Replace(msg.ErrorMsg, "\r\n", "\n", -1), "\n")
	}
	cmdStatus.ExitCode, err = strconv.ParseInt(msg.ExecuteCode, 10, 64)
	if err != nil {
		cmdStatus.ExitCode = -1
		cmdStatus.Stderr = []string{"Invalid shell exit code " + msg.ExecuteCode}
	}
	return cmdStatus, nil
}
