package compile

import (
	"encoding/json"
	"errors"
	"fmt"
	"hercules_compiler/engine/syntax"
	"hercules_compiler/rce-executor/executor"
	"hercules_compiler/rce-executor/modules"
	"strconv"
	"strings"
)

//walkExpr statement runner
func (e *Engineer) walkExpr(expr *syntax.ExprStmt) (interface{}, error) {
	return e.expr(expr.X)
}

//walkexpr walk step
func (e *Engineer) expr(expr syntax.Expr) (interface{}, error) {

	switch ex := expr.(type) {
	case *syntax.BasicLit:
		return basicLit(ex), nil
	case *syntax.CompositeLit:
		return ex, nil
	case *syntax.CallExpr:
		return e.callExpr(ex)
	case *syntax.Name:
		return e.nameVal(ex.Value)
	case *syntax.Operation:
		return e.operationExpr(ex)
	case *syntax.ParenExpr:
		return e.expr(ex.X)
	case *syntax.IndexExpr:
		return e.indexExpr(ex)
	default:
		return nil, nil
	}
}

func basicLit(lit *syntax.BasicLit) interface{} {
	switch lit.Kind {
	case syntax.IntLit:
		intV, _ := strconv.Atoi(lit.Value)
		return intV
	case syntax.StringLit:
		return strings.Trim(lit.Value, "\"")
	case syntax.BoolLit:
		if lit.Value == "true" {
			return true
		}
		return false
	}
	return nil
}

func (e *Engineer) callExpr(expr *syntax.CallExpr) (interface{}, error) {

	switch t := expr.Fun.(type) {
	case *syntax.Name:
		return e.nameCall(t, expr.ArgList)
	case *syntax.SelectorExpr:
		return e.selectorCall(t, expr.ArgList)
	default:
		return nil, nil
	}
}

func (e *Engineer) nameCall(name *syntax.Name, args []syntax.Expr) (interface{}, error) {
	fname := name.Value
	if syntax.FuncExsit(name.Value) {
		//user define function call
		return e.namecall(name, args)
	}
	fuc := factory(fname)
	if fuc != nil {
		return fuc.call(e, args)
	}
	//error
	return nil, fmt.Errorf("unkhown function name '%s'", name.Value)
}

func (e *Engineer) namecall(name *syntax.Name, args []syntax.Expr) (interface{}, error) {
	//user define function call
	fun := syntax.GetFunc(name.Value)
	assign, err := e.defineFormParameter(fun, args)
	if err != nil {
		return nil, err
	}
	e.startFunc(name.Value)
	defer func() {
		e.endFunc()
	}()
	if len(assign) != 0 {
		for key, value := range assign {
			_, err := e.assign(key, value)
			if err != nil {
				return nil, err
			}
		}
	}
	for _, stmt := range fun.Body.List {
		_, err := e.walk(stmt)
		err = e.checkError(stmt.Pos(), err)
		if err != nil {
			return nil, err
		}
		if e.checkStatus() {
			break
		}
		if e.alreadyReturn() {
			return e.fnReturn(), nil
		}
	}
	//error
	return nil, nil
}

func (e *Engineer) selectorCall(fun *syntax.SelectorExpr, args []syntax.Expr) (interface{}, error) {
	if !argsCheck(args) {
		return nil, fmt.Errorf("selector call function must be key:value params")
	}
	if e.executor == nil {
		return nil, fmt.Errorf("function called must set target executor")
	}
	x, ok := fun.X.(*syntax.Name)
	if !ok {
		return nil, fmt.Errorf("selector call function is valid")
	}

	if !modules.CheckModuleExsit(x.Value) {
		return nil, fmt.Errorf("selector call function is valid")
	}

	commandName := fmt.Sprintf("%s.%s", x.Value, fun.Sel.Value)

	cmdFunc, ok := executor.GetCmdByModuleAndName(x.Value, fun.Sel.Value)

	if !ok {
		return nil, fmt.Errorf("execute command: can not found command '%s'", commandName)
	}
	params, err := e.agrsConverseMap(args)
	if err != nil {
		return nil, fmt.Errorf("execute command: incorrect execute statement, parse command parameter %s", err.Error())
	}
	execParams := executor.ExecutorCmdParams{}
	for k, v := range params {
		execParams[k] = v
	}
	// for k, v := range threadContext.envs {
	// 	target.executor.SetEnv(k, v)
	// }

	//add timeout param
	if timeout, exsit := params["timeout"]; exsit {
		timeoutInt, _ := strconv.Atoi(timeout.(string))
		e.executor.executor.SetTimeOut(timeoutInt)
	}
	er := cmdFunc(e.executor.executor, &execParams)

	if !er.Successful {
		if er.ExitCode == 124 {
			er.Message = "execute timeout"
		}
		return nil, fmt.Errorf("execute command failed: '%s'", er.Message)
	}

	// threadContext.resultData = er.ResultData
	lit := assembleMapCompositeLit(er.ResultData)
	e.edone(fmt.Sprintf("%s on target host '%s': %s", commandName, e.target, er.Message))
	return lit, nil
}

