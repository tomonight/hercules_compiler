package compile

import (
	"fmt"
	"hercules_compiler/engine/syntax"
)

//decl statement runner
//only support var declire
func (e *Engineer) walkDecl(stmt *syntax.DeclStmt) (interface{}, error) {

	switch de := stmt.Decl.(type) {
	case *syntax.VarDecl:
		return e.decl(de)
	default:
		return nil, notSupportVarType
	}
}

//decl statement runner
func (e *Engineer) decl(decl *syntax.VarDecl) (interface{}, error) {
	expr := decl.Values
	key := decl.NameList[0].Value
	value, err := e.expr(expr)
	if err != nil {
		return nil, err
	}
	return e.assign(key, value)
}

//assin value into engineer var
func (e *Engineer) assign(key string, value interface{}) (interface{}, error) {
	if syntax.CheckIsConst(key) {
		return nil, fmt.Errorf("can not assign const value")
	}
	fn := e.currentFuncVar()
	bn := e.currentBlock()
	val := &Var{
		fn:    fn,
		key:   key,
		value: value,
	}
	if bn == nil {
		if fn != nil {
			fn.vars[key] = val
		} else {
			e.vars[key] = val
			return nil, nil
		}
	}
	if bn != nil {
		currentExec := e.currentExcute()
		cv := e.funcVars[currentExec]
		list := cv.blocks
		cur := list.Back()
		for cur != nil {
			if _, ok := cur.Value.(*BlockVars).vars[key]; ok {
				cur.Value.(*BlockVars).vars[key] = val
				return nil, nil
			}
			cur = cur.Prev()
		}
		if fn != nil {
			if _, ok := fn.vars[key]; ok {
				fn.vars[key] = val
				return nil, nil
			}
		}
		if _, ok := e.vars[key]; ok {
			e.vars[key] = val
			return nil, nil
		}
		e.addBlockVar(val)
	}
	return nil, nil
}

func (e *Engineer) addBlockVar(val *Var) {
	fn := e.currentExcute()
	bn := e.currentBlock()
	val.block = bn
	val.fn = e.funcVars[fn]
	bn.vars[val.key] = val
}

func checkExsit(l []*Var, key string) (int, bool) {
	if len(l) == 0 {
		return 0, false
	}
	for index, v := range l {
		if v.key == key {
			return index, true
		}
	}
	return 0, false
}

//assin value into engineer var
func (e *Engineer) nameVal(key string) (interface{}, error) {
	if key == "true" {
		return true, nil
	}
	if key == "false" {
		return false, nil
	}

	fn := e.currentExcute()
	list := e.funcVars[fn].blocks
	cur := list.Back()
	// cur := cur.Prev()
	for cur != nil {
		blockVar := cur.Value.(*BlockVars)
		if v, ok := blockVar.vars[key]; ok {
			return v.value, nil
		}
		cur = cur.Prev()
	}
	if e.funcVars[fn] != nil {
		for mk, m := range e.funcVars[fn].vars {
			if mk == key {
				return m.value, nil
			}
		}
	}
	if val, ok := e.vars[key]; ok {
		return val.value, nil
	}
	if constValue := syntax.GetConst(key); constValue != nil {
		return e.expr(constValue.Values)
	}
	return nil, fmt.Errorf("cannot recognize name 【%s】", key)
}

func (e *Engineer) endScope(key string) {
	// if _, ok := e.funcVars[fn]; !ok {
	// 	return
	// }
	// v := e.funcVars[fn]
	// index, ok := checkExsit(v, key)
	// if !ok {
	// 	return
	// }
	// keyv := v[index]
	// keyv.key = ""
	// keyv.value = nil
	// e.funcVars[fn] = v
	return
}

func (e *Engineer) startScope(key string) {

}
