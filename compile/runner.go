package compile

import (
	"container/list"

	"github.com/google/uuid"
)

//Engineer Engine runner
type Engineer struct {
	node       *noder
	errC       chan error
	vars       map[string]*Var
	curFn      string
	fnStack    *list.List
	blockStack map[string]*list.List
	funcVars   map[string]*funcVars
	infun      bool
	inblock    bool
}

//Run run
//Engineer start gate
func (e *Engineer) Run() error {
	err := e.Walk()
	return err
}

//NewEngineer initialize Engineer
//file params is the script file AST
func NewEngineer(node *noder, params string) (*Engineer, error) {
	if node == nil {
		return nil, noASTError
	}
	e := &Engineer{
		node:       node,
		vars:       make(map[string]*Var),
		fnStack:    list.New(),
		blockStack: make(map[string]*list.List),
		funcVars:   make(map[string]*funcVars)}
	e.node = node
	e.blockStack[""] = list.New()
	fv := &funcVars{
		name:   "",
		blocks: list.New(),
		vars:   make(map[string]*Var),
	}
	e.funcVars[""] = fv
	err := e.initialParams(params)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Engineer) initialParams(params string) error {

	return nil
}

//func inner param up to script body param
func (e *Engineer) addBodyParam(vars *Var) {

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

func (e *Engineer) startFunc() {
	name := uuid.New().String()
	fv := &funcVars{
		name:   name,
		blocks: list.New(),
		vars:   make(map[string]*Var),
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
