package syntax

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func Test_parse(t *testing.T) {
	type args struct {
		src io.Reader
	}
	f, _ := os.Open("test.script")
	tests := []struct {
		name    string
		args    args
		want    *File
		wantErr bool
	}{
		{name: "success", args: args{src: f}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := Parse(tt.args.src)
			fmt.Println(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
