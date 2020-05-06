package main

import (
	"encoding/json"
	"fmt"
	"hercules_compiler/engine/compile"
	_ "hercules_compiler/rce-executor/modules/cgroup"
	_ "hercules_compiler/rce-executor/modules/etcd"
	_ "hercules_compiler/rce-executor/modules/http"
	_ "hercules_compiler/rce-executor/modules/oscmd"
	_ "hercules_compiler/rce-executor/modules/osservice"
	_ "hercules_compiler/rce-executor/modules/zdata"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "success", args: args{path: "test_script"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.path)
		})
		noders := compile.Noders

		for _, node := range noders {
			for _, err := range node.Errors() {
				fmt.Println(err.Error())
			}
		}
		fmt.Println(noders)
	}
}

func TestRunScript(t *testing.T) {

	kk := &KK{Name: "kk", ID: 123}
	c1 := &Class{Name: "c1"}
	c2 := &Class{Name: "c2"}
	c3 := &Class{Name: "c3"}
	kk.Class = []*Class{}
	kk.Class = append(kk.Class, c1)
	kk.Class = append(kk.Class, c2)
	kk.Class = append(kk.Class, c3)

	aa, _ := json.Marshal(kk)
	type args struct {
		name   string
		params interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{name: "backup", params: string(aa)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init("test_script")
			if err := RunScript(tt.args.name, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("RunScript() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
