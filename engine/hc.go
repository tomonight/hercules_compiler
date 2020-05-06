package main

import (
	"fmt"
	"hercules_compiler/engine/compile"
	_ "hercules_compiler/rce-executor/modules/cgroup"
	_ "hercules_compiler/rce-executor/modules/etcd"
	_ "hercules_compiler/rce-executor/modules/http"
	_ "hercules_compiler/rce-executor/modules/oscmd"
	_ "hercules_compiler/rce-executor/modules/osservice"
	_ "hercules_compiler/rce-executor/modules/zdata"
	"io/ioutil"
	pa "path/filepath"
	"strings"
)

type KK struct {
	Name  string   `json:"name"`
	ID    uint     `json:"id"`
	Class []*Class `json:"class"`
}

type Class struct {
	Name string `json:"name"`
}

//Init initialial
//set complile path which need to complile
//parse AST to memory
func Init(path string) error {
	fileNames := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".hercules") {
			continue
		}
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
func RunScript(name string, params interface{}) error {
	node := compile.GetNoder(name)
	if node == nil {
		return fmt.Errorf("no script exsit")
	}
	engineer, err := compile.NewEngineer(node, params, newCallback())
	if err != nil {
		return err
	}
	engineer.Run()
	return nil
}

func main() {
	kk := &KK{Name: "kk", ID: 123}
	c1 := &Class{Name: "c1"}
	c2 := &Class{Name: "c2"}
	c3 := &Class{Name: "c3"}
	kk.Class = []*Class{}
	kk.Class = append(kk.Class, c1)
	kk.Class = append(kk.Class, c2)
	kk.Class = append(kk.Class, c3)
	Init("test_script")
	_ = RunScript("backup", kk)
}

func newCallback() compile.Callback {
	return func(scriptExecID string, callbackResultType int, callbackResultInfo string) {
		fmt.Println(fmt.Sprintf("scriptExecID:%s-----callbackResultType:%d-------callbackResultInfo:%s", scriptExecID, callbackResultType, callbackResultInfo))
	}
}
