package utils

import (
	"os"
	"path"
	"strings"
)

//FileExist 判断文件是否存在
func FileExist(filePath string) (bool, error) {
	var file *os.File
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	return true, nil
}

//DirExist 判断文件夹是否存在
func DirExist(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

//CreateDir 创建文件夹
func CreateDir(dir string) error {
	err := os.Mkdir(dir, os.ModePerm)
	return err
}

//DeleteFile 删除文件
func DeleteFile(filePath string) error {
	exist, _ := FileExist(filePath)
	if exist {
		return os.Remove(filePath)
	}
	return nil
}

//GetFileName 获取文件名
func GetFileName(fullFilename string) string {
	var filenameWithSuffix, fileSuffix, filenameOnly string
	filenameWithSuffix = path.Base(fullFilename)
	fileSuffix = path.Ext(filenameWithSuffix)
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameOnly
}
