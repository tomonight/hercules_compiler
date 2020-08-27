package compile

import (
	"container/list"
	"encoding/json"
	"fmt"
	"hercules_compiler/engine/syntax"
	"runtime"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Callback func(scriptExecID string, callbackResultType int, callbackResultInfo string)

//Engineer Engine runner
type Engineer struct {
	parent      *Engineer
	node        *noder
	errC        chan error
	vars        map[string]*Var
	curFn       string
	fnStack     *list.List
	callback    Callback
	blockStack  map[string]*list.List
	funcVars    map[string]*funcVars
	infun       bool
	inblock     bool
	wg          sync.WaitGroup
	parse       bool
	done        bool
	execStatus  int
	setStatus   bool
	nopauseC    chan struct{}
	parellelC   chan struct{}
	executor    *engineExecutor
	allExecutor []*engineExecutor
	target      string
}

//Run run
//Engineer start gate
func (e *Engineer) Run() {
	defer func() {
		e.disconnectExecutor()
		if e.parent != nil {
			e.parent.wg.Done()
		} else {
			e.walkWait(nil)
		}
	}()
	err := e.Walk()
	if err != nil {
		e.callback(e.node.filename, CRT_SCRIPT_FAILED, "script run with failed")
		return
	}
	e.callback(e.node.filename, CRT_SCRIPT_COMPLETED, "script run completed")
}

//NewEngineer initialize Engineer
//file params is the script file AST
func NewEngineer(node *noder, params interface{}, callback Callback) (*Engineer, error) {
	if node == nil {
		return nil, noASTError
	}
	e := &Engineer{
		node:        node,
		vars:        make(map[string]*Var),
		fnStack:     list.New(),
		blockStack:  make(map[string]*list.List),
		funcVars:    make(map[string]*funcVars),
		callback:    callback,
		nopauseC:    make(chan struct{}, 1),
		allExecutor: []*engineExecutor{},
		parellelC:   make(chan struct{}, runtime.GOMAXPROCS(0)+10)}
	e.blockStack[""] = list.New()
	e.nopauseC <- struct{}{}
	e.execStatus = 10
	fv := &funcVars{
		name:   "",
		blocks: list.New(),
		vars:   make(map[string]*Var),
	}
	e.funcVars[""] = fv
	err := e.addBodyParam(params)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Engineer) initialParams() error {
	if e.parent == nil {
		return nil
	}
	for key, value := range e.parent.vars {
		e.vars[key] = value
	}
	fv := e.parent.currentFuncVar()
	if fv == nil {
		return nil
	}
	for key, value := range fv.vars {
		e.vars[key] = value
	}
	if fv.blocks.Len() == 0 {
		return nil
	}
	list := fv.blocks
	cur := list.Back()
	for cur != nil {
		bv := cur.Value.(*BlockVars)
		for key, value := range bv.vars {
			e.vars[key] = value
		}
		cur = cur.Prev()
	}
	return nil
}

//func inner param up to script body param
func (e *Engineer) addBodyParam(params interface{}) error {
	m := make(map[string]interface{})
	j, err := paramValid(params)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		return err
	}
	for key, value := range m {
		if value == nil {
			continue
		}
		var val *Var
		switch t := value.(type) {
		case bool, string, int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
			val = &Var{
				fn:    nil,
				key:   key,
				value: t,
			}
		case []interface{}:
			lm, err := assembleListMapCompositeLit(value.([]interface{}))
			if err != nil {
				return err
			}
			val = &Var{
				fn:    nil,
				key:   key,
				value: lm,
			}
		case map[string]interface{}:
			lit, err := assembleMapCompositeLitWithInterface(value.(map[string]interface{}))
			if err != nil {
				return err
			}
			val = &Var{
				fn:    nil,
				key:   key,
				value: lit,
			}
		default:
			fmt.Println("unsupport param type,only support key:value,value only support list[map],or string or int")
			continue
		}
		e.vars[key] = val
	}
	return nil
}

func assembleInterfaceParams(data interface{}) (syntax.Expr, error) {
	switch data.(type) {
	case bool, string, int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
		return getBasicLit(data)
	case map[string]interface{}:
		dataV := data.(map[string]interface{})
		return assembleMapCompositeLitWithInterface(dataV)
	case []interface{}:

	}
	return nil, fmt.Errorf("unsupport param type")
}

func assembleMapCompositeLitWithInterface(params map[string]interface{}) (*syntax.CompositeLit, error) {
	lit := &syntax.CompositeLit{}
	t := &syntax.MapType{}
	lit.Type = t
	lit.NKeys = len(params)
	lit.ElemList = []syntax.Expr{}
	for k, v := range params {
		kv := &syntax.KeyValueExpr{}
		kv.Key = &syntax.Name{Value: k}
		value, err := assembleInterfaceParams(v)
		if err != nil {
			return nil, err
		}
		kv.Value = value
		lit.ElemList = append(lit.ElemList, kv)
	}
	return lit, nil
}

