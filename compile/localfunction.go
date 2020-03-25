package compile

import (
	"fmt"
	"hercules_compiler/syntax"
)

type function interface {
	call(e *Engineer, args []syntax.Expr) (interface{}, error)
}

type basic struct {
	e    *Engineer
	args []syntax.Expr
}

type oprint struct {
	basic
}

type olength struct {
	basic
}

type oappend struct {
	basic
}

func factory(name string) function {
	switch name {
	case nprint:
		return &oprint{}
	case nlength:
		return &olength{}
	case nappend:
		return &oappend{}
	default:
		//error
		return nil
	}
}

//print function
func (o *oprint) call(e *Engineer, args []syntax.Expr) (interface{}, error) {

	//do nothing
	if len(args) == 0 {
		return nil, nil
	}
	fomats := []interface{}{}
	for _, arg := range args {
		v, err := e.expr(arg)
		if err != nil {
			//error
			return nil, err
		}
		fomats = append(fomats, v)
	}
	fmt.Println(fmt.Sprint(fomats...))
	return nil, nil
}

func (o *olength) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	return nil, nil
}

func (o *oappend) call(e *Engineer, args []syntax.Expr) (interface{}, error) {
	return nil, nil
}
