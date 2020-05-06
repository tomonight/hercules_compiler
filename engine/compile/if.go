package compile

import (
	"errors"
	"hercules_compiler/engine/syntax"
)

//walkIf statement runner
func (e *Engineer) walkIf(stmt *syntax.IfStmt) (interface{}, error) {

	cond, err := e.expr(stmt.Cond)
	if err != nil {
		return nil, err
	}
	ok := cond.(bool)
	if ok {
		return e.walkBlock(stmt.Then)
	}
	if !ok && stmt.Else != nil {
		return e.walkElse(stmt.Else)
	}
	return nil, nil
}

func (e *Engineer) walkElse(stmt syntax.Stmt) (interface{}, error) {
	if stmt == nil {
		return nil, nil
	}
	switch st := stmt.(type) {
	case *syntax.BlockStmt:
		return e.walkBlock(st)
	case *syntax.IfStmt:
		return e.walkIf(st)
	default:
		return nil, errors.New("error")
	}
}
