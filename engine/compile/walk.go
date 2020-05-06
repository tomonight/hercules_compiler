// Copyright 2020 tomonight. All rights reserved.

// This file is the gate of the abstract syntax tree to parse to golang code
// Engineer is the excutor to parse tree
// Engineer support:
// (just like an advanced computer language: mostly like javascript)

// if [codition] else if [codition] else [condition]
// for var i=0;i<10;i++{...}  of cause break  continue is support anyway
// function syntax : func a(x,y){...}   no need to define x or y type

// What's new;
// parallel : begin parallel with sub Engineer
// wait : just like lock or synchronized,wait all parallel run over

// system default function:
// length
// sleep
// setTarget
// setRunmode
// connect
// append
// print
// ... on the way
package compile

import (
	"fmt"
	"hercules_compiler/engine/syntax"
	"strings"
)

//Walk is a gate,the main function to run script
func (e *Engineer) Walk() error {
	//TODO
	// consider type check before run script, find static syntax error as fast
	//
	// if e.parse {

	// 	//typecheck
	// 	for _, stmt := range e.node.file.Blocks {
	// 		e.walk(stmt)
	// 	}
	// }

	//the script is a block , we execute these block in a range function
	// check error depends on Engineer run mode
	// checkStatus is check Engineer is pause or need to quit
	for _, stmt := range e.node.file.Blocks {
		_, err := e.walk(stmt)
		err = e.checkError(stmt.Pos(), err)
		if err != nil {
			return err
		}
		if e.checkStatus() {
			break
		}
	}
	return nil
}

//Walk AST tree
func (e *Engineer) walkBlock(stmt *syntax.BlockStmt) (interface{}, error) {
	e.startBlock()
	defer func() {
		e.endBlock()
	}()
	for _, sm := range stmt.List {
		_, err := e.walk(sm)
		err = e.checkError(sm.Pos(), err)
		if err != nil {
			return nil, err
		}
		if e.checkStatus() {
			break
		}
		if e.alreadyReturn() {
			return nil, nil
		}
		if e.isBranch() {
			return nil, nil
		}
	}
	return nil, nil
}

func (e *Engineer) walk(stmt interface{}) (interface{}, error) {
	if stmt == nil {
		return nil, nil
	}
	switch stm := stmt.(type) {
	case *syntax.EmptyStmt:
		return nil, nil
	case *syntax.BlockStmt:
		return e.walkBlock(stm)
	case *syntax.ExprStmt:
		return e.walkExpr(stm)
	case *syntax.DeclStmt:
		return e.walkDecl(stm)
	case *syntax.AssignStmt:
		return e.walkAssign(stm)
	case *syntax.BranchStmt:
		return e.walkBranch(stm)
	case *syntax.ReturnStmt:
		return e.walkReturn(stm)
	case *syntax.IfStmt:
		return e.walkIf(stm)
	case *syntax.ForStmt:
		return e.walkFor(stm)
	case *syntax.ParallelStmt:
		return e.walkParallel(stm)
	case *syntax.WaitStmt:
		return e.walkWait(stm)
	}
	panic("unhandled Stmt")
}

//checkPause or done
func (e *Engineer) checkStatus() bool {
	if e.done {
		return true
	}
	select {
	case <-e.nopauseC:
		e.nopauseC <- struct{}{}
		return false
	}
}

//checkPause or done
func (e *Engineer) checkError(pos syntax.Pos, err error) error {
	if err == nil && !e._EFS_STOP() {
		return nil
	}

	host := ""
	if e.executor != nil {
		host = e.executor.Host
	}
	fun := e.currentFunc()
	packName := ""
	funcName := ""
	if fun != nil {
		funcName = fun.Name.Value
		packName = fun.Package
	}
	errL := []string{"execute script statement on target host '%s' error at "}
	if funcName != "" {
		errL = append(errL, fmt.Sprintf(" %s.%s ", packName, funcName))
	}
	errL = append(errL, "line %d:")
	errMsg := fmt.Sprintf(strings.Join(errL, ""), host, pos.Line())
	fmt.Println(errMsg)
	if err == nil && e._EFS_STOP() {
		e.done = true
		msg := "script execute successful but end with runmode>=10"
		e.callback(e.node.filename, CRT_STATEMENT_FAILED, errMsg+msg)
		// e.callback(e.node.filename, CRT_SCRIPT_FAILED, fmt.Sprintf("execute script statement on target host '%s' error at line %d: %s, statement='%s':script execute successful but end with runmode>=10", pos.Line()))
		return fmt.Errorf(msg)
	}

	if e._EFF_STOP() {
		e.done = true
		e.callback(e.node.filename, CRT_STATEMENT_FAILED, errMsg+err.Error())
		return err
	}
	e.callback(e.node.filename, CRT_STATEMENT_COMPLETED, errMsg+err.Error())
	return nil
}

func (e *Engineer) _EFF_STOP() bool {
	return e.execStatus&1 == 0
}

func (e *Engineer) _EFS_STOP() bool {
	return e.execStatus&(1<<1) == 0
}

