package compile

import "hercules_compiler/engine/syntax"

//decl statement runner
func (e *Engineer) walkReturn(stmt *syntax.ReturnStmt) (interface{}, error) {

	value, err := e.expr(stmt.Results)
	if err != nil {
		return nil, err
	}
	e.setReturn(value)
	return nil, nil
}
