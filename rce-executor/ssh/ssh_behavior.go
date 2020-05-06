package ssh

import (
	"strings"
	"time"
)

const (
	DefaultExecuteTimeout = time.Second * 5
)

//ExecuteSetTimeout 判断是否需要设置执行超时
func ExecuteSetTimeout(cmd string) bool {
	if strings.Contains(cmd, "wget") { //对于下载不设置执行超时
		return false
	}
	return true
}
