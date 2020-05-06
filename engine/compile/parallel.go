package compile

import (
	"errors"
	"hercules_compiler/engine/src"
	"hercules_compiler/engine/syntax"
)

//walkParallel statement runner
func (e *Engineer) walkParallel(stmt *syntax.ParallelStmt) (interface{}, error) {
	subAST := &noder{
		file:    &syntax.File{Blocks: stmt.Body.List},
		basemap: make(map[*syntax.PosBase]*src.PosBase),
		err:     make(chan syntax.Error),
	}
	subEngineer, err := NewEngineer(subAST, "", e.callback)
	if err != nil {
		return nil, errors.New("Parallel error")
	}
	subEngineer.parent = e
	err = subEngineer.initialParams()
	if err != nil {
		return nil, err
	}
	e.wg.Add(1)
	go subEngineer.Run()
	return nil, nil
}
