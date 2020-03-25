package compile

import (
	"encoding/json"
	"errors"
	"hercules_compiler/syntax"
	"strconv"
)

//walkExpr statement runner
func (e *Engineer) walkExpr(expr *syntax.ExprStmt) (interface{}, error) {
	return e.expr(expr.X)
}

//walkexpr walk step
func (e *Engineer) expr(expr syntax.Expr) (interface{}, error) {

	switch ex := expr.(type) {
	case *syntax.BasicLit:
		return ex.Value, nil
	case *syntax.CallExpr:
		return e.callExpr(ex)
	case *syntax.Name:
		return e.nameVal(ex.Value)
	case *syntax.Operation:
		return e.operationExpr(ex)
	case *syntax.ParenExpr:
		return e.expr(ex.X)
	default:
		return nil, nil
	}
}

func (e *Engineer) callExpr(expr *syntax.CallExpr) (interface{}, error) {

	switch expr.Fun.(type) {
	case *syntax.Name:
		return e.nameCall(expr.Fun.(*syntax.Name), expr.ArgList)
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
	return nil, nil
}

func (e *Engineer) namecall(name *syntax.Name, args []syntax.Expr) (interface{}, error) {
	//user define function call
	fun := syntax.GetFunc(name.Value)
	assign, err := e.defineFormParameter(fun, args)
	if err != nil {
		return nil, err
	}
	e.startFunc()
	defer func() {
		e.endFunc()
	}()
	if len(assign) != 0 {
		for key, value := range assign {
			e.assign(key, value)
		}
	}
	for _, stmt := range fun.Body.List {
		e.walk(stmt)
		if e.alreadyReturn() {
			return e.fnReturn(), nil
		}
	}
	//error
	return nil, nil
}

//define form parameter into enginner
func (e *Engineer) defineFormParameter(fun *syntax.FuncDecl, args []syntax.Expr) (map[string]interface{}, error) {
	if len(fun.Type.ParamList) == 0 {
		return nil, nil
	}
	if len(args) != len(fun.Type.ParamList) {
		return nil, funPramsIsvalid
	}
	assign := make(map[string]interface{})
	for index, p := range fun.Type.ParamList {
		value, err := e.expr(args[index])
		if err != nil {
			return nil, err
		}
		assign[p.Name.Value] = value
	}
	return assign, nil
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
		return xv.(string) + yv.(string), nil
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
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8:
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
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8:
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

func getDigit(value interface{}) (int64, error) {
	switch value.(type) {
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8:
		return value.(int64), nil
	case string:
		if isDigit(value) {
			xv, err := strconv.Atoi(value.(string))
			if err != nil {
				return 0, err
			}
			return int64(xv), nil
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