//define form parameter into enginner
func (e *Engineer) defineFormParameter(fun *syntax.FuncDecl, args []syntax.Expr) (map[string]interface{}, error) {
	if len(fun.Type.ParamList) == 0 {
		return nil, nil
	}
	if len(args) != len(fun.Type.ParamList) {
		return nil, funPramsIsvalid
	}
	hasKey := argsCheck(args)
	assign := make(map[string]interface{})
	if hasKey {
		for _, arg := range args {
			argv := arg.(*syntax.KeyValueExpr)
			value, err := e.expr(argv.Value)
			if err != nil {
				return nil, err
			}
			key, ok := argv.Key.(*syntax.Name)
			if !ok {
				return nil, errors.New("function call args is valid")
			}
			assign[key.Value] = value
		}
	} else {
		for index, p := range fun.Type.ParamList {
			argv, ok := args[index].(*syntax.KeyValueExpr)
			if !ok {
				return nil, fmt.Errorf("params parse error")
			}
			value, err := e.expr(argv.Value)
			if err != nil {
				return nil, err
			}
			assign[p.Name.Value] = value
		}
	}
	return assign, nil
}

func argsCheck(args []syntax.Expr) bool {
	if len(args) == 0 {
		return true
	}
	expr := args[0].(*syntax.KeyValueExpr)
	return expr.Key != nil
}

func (e *Engineer) nameExpr(expr *syntax.Name) (interface{}, error) {
	val, err := e.nameVal(expr.Value)
	if err != nil {
		return nil, err
	}
	if nil == val {
		//err
		return nil, nil
	}
	return val, nil
}

func (e *Engineer) operationExpr(expr *syntax.Operation) (interface{}, error) {

	xv, err := e.expr(expr.X)
	if err != nil {
		return nil, err
	}
	yv, err := e.expr(expr.Y)
	if err != nil {
		return nil, err
	}
	if isBoolLogicOp(expr.Op) {
		return e.logicop(xv, yv, expr.Op)
	}
	if isLitLogicop(expr.Op) {
		return e.litlogicop(xv, yv, expr.Op)
	}
	if isDigit(xv) && isDigit(yv) {
		return mathCalc(xv, yv, expr.Op)
	}
	if expr.Op == syntax.Add {
		return fmt.Sprintf("%v", xv) + fmt.Sprintf("%v", yv), nil
	}
	return nil, nil
}

//logic operator
const logicop uint64 = 1<<syntax.OrOr |
	1<<syntax.AndAnd |
	1<<syntax.Not

const litLogicop = 1<<syntax.Eql | // ==
	1<<syntax.Neq | // !=
	1<<syntax.Lss | // <
	1<<syntax.Leq | // <=
	1<<syntax.Gtr | // >
	1<<syntax.Geq // >=

//assert op is logic operator
func isBoolLogicOp(op syntax.Operator) bool {
	return logicop&(1<<op) > 0
}

//assert op is logic operator
func isLitLogicop(op syntax.Operator) bool {
	return litLogicop&(1<<op) > 0
}

func (e *Engineer) logicop(x, y interface{}, op syntax.Operator) (interface{}, error) {
	var xb, yb bool
	if _, ok := x.(bool); !ok {
		return nil, boolPramsIsvalid
	}
	if _, ok := y.(bool); !ok {
		return nil, boolPramsIsvalid
	}
	switch op {
	case syntax.OrOr:
		return xb || yb, nil
	case syntax.AndAnd:
		return xb && yb, nil
	case syntax.Not:
		return !xb, nil
	}
	return nil, nil
}

func (e *Engineer) litlogicop(xv, yv interface{}, op syntax.Operator) (interface{}, error) {
	t, ok := matchType(xv, yv)
	if !ok {
		return nil, errors.New("mix with compare type")
	}
	switch t {
	case ostr:
		x, y := xv.(string), yv.(string)
		switch op {
		case syntax.Eql:
			return x == y, nil
		case syntax.Neq:
			return x != y, nil
		case syntax.Lss:
			return x < y, nil
		case syntax.Leq:
			return x <= y, nil
		case syntax.Gtr:
			return x > y, nil
		case syntax.Geq:
			return x >= y, nil
		}
	case oint:
		x, _ := getDigit(xv)
		y, _ := getDigit(yv)
		switch op {
		case syntax.Eql:
			return x == y, nil
		case syntax.Neq:
			return x != y, nil
		case syntax.Lss:
			return x < y, nil
		case syntax.Leq:
			return x <= y, nil
		case syntax.Gtr:
			return x > y, nil
		case syntax.Geq:
			return x >= y, nil
		}
	}
	return nil, nil
}

