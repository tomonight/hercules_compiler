package compile

import (
	"fmt"
	"hercules_compiler/engine/src"
	"hercules_compiler/engine/syntax"
	"os"
	"runtime"
	"strings"
	"sync"
)

// A ScopeID represents a lexical scope within a function.
type ScopeID int32

//Noders memory all noder if parse success
var Noders []*noder

var Nodes []*syntax.File

// noder transforms package syntax's AST into a Node tree.
type noder struct {
	basemap   map[*syntax.PosBase]*src.PosBase
	basecache struct {
		last *syntax.PosBase
		base *src.PosBase
	}
	filename string
	file     *syntax.File
	// linknames  []linkname
	pragcgobuf [][]string
	err        chan syntax.Error
	scope      ScopeID

	// scopeVars is a stack tracking the number of variables declared in the
	// current function at the moment each open scope was opened.
	scopeVars []int

	lastCloseScopePos syntax.Pos
}

// ParseFiles concurrently parses files into *syntax.File structures.
// Each declaration in every *syntax.File is converted to a syntax tree
// and its root represented by *Node is appended to xtop.
func ParseFiles(filenames []string) uint {
	// Returns the total count of parsed lines.
	// Limit the number of simultaneously open files.
	sem := make(chan struct{}, runtime.GOMAXPROCS(0)+10)
	var wg sync.WaitGroup
	wg.Add(len(filenames))
	for _, filename := range filenames {
		fn := filename[strings.LastIndex(filename, "/")+1 : strings.LastIndex(filename, ".")]
		p := &noder{
			basemap:  make(map[*syntax.PosBase]*src.PosBase),
			err:      make(chan syntax.Error, 1000),
			filename: fn,
		}
		Noders = append(Noders, p)
		go func(filename string) {
			defer func() { wg.Done() }()
			sem <- struct{}{}
			defer func() { <-sem }()
			defer close(p.err)
			base := syntax.NewFileBase(fn)

			f, err := os.Open(filename)
			if err != nil {
				p.error(syntax.Error{Pos: syntax.MakePos(base, 0, 0), Msg: err.Error()})
				return
			}
			defer f.Close()

			p.file, _ = syntax.Parse(base, f, p.error) // errors are tracked via p.error
		}(filename)
	}
	wg.Wait()
	var lines uint
	for _, p := range Noders {
		le := len(p.err)
		for e := range p.err {
			p.yyerrorpos(e.Pos, "%s", e.Msg)
		}
		if le > 0 {
			panic("parse script error")
		}
		// p.node()
		// lines += p.file.Lines
		// p.file = nil // release memory

		// if nsyntaxerrors != 0 {
		// 	errorexit()
		// }
		// // Always run testdclstack here, even when debug_dclstack is not set, as a sanity measure.
		// testdclstack()
	}

	// localpkg.Height = myheight

	return lines
}

func (p *noder) yyerrorpos(pos syntax.Pos, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(fmt.Sprintf("parse error at  %d:%d,error:%s", pos.Line(), pos.Col(), msg))
}

// error is called concurrently if files are parsed concurrently.
func (p *noder) error(err error) {
	p.err <- err.(syntax.Error)
}

//Errors find all syntax errors
func (p *noder) Errors() []error {
	errs := []error{}
	for len(p.err) > 0 {
		errs = append(errs, <-p.err)
	}
	return errs
}

//GetNoder find noder by pkg name
func GetNoder(name string) *noder {
	for _, node := range Noders {
		if node.filename == name {
			return node
		}
	}
	return nil
}
