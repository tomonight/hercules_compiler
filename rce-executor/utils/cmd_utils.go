package utils

import (
	"fmt"
	"os"
	"strings"
)

//EscapeShellCmd 转换双引号内容的shell语句
func EscapeShellCmd(shell string) string {
	shell = strings.TrimSpace(shell)
	r := ""
	for _, c := range shell {
		switch c {
		case '"':
			r += "\\"
		case '\\':
			r += "\\"
		case '$':
			r += "\\"
		}
		r += string(c)
	}
	return r
}

//EscapeSingleQuoteShellCmd 转换单引号内容的shell语句
func EscapeSingleQuoteShellCmd(shell string) string {
	shell = strings.TrimSpace(shell)
	r := ""
	for _, c := range shell {
		switch c {
		//case '"':
		case '\'':
			r += "'\\'"
			//case '$':
			//	r += `\`
		}
		r += string(c)
	}
	return r
}
func CheckShellCmd(shell string) error {

	shell = strings.TrimSpace(shell)

	if len(shell) == 0 {
		return nil
	}
	cnt := strings.Count(shell, "\"") //双引号不能为单个
	if cnt%2 != 0 {
		return fmt.Errorf("Quote character count can not bee odd number")
	}
	cnt = strings.Count(shell, "'") //单引号不能为单个
	if cnt%2 != 0 {
		return fmt.Errorf("Quote character count can not bee odd number")
	}

	if shell[len(shell)-1] == '\\' {
		return fmt.Errorf("Last character can not be \\")
	}
	return nil
}

// PathExist 检查文件或目录是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
