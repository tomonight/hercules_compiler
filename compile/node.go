package compile

import (
	"hercules_compiler/src"
	"hercules_compiler/syntax"
	"os"
	"runtime"
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
		p := &noder{
			basemap:  make(map[*syntax.PosBase]*src.PosBase),
			err:      make(chan syntax.Error),
			filename: filename,
		}
		Noders = append(Noders, p)
		go func(filename string) {
			defer func() { wg.Done() }()
			sem <- struct{}{}
			defer func() { <-sem }()
			defer close(p.err)
			base := syntax.NewFileBase(filename)

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
	// for _, p := range Noders {
	// for _ := range p.err {
	// p.yyerrorpos(e.Pos, "%s", e.Msg)
	// }

	// p.node()
	// lines += p.file.Lines
	// p.file = nil // release memory

	// if nsyntaxerrors != 0 {
	// 	errorexit()
	// }
	// // Always run testdclstack here, even when debug_dclstack is not set, as a sanity measure.
	// testdclstack()
	// }

	// localpkg.Height = myheight

	return lines
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
