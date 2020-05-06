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
	name       string
	basicName  string
	returnVar  interface{}
	isReturn   bool
	isContinue bool
	isBreak    bool
	blocks     *list.List
	vars       map[string]*Var
}

//Target Target
type Target struct {
	protocol string
	host     string
	port     int
	username string
	password string
	keyfile  string
}
