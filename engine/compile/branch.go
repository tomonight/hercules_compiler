package compile

import "hercules_compiler/engine/syntax"

//decl statement runner
//only support var declire
func (e *Engineer) walkBranch(stmt *syntax.BranchStmt) (interface{}, error) {

	switch stmt.Tok {
	case syntax.Break:
		e.setBreak()
		return nil, nil
	case syntax.Continue:
		e.setContinue()
		return nil, nil
	default:
		return nil, notSupportVarType
	}
}
