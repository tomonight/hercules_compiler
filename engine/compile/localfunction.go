package compile

import (
	"errors"
	"fmt"
	"hercules_compiler/engine/syntax"
	executorPkg "hercules_compiler/rce-executor/executor"
	"strings"
	"time"
)

type engineExecutor struct {
	Host     string
	port     int
	executor executorPkg.Executor
}

type function interface {
	call(e *Engineer, args []syntax.Expr) (interface{}, error)
}

type basic struct {
	e    *Engineer
	args []syntax.Expr
}

type oprint struct {
	basic
}

type olength struct {
	basic
}

type oappend struct {
	basic
}

type osleep struct {
	basic
}

type oconnect struct {
	basic
}

type osetTarget struct {
	basic
}

type osetRunMode struct {
	basic
}

func factory(name string) function {
	switch name {
	case nprint:
		return &oprint{}
	case nlength:
		return &olength{}
	case nappend:
		return &oappend{}
	case nsleep:
		return &osleep{}
	case nconnect:
		return &oconnect{}
	case nsetTarget:
		return &osetTarget{}
	case nsetRunMode:
		return &osetRunMode{}
	default:
		//error
		return nil
	}
}

//print function
func (o *oprint) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	//do nothing
	if len(args) == 0 {
		return nil, nil
	}
	args = agrsConverse(args)
	fomats := []interface{}{}
	for _, arg := range args {
		v, err := e.expr(arg)
		if err != nil {
			//error
			return nil, err
		}
		fomats = append(fomats, v)
	}
	e.edone(fmt.Sprint(fomats...))
	return nil, nil
}

func (o *olength) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	args = agrsConverse(args)
	if len(args) != 1 {
		return nil, errors.New("expect only 1 param")
	}
	value, err := e.expr(args[0])
	if err != nil {
		return nil, err
	}
	basicValue := 0
	switch xv := value.(type) {
	case *syntax.CompositeLit:
		basicValue = len(xv.ElemList)
	case string:
		basicValue = len(xv)
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
		basicValue = 1
	case *syntax.BasicLit:
		if xv.Kind == syntax.StringLit {
			basicValue = len(xv.Value)
		} else {
			basicValue = 1
		}
	}
	return basicValue, nil
	// return &syntax.BasicLit{Value: fmt.Sprintf("%d", basicValue), Kind: syntax.IntLit}, nil
}

func (o *oappend) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	return nil, nil
}

func (o *osleep) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	args = agrsConverse(args)
	if len(args) != 1 {
		return nil, errors.New("unexpect sleep time")
	}
	value, err := e.expr(args[0])
	if err != nil {
		return nil, err
	}
	vt, ok := value.(int)
	if !ok {
		return nil, errors.New("unexpect sleep time")
	}
	e.edone(fmt.Sprintf("sleep %dms", vt))
	time.Sleep(time.Duration(int64(vt)) * time.Millisecond)
	return nil, nil
}

func (o *oconnect) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	argsMap, err := e.agrsConverseMap(args, []string{"protocol", "host", "port", "username", "password", "keyfile"})
	if err != nil {
		return nil, err
	}
	return e.connectTarget(argsMap)
}

func agrsConverse(args []syntax.Expr) []syntax.Expr {
	exprs := []syntax.Expr{}
	for _, arg := range args {
		argv := arg.(*syntax.KeyValueExpr)
		exprs = append(exprs, argv.Value)
	}
	return exprs
}

func (e *Engineer) agrsConverseMap(args []syntax.Expr, keys ...[]string) (map[string]interface{}, error) {
	argsMap := make(map[string]interface{})
	hasKey := argsCheck(args)
	if !hasKey && (len(keys) == 0 || len(keys[0]) < len(args)) {
		return nil, fmt.Errorf("invalid params")
	}
	for index, arg := range args {
		argv := arg.(*syntax.KeyValueExpr)
		keyStr := ""
		if hasKey {
			key := argv.Key.(*syntax.Name)
			keyStr = key.Value
		} else {
			keyStr = keys[0][index]
		}
		value, err := e.expr(argv.Value)
		if err != nil {
			return nil, err
		}

		argsMap[keyStr] = value
	}
	return argsMap, nil
}

func sshTargetValid(target *Target) error {
	if target.username == "" || target.host == "" {
		return fmt.Errorf("connect target:incorrect connect statement, require host and username parameters")
	}

	if target.password == "" && target.keyfile == "" {
		return fmt.Errorf("connect target:incorrect connect statement, require password or keyfile parameters")
	}

	if target.port == 0 {
		target.port = 22
	}
	return nil
}

