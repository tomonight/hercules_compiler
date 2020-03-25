package main

import (
	"hercules_compiler/compile"
	"io/ioutil"
	pa "path/filepath"
)

var scriptpath string

//Init initialial
//set complile path which need to complile
//parse AST to memory
func Init(path string) error {
	scriptpath = path
	fileNames := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		filepath := getfilePath(path, f.Name())
		fileNames = append(fileNames, filepath)
	}

	_ = compile.ParseFiles(fileNames)
	return nil
}

//getfilePath
//combine file path and filename with diffrent os
func getfilePath(path, name string) string {
	filePath := pa.Join(path, name)
	return filePath
}

//RunScript run single script synchronize
//name script name
//params json string ,allow [string int bool array map]
func RunScript(name, params string) error {
	node := compile.GetNoder(getfilePath(scriptpath, name))

	engineer, err := compile.NewEngineer(node, params)
	if err != nil {
		return err
	}
	return engineer.Run()
}

func main() {
	Init("test_script")
	_ = RunScript("backup.script", "")
}
