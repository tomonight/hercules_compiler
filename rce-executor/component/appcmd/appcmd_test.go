//安装agent 单元测试
package appcmd

import (
	"testing"
)

func TestAgentInstall(t *testing.T) {
	sourcePath := "/Users/daiwei/work/tmp/rce-agent"
	tagetPath := "/home/daiwei/tmp"

	err := InstallAgentInTarget(sourcePath, tagetPath, "192.168.99.254", "root", "123456", 22)
	if err != nil {
		t.Error("InstallAgentInTarget failed :", err)
		return
	}
}
