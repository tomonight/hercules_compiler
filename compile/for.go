package compile

import (
	"errors"
	"hercules_compiler/syntax"
)

//walkFor statement runner
func (e *Engineer) walkFor(stmt *syntax.ForStmt) (interface{}, error) {
	e.startBlock()
	defer func() {
		e.endBlock()
	}()
	_, err := e.forInit(stmt.Init)
	if err != nil {
		return nil, err
	}
	for {
		cond, err := e.expr(stmt.Cond)
		if err != nil {
			return nil, err
		}

		if _, ok := cond.(bool); !ok {
			return nil, errors.New("syntax error")
		}

		if !cond.(bool) {
			break
		}

		_, err = e.walkBlock(stmt.Body)
		if err != nil {
			return nil, err
		}
		if e.alreadyReturn() {
			return nil, nil
		}
		_, err = e.walk(stmt.Post)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (e *Engineer) forInit(stmt syntax.SimpleStmt) (interface{}, error) {

	return e.walk(stmt)
}
