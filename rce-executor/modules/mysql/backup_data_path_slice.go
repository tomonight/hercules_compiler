package mysql

import (
	"time"
)

//BackUpDataPathSlice 定义备份文件路径动态数组
type BackUpDataPathSlice []string

//Valid 判断字符格式是否合法
func (c BackUpDataPathSlice) Valid() bool {
	for _, path := range c {
		_, err := time.Parse("20060102150405", path)
		if err != nil {
			return false
		}
	}
	return true
}

//Len 增量备份个数
func (c BackUpDataPathSlice) Len() int {
	return len(c)
}

//Swap 交换位置
func (c BackUpDataPathSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

//Less 比较两个数字字符串的大小
func (c BackUpDataPathSlice) Less(i, j int) bool {
	iTime, _ := time.Parse("20060102150405", c[i])
	jTime, _ := time.Parse("20060102150405", c[j])
	return iTime.Unix() < jTime.Unix()
}
