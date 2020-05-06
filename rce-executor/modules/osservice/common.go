//Package  osservice common.go 定义服务命令的相关常量、以及通用处理函数

package osservice

import (
	"strings"
)

//定义服务操作常见的错误信息
const (
	//ServiceNotLoadedErrMsg 服务未加载
	ServiceNotLoadedErrMsg = "not loaded"
	//ServiceAccessDeniedErrMsg 服务无权限访问
	ServiceAccessDeniedErrMsg = "Access denied"
	//ServiceFileNotExistErrMsg 服务文件不存在
	ServiceFileNotExistErrMsg = "No such file or directory"
)

//定义服务操作
const (
	//ServiceActionStart 启动服务
	ServiceActionStart = "start"
	//ServiceActionStop 停止服务
	ServiceActionStop = "stop"
	//SerciceActionEnable 服务自动启动
	SerciceActionEnable = "enable"
	//ServiceActionDisable 删除服务
	ServiceActionDisable = "disable"
)

//ServiceErrorCanIgnore 判断服务操作返回的错误信息是否可忽略
//目前支持 stop disable 操作的错误信息
func ServiceErrorCanIgnore(action, msg string) bool {
	ignoreErrMsgList := []string{
		ServiceAccessDeniedErrMsg,
		ServiceNotLoadedErrMsg,
		ServiceFileNotExistErrMsg,
	}

	if action == ServiceActionDisable ||
		action == ServiceActionStop {
		for _, errMsg := range ignoreErrMsgList {
			if strings.Contains(msg, errMsg) {
				return true
			}
		}
	}
	return false
}