func (e *Engineer) connectTarget(params map[string]interface{}) (interface{}, error) {

	target, err := checkTarget(params)
	if err != nil {
		return nil, fmt.Errorf("connect target:incorrect connect statement,%s", err.Error())
	}

	// delete(params, "protocol")
	// targetName := replaceToken(tc, threadContext, tokens[2])

	switch strings.ToLower(target.protocol) {
	case "ssh":
		if err := sshTargetValid(target); err != nil {
			return nil, err
		}

		if ex := e.checkExecutorExsit(target.host, target.port); ex != nil {
			e.edone(fmt.Sprintf("connect to ssh target statement: connect to target '%s:%d' success", target.host, target.port))
			return ex, nil
		}

		executor, err := executorPkg.NewSSHAgentExecutor(target.host, target.username, target.password, target.port, target.keyfile)
		if err != nil {
			return nil, fmt.Errorf("connect ssh target '%s:%d' failed: %s", target.host, target.port, err.Error())
		}

		newTarget := &engineExecutor{executor: executor, Host: target.host, port: target.port}
		e.addExecutor(newTarget)
		e.edone(fmt.Sprintf("connect to ssh target statement: connect to target '%s:%d' success", target.host, target.port))
		return newTarget, nil
	case "rce":
		if target.host == "" {
			return nil, fmt.Errorf("connect target:incorrect connect statement, require host parameters")
		}
		port := 5051
		if target.port != 0 {
			port = target.port
		}
		executor, err := executorPkg.NewRCEAgentExecutor(target.host, fmt.Sprintf("%d", target.port))
		if err != nil {
			return nil, fmt.Errorf("connect rce target '%s:%d' failed: %s", params["host"], port, err.Error())
		}

		if ex := e.checkExecutorExsit(target.host, port); ex != nil {
			e.edone(fmt.Sprintf("connect to rce target statement: connect to target '%s:%d' success", target.host, port))
			return ex, nil
		}

		newTarget := &engineExecutor{executor: executor, Host: target.host, port: port}
		e.addExecutor(newTarget)
		e.edone(fmt.Sprintf("connect to rce target statement: connect to target '%s:%d' success", target.host, port))
		return newTarget, nil

	//support zcloud agent
	case "zcagent":
		if target.host == "" {
			return nil, fmt.Errorf("connect target:incorrect connect statement, require host parameters")
		}
		//zcloud default agent port is 8100
		port := 8100
		if target.port != 0 {
			port = target.port
		}
		if ex := e.checkExecutorExsit(target.host, port); ex != nil {
			e.edone(fmt.Sprintf("connect to zcagent target statement: connect to target '%s:%d' success", target.host, port))
			return ex, nil
		}
		executor, err := executorPkg.NewZCAgentExecutor(target.host, target.port)
		if err != nil {
			return nil, fmt.Errorf("connect rce target '%s:%d' failed: %s", target.host, port, err.Error())
		}

		newTarget := &engineExecutor{executor: executor, Host: target.host, port: port}
		e.addExecutor(newTarget)
		e.edone(fmt.Sprintf("connect to zcagent target statement: connect to target '%s:%d' success", target.host, port))
		return newTarget, nil

		//连接到本地
	case "local":
		if ex := e.checkExecutorExsit(target.host, 0); ex != nil {
			e.edone(fmt.Sprintf("connect to local target statement success"))
			return ex, nil
		}
		executor, err := executorPkg.NewLocalExecutor()
		if err != nil {
			return nil, fmt.Errorf("connect local target failed:%s", err.Error())
		}
		newTarget := &engineExecutor{executor: executor, Host: "localhost", port: 0}
		e.addExecutor(newTarget)
		e.edone("connect to local target statement success")
		return newTarget, nil
	default:
		return nil, fmt.Errorf("connect target:not supported connect protocol")
	}
	return nil, fmt.Errorf("connect target failed:unknown error")
}

func (o *osetTarget) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	argsMap, err := e.agrsConverseMap(args, []string{"target"})
	if err != nil {
		return nil, err
	}
	if targetIn, ok := argsMap["target"]; ok {
		if executor, ok := targetIn.(*engineExecutor); ok {
			e.executor = executor
			e.target = executor.Host
			e.edone(fmt.Sprintf("set target statement: switch to target '%s'", executor.Host))
			return nil, nil
		}
	}
	return nil, fmt.Errorf("unknown target")
}

func (o *osetRunMode) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	args = agrsConverse(args)
	if len(args) != 1 {
		return nil, errors.New("unexpect mode value")
	}
	value, err := e.expr(args[0])
	if err != nil {
		return nil, err
	}
	vt, ok := value.(int)
	if !ok {
		return nil, errors.New("unexpect mode value")
	}
	if !checkStatusValid(vt) {
		return nil, errors.New("unexpect mode value")
	}
	e.execStatus = vt
	return nil, nil
}

func checkStatusValid(status int) bool {
	return status == 0 || status == 1 || status == 10 || status == 11
}
