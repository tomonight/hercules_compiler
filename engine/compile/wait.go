package compile

import (
	"hercules_compiler/engine/syntax"
)

//walkWait statement runner
func (e *Engineer) walkWait(stmt *syntax.WaitStmt) (interface{}, error) {

	e.wg.Wait()
	return nil, nil
}
