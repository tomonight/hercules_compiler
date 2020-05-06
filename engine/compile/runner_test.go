package compile

import (
	"testing"
)

type person struct {
	A string              `json:"a"`
	B string              `json:"b"`
	C int                 `json:"c"`
	D []map[string]string `json:"d"`
}

func Test_getParams(t *testing.T) {
	type args struct {
		params interface{}
	}
	m := make(map[string]string)
	m["d"] = "d"
	m["dd"] = "dd"
	m1 := make(map[string]string)
	m1["d"] = "d1"
	m1["dd"] = "dd1"
	r := []map[string]string{}
	r = append(r, m)
	r = append(r, m1)
	per := person{A: "a", B: "b", C: 3, D: r}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{name: "success", args: args{params: per}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getParams(tt.args.params)
		})
	}
}
