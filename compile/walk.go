package compile

import (
	"fmt"
	"hercules_compiler/syntax"
)

//Walk AST tree
func (e *Engineer) Walk() error {

	for _, stmt := range e.node.file.Blocks {
		e.walk(stmt)
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
		e.walk(sm)
		if e.alreadyReturn() {
			return nil, nil
		}
	}
	return nil, nil
}

func (e *Engineer) walk(stmt interface{}) (interface{}, error) {
	switch stm := stmt.(type) {
	case *syntax.EmptyStmt:
		return nil, nil
	case *syntax.LabeledStmt:
		fmt.Println("label")
	case *syntax.BlockStmt:
		fmt.Println("block")
		return e.walkBlock(stm)
	case *syntax.ExprStmt:
		return e.walkExpr(stm)
	case *syntax.DeclStmt:
		return e.walkDecl(stm)
	case *syntax.AssignStmt:
		return e.walkAssign(stm)
	case *syntax.BranchStmt:
		fmt.Println("branch")
	case *syntax.CallStmt:
		return e.walkCall(stm)
	case *syntax.ReturnStmt:
		fmt.Println("return")
		return e.walkReturn(stm)
	case *syntax.IfStmt:
		fmt.Println("if")
		return e.walkIf(stm)
	case *syntax.ForStmt:
		fmt.Println("for")
		return e.walkFor(stm)
	case *syntax.ParallelStmt:
		fmt.Println("parallel")
		return e.walkParallel(stm)
	case *syntax.WaitStmt:
		fmt.Println("wait")
		return e.walkWait(stm)
	}
	// panic("unhandled Stmt")

	return nil, nil
}