func assembleListMapCompositeLit(params []interface{}) (*syntax.CompositeLit, error) {
	lit := &syntax.CompositeLit{}
	list := &syntax.SliceType{}
	lit.Type = list
	lit.NKeys = len(params)
	lit.ElemList = []syntax.Expr{}
	for _, v := range params {
		eachMap, err := assembleInterfaceParams(v)
		if err != nil {
			return nil, err
		}
		// kv := &syntax.KeyValueExpr{}
		// kv.Key = &syntax.Name{Value: k}
		// kv.Value = &syntax.BasicLit{Value: v, Kind: syntax.StringLit}
		lit.ElemList = append(lit.ElemList, eachMap)
	}
	return lit, nil
}

func paramValid(params interface{}) ([]byte, error) {
	switch p := params.(type) {
	case string:
		pstr := strings.Trim(strings.Trim(strings.Trim(p, "\n\t"), "\n"), "")
		fmt.Println(pstr[0:1])
		if strings.HasPrefix(pstr, "{") && strings.HasSuffix(pstr, "}") {
			return []byte(pstr), nil
		}
	case []byte:
		return p, nil
	default:
		return json.Marshal(params)
	}
	return nil, fmt.Errorf("unsupport script params")
}

//func inner param up to script body param
func (e *Engineer) addFuncParam(vars *Var) {

}

func (e *Engineer) currentExcute() string {
	el := e.fnStack.Back()
	if el == nil {
		return ""
	}
	return el.Value.(string)
}

func (e *Engineer) currentFunc() *syntax.FuncDecl {
	fn := e.currentFuncVar()
	if fn == nil {
		return nil
	}
	return syntax.GetFunc(fn.basicName)
}

func (e *Engineer) currentFuncVar() *funcVars {
	fn := e.currentExcute()
	if fn == "" {
		return nil
	}
	return e.funcVars[fn]
}

func (e *Engineer) fnReturn() interface{} {
	fn := e.currentExcute()
	if fn == "" {
		return nil
	}
	return e.funcVars[fn].returnVar
}

func (e *Engineer) currentBlock() *BlockVars {
	el := e.funcVars[e.currentExcute()].blocks.Back()
	if el == nil {
		return nil
	}
	return el.Value.(*BlockVars)
}

func (e *Engineer) startFunc(basciName string) {
	name := uuid.New().String()
	fv := &funcVars{
		name:      name,
		blocks:    list.New(),
		basicName: basciName,
		vars:      make(map[string]*Var),
	}
	e.funcVars[name] = fv
	e.fnStack.PushBack(name)
}

func (e *Engineer) endFunc() {
	e.fnStack.Remove(e.fnStack.Back())
}

func (e *Engineer) startBlock() {
	bl := &BlockVars{
		fn:   e.currentExcute(),
		vars: make(map[string]*Var),
	}
	e.funcVars[e.currentExcute()].blocks.PushBack(bl)
}

func (e *Engineer) endBlock() {
	fn := e.currentExcute()
	back := e.funcVars[fn].blocks.Back()
	e.funcVars[fn].blocks.Remove(back)
}

func (e *Engineer) alreadyReturn() bool {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	return back.isReturn
}

func (e *Engineer) setReturn(value interface{}) {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	back.isReturn = true
	back.returnVar = value
}

func (e *Engineer) setBreak() {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	back.isBreak = true
}

func (e *Engineer) endBreak() {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	back.isBreak = false
}

func (e *Engineer) isBreak() bool {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	return back.isBreak
}

func (e *Engineer) setContinue() {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	back.isContinue = true
}

func (e *Engineer) endContinue() {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	back.isContinue = false
}

func (e *Engineer) isContinue() bool {
	fn := e.currentExcute()
	back := e.funcVars[fn]
	return back.isContinue
}

func (e *Engineer) isBranch() bool {
	return e.isContinue() || e.isBreak()
}

func (e *Engineer) endBranch() {
	e.endContinue()
	e.endBreak()
}

func (e *Engineer) perror(msg string) {
	e.done = true
	e.callback(e.node.filename, CRT_PARSE_FAILED, msg)
	e.callback(e.node.filename, CRT_SCRIPT_FAILED, msg)
}

func (e *Engineer) runError(msg string) {
	e.done = true
	e.callback(e.node.filename, CRT_STATEMENT_FAILED, msg)
	e.callback(e.node.filename, CRT_SCRIPT_FAILED, msg)
}

func (e *Engineer) edone(msg string) {
	e.callback(e.node.filename, CRT_STATEMENT_COMPLETED, msg)
}

func (e *Engineer) Pause(msg string) {
	<-e.nopauseC
}

func (e *Engineer) Continue(msg string) {
	e.nopauseC <- struct{}{}
}

func (e *Engineer) addExecutor(executor *engineExecutor) {
	e.allExecutor = append(e.allExecutor, executor)
}

func (e *Engineer) checkExecutorExsit(host string, port int) *engineExecutor {
	for _, ex := range e.allExecutor {
		if host == ex.Host && port == ex.port {
			return ex
		}
	}
	return nil
}

func (e *Engineer) disconnectExecutor() {
	for _, executor := range e.allExecutor {
		executor.executor.Close()
	}
}
