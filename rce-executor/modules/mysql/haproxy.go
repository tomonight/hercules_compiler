package mysql

import (
	"fmt"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/utils"
)

//定义RCE 默认端口
const (
	DefaultRCEPort = "5051"
)

//stopHAProxyMySQLProxy 停止HAProxy mysql代理流量
func stopHAProxyMySQLProxy(ip, port string, servicePort uint) error {
	if port == "" {
		port = DefaultRCEPort
	}

	rceExecutor, err := executor.NewRCEAgentExecutor(ip, port)
	if err != nil {
		return err
	}

	//eg command echo "disable server mysql-3378/mysql-3378-1" | socat  unix-connect:/var/lib/haproxy/stats stdio
	command := `echo "disable server mysql-%d/mysql-%d-1" | socat  unix-connect:/var/lib/haproxy/stats stdio`
	command = fmt.Sprintf(command, servicePort, servicePort)
	es, err := rceExecutor.ExecShell(command)
	if err != nil {
		return err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		return err
	}

	//eg command echo "shutdown sessions server mysql-3378/mysql-3378-1" |   socat unix-connect:/var/lib/haproxy/stats stdio
	command = `echo "shutdown sessions server mysql-%d/mysql-%d-1" |   socat unix-connect:/var/lib/haproxy/stats stdio`
	command = fmt.Sprintf(command, servicePort, servicePort)
	es, err = rceExecutor.ExecShell(command)
	if err != nil {
		return err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		return err
	}
	return nil
}

//reloadHAProxy 重载HAProxy
func reloadHAProxy(ip, port, confText, confPath string) error {
	if port == "" {
		port = DefaultRCEPort
	}

	rceExecutor, err := executor.NewRCEAgentExecutor(ip, port)
	if err != nil {
		return err
	}
	command := fmt.Sprintf(`echo -e "%s" > %s`, utils.EscapeShellCmd(confText), confPath)
	es, err := rceExecutor.ExecShell(command)
	if err != nil {
		return err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		return err
	}

	command = "systemctl reload haproxy"
	es, err = rceExecutor.ExecShell(command)
	if err != nil {
		return err
	}
	err = executor.GetExecResult(es)
	if err != nil {
		return err
	}
	return nil
}
