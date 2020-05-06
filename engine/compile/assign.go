package compile

import (
	"errors"
	"hercules_compiler/engine/syntax"
)

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
	switch lv := lhs.(type) {
	case *syntax.Name:
		return e.assign(lv.Value, rv)
	case *syntax.IndexExpr:
		return e.assignIndex(lv, rhs)
	}
	return nil, nil
}

func (e *Engineer) assignIndex(lhs *syntax.IndexExpr, rhs syntax.Expr) (interface{}, error) {
	xv, err := e.expr(lhs.X)
	if err != nil {
		return nil, err
	}

	if _, ok := xv.(*syntax.CompositeLit); !ok {
		return nil, errors.New("unsurport type ")
	}
	composite := xv.(*syntax.CompositeLit)
	indexv, err := e.expr(lhs.Index)
	if err != nil {
		return nil, err
	}

	switch composite.Type.(type) {
	case *syntax.MapType:
		for k, exp := range composite.ElemList {
			if _, ok := exp.(*syntax.KeyValueExpr); !ok {
				return nil, errors.New("error map type")
			}
			keyExpr := exp.(*syntax.KeyValueExpr)
			exprv, err := e.expr(keyExpr.Key)
			if err != nil {
				return nil, err
			}

			if indexv == exprv {
				keyExpr.Value = rhs
				composite.ElemList[k] = keyExpr
				return nil, nil
			}
		}
		keyvalue := &syntax.KeyValueExpr{Key: lhs.Index, Value: rhs}
		composite.ElemList = append(composite.ElemList, keyvalue)
	case *syntax.SliceType:
		if !isDigit(indexv) {
			return nil, errors.New("unsurport slice index ")
		}
		indexInt, _ := getDigit(indexv)
		if int(indexInt) < len(composite.ElemList) {
			composite.ElemList[indexInt] = rhs
		} else {
			composite.ElemList = append(composite.ElemList, rhs)
		}
	default:
		return nil, errors.New("unsurport type ")
	}
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
	return e.assign(lv.Value, rv)
}
