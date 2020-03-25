package compile

import (
	"hercules_compiler/syntax"
	"testing"
)

func Test_isLogicOp(t *testing.T) {
	type args struct {
		op syntax.Operator
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "success", args: args{op: syntax.AndAnd}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLogicOp(tt.args.op); got != tt.want {
				t.Errorf("isLogicOp() = %v, want %v", got, tt.want)
			}
		})
	}
}
