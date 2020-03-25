package main

import (
	"fmt"
	"hercules_compiler/compile"
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
	type args struct {
		name   string
		params string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{name: "backup.script", params: "{\"a\":\"123\"}"}, wantErr: false},
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
