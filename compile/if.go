package compile

import (
	"errors"
	"hercules_compiler/syntax"
)

//walkIf statement runner
func (e *Engineer) walkIf(stmt *syntax.IfStmt) (interface{}, error) {

	var init interface{}
	var err error
	switch ex := stmt.Init.(type) {
	case *syntax.ExprStmt:
		init, err = e.expr(ex.X)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("error")
	}
	ok := init.(bool)
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
