package compile

import "container/list"

type Var struct {
	scope Scope
	key   string
	fn    *funcVars
	block *BlockVars
	value interface{}
}

type BlockVars struct {
	fn   string
	name string
	vars map[string]*Var
}

type funcVars struct {
	name      string
	returnVar interface{}
	isReturn  bool
	blocks    *list.List
	vars      map[string]*Var
}
