package compile

import "hercules_compiler/syntax"

//decl statement runner
func (e *Engineer) walkAssign(stmt *syntax.AssignStmt) (interface{}, error) {
	if stmt.Op == 0 {
		return e.assignOp0(stmt.Lhs, stmt.Rhs)
	}
	return e.assignOpNot0(stmt)
}

func (e *Engineer) assignOp0(lhs syntax.Expr, rhs syntax.Expr) (interface{}, error) {
	rv, err := e.expr(rhs)
	if err != nil {
		return nil, err
	}
	lv := lhs.(*syntax.Name)
	e.assign(lv.Value, rv)
	return nil, nil
}

func (e *Engineer) assignOpNot0(stmt *syntax.AssignStmt) (interface{}, error) {
	expr := &syntax.Operation{
		Op: stmt.Op,
		X:  stmt.Lhs,
		Y:  stmt.Rhs,
	}
	lv := stmt.Lhs.(*syntax.Name)
	rv, err := e.expr(expr)
	if err != nil {
		return nil, err
	}
	e.assign(lv.Value, rv)
	return nil, nil
}