func matchType(x, y interface{}) (uint, bool) {
	if isDigit(x) && isDigit(y) {
		return oint, true
	}
	switch x.(type) {
	case string:
		if isDigit(x) && isDigit(y) {
			return oint, true
		}
		if _, ok := y.(string); ok {
			return ostr, true
		}
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
		if _, ok := y.(int64); ok {
			return oint, true
		}
	case bool:
		if _, ok := y.(bool); ok {
			return obool, true
		}
	}
	return 0, false
}

func isDigit(value interface{}) bool {
	switch value.(type) {
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
		return true
	}
	bv, err := json.Marshal(value)
	if err != nil {
		return false
	}
	for _, b := range bv {
		if rune(b) != 34 && (rune(b) < rune(48) || rune(b) > rune(57)) {
			return false
		}
	}
	return true
}

func getDigit(value interface{}) (int, error) {
	switch value.(type) {
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8:
		return value.(int), nil
	case string:
		if isDigit(value) {
			xv, err := strconv.Atoi(value.(string))
			if err != nil {
				return 0, err
			}
			return int(xv), nil
		}
		return 0, notSupportOprationType
	default:
		return 0, notSupportOprationType
	}
}

func mathCalc(x, y interface{}, op syntax.Operator) (interface{}, error) {
	xv, err := getDigit(x)
	if err != nil {
		return nil, err
	}
	yv, err := getDigit(y)
	if err != nil {
		return nil, err
	}

	switch op {
	case syntax.Add:
		return xv + yv, nil
	case syntax.Sub:
		return xv - yv, nil
	case syntax.Mul:
		return xv * yv, nil
	case syntax.Div:
		return xv / yv, nil
	default:
		return nil, notSupportOprationType
	}
}

func (e *Engineer) indexExpr(expr *syntax.IndexExpr) (interface{}, error) {
	xv, err := e.expr(expr.X)
	if err != nil {
		return nil, err
	}

	if _, ok := xv.(*syntax.CompositeLit); !ok {
		return nil, errors.New("unsurport type ")
	}
	composite := xv.(*syntax.CompositeLit)
	indexv, err := e.expr(expr.Index)
	if err != nil {
		return nil, err
	}
	switch composite.Type.(type) {
	case *syntax.MapType:
		for _, exp := range composite.ElemList {
			if _, ok := exp.(*syntax.KeyValueExpr); !ok {
				return nil, errors.New("error map type")
			}
			keyExpr := exp.(*syntax.KeyValueExpr)
			// exprv, ok := keyExpr.Key.(*syntax.Name)
			// if !ok {
			// 	return nil, fmt.Errorf("key is not name type")
			// }
			var value interface{}
			switch ktype := keyExpr.Key.(type) {
			case *syntax.Name:
				value = ktype.Value
			case *syntax.BasicLit:
				value = basicLit(ktype)
			default:
				return nil, fmt.Errorf("key type got run type")
			}
			if indexv == value {
				return e.expr(keyExpr.Value)
			}
		}
	case *syntax.SliceType:
		if !isDigit(indexv) {
			return nil, errors.New("unsurport slice index ")
		}
		indexInt, _ := getDigit(indexv)
		return e.expr(composite.ElemList[indexInt])
	default:
		return nil, errors.New("unsurport type ")
	}
	return nil, nil
}

func assembleMapCompositeLit(params map[string]string) *syntax.CompositeLit {
	lit := &syntax.CompositeLit{}
	t := &syntax.MapType{}
	lit.Type = t
	lit.NKeys = len(params)
	lit.ElemList = []syntax.Expr{}
	for k, v := range params {
		kv := &syntax.KeyValueExpr{}
		kv.Key = &syntax.Name{Value: k}
		kv.Value = &syntax.BasicLit{Value: v, Kind: syntax.StringLit}
		lit.ElemList = append(lit.ElemList, kv)
	}
	return lit
}

func getBasicLit(value interface{}) (*syntax.BasicLit, error) {
	lit := &syntax.BasicLit{}
	switch t := value.(type) {
	case string:
		lit = &syntax.BasicLit{Value: t, Kind: syntax.StringLit}
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
		lit = &syntax.BasicLit{Value: fmt.Sprintf("%v", t), Kind: syntax.IntLit}
	case bool:
		lit = &syntax.BasicLit{Value: fmt.Sprintf("%v", t), Kind: syntax.BoolLit}
	default:
		return nil, fmt.Errorf("unsupport basic lit ")
	}
	return lit, nil
}
