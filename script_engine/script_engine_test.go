package script_engine

import (
	"encoding/json"
	"fmt"
	"hercules_compiler/rce-executor/log"
	beego "mydata-beego"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// var a = `[{\"slaveAccount\":\"123456\",\"slaveHost\":\"192.168.11.222\",\"slavePort":"10003","slaveUsername":"root"},{"slaveAccount":"123456","slaveHost":"192.168.11.223","slavePort":"10003","slaveUsername":"root"}]`
var a = `[{"b":"1"},{"b":"2"}]`

func TestAddDbUser(t *testing.T) {

	Convey("测试能否添加成功", t, func() {
		// param := []map[string]string{}
		// map1 := map[string]string{"a": "1"}
		// param = append(param, map1)
		// map2 := map[string]string{"a": "2"}
		// param = append(param, map2)
		// list, _ := ListToJSONString(param)
		script := []string{}
		script = append(script, `connect target target1 protocol=ssh host=192.168.11.181 username=root port=22 password=root123`)
		script = append(script, `set target target1`)
		script = append(script, `execute oscmd.GetMGRWhiteList`)
		script = append(script, "print ${{value}}")
		script = append(script, "set var b=${{value}}")
		script = append(script, `for ${b} `)
		script = append(script, `print ${a} `)
		script = append(script, `endfor `)
		// script = append(script, `finally failed`)
		// script = append(script, `if "111" == "111"`)
		// script = append(script, `test "dakwhdkjawhdkahwdawhdjahwdhawddddddddddd"`)
		// script = append(script, `endif`)
		// script = append(script, `endfinally`)
		// script = append(script, `finally success`)
		// script = append(script, `test "112312"`)
		// script = append(script, `endfinally`)
		// script = append(script, `begin parallel`)
		// script = append(script, `test "223333"`)
		// script = append(script, `parallel target target1`)
		// script = append(script, `print "aaaaaaaa"`)
		// // script = append(script, `test "11"`)
		// script = append(script, `parallel target target2`)
		// // script = append(script, `return "11111"`)
		// script = append(script, `test "111"`)
		// script = append(script, `test "2222"`)
		// script = append(script, `test "3333"`)
		// script = append(script, `test "4444"`)
		// script = append(script, `test "5555"`)
		// script = append(script, `return "22222"`)
		// script = append(script, `test "7777"`)
		// script = append(script, `if "1" == "2"`)
		// script = append(script, `test "555555555"`)
		// script = append(script, `endif`)
		// script = append(script, `test "8888"`)
		// script = append(script, `test "9999"`)
		// script = append(script, `end parallel`)
		tc := scriptContext{}
		tc.staging = TCS_PARSING
		tc.mainThreadContext = &scriptThreadContext{}
		currentThreadContext := tc.mainThreadContext
		currentThreadContext.envs = make(map[string]string)
		currentThreadContext.vars = make(map[string]string)
		currentThreadContext.executeFailedFlag = EFF_STOP
		currentThreadContext.executeSuccessfulFlag = EFS_CONTINUE
		tc.targets = make(map[string]*engineExecutor)
		tc.inParallel = false
		tc.threadContexts = nil
		tc.expressBlock = make(map[int]*scriptBlockTokens)
		allTokens, err, _ := parseScript(script, &tc)
		fmt.Println(allTokens)
		So(err, ShouldBeNil)
		params := make(map[string]string)
		aa := []map[string]string{}
		bb := map[string]string{"kk": "1awdawd1"}
		cc := map[string]string{"kk": "11awdawd"}
		dd := map[string]string{"kk": "1awdawd1"}
		aa = append(aa, bb)
		aa = append(aa, cc)
		aa = append(aa, dd)
		b, err := json.Marshal(aa)
		if err != nil {
			log.Error("json.Marshal failed:", err)
		}
		params["aa"] = string(b)

		ExecuteScript(script, params, "1", nil)
	})
}

func ListToJSONString(param []map[string]string) (str string, err error) {

	b, err := json.Marshal(param)
	if err != nil {
		beego.Error("json.Marshal failed:", err)
		return "", err
	}

	return string(b), nil
}
