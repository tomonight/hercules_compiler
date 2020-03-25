package compile

import "errors"

//define error type
var (
	noASTError             = errors.New("ast file initialize fail,have no AST file")
	notSupportVarType      = errors.New("not support var type")
	notSupportOprationType = errors.New("not support operation type")
	assignIsvalid          = errors.New("cannot recognize name")
	funPramsIsvalid        = errors.New("function prams is valid")
	boolPramsIsvalid       = errors.New("bool operation need bool param")
)
